package server

import (
	"context"
	api "github.com/rezaAmiri123/proglog/api/v1"
	"github.com/rezaAmiri123/proglog/internal/auth"
	"github.com/rezaAmiri123/proglog/internal/config"
	"github.com/rezaAmiri123/proglog/internal/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T, rootClient api.LogClient, nobodyClient api.LogClient, config *Config){
		"produce/consume a message to/form the log succeeds": testProduceConsume,
		"produce/consume stream succeeds":                    testProduceConsumeStream,
		"consume past log boundary fails":                    testConsumePastBoundary,
		"unauthorized fails":                                 testUnauthorized,
	} {
		t.Run(scenario, func(t *testing.T) {
			rootClient, nobodyClient, cfg, teardown := setupTest(t, nil)
			defer teardown()
			fn(t, rootClient, nobodyClient, cfg)
		})
	}

}

func setupTest(t *testing.T, fn func(*Config)) (
	rootClient api.LogClient,
	nobodyClient api.LogClient,
	cfg *Config,
	teardown func(),
) {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	newClient := func(crtPath, keyPath string) (
		*grpc.ClientConn,
		api.LogClient,
		[]grpc.DialOption,
	) {
		tlsConfig, err := config.SetupTLSConfig(config.TLSConfig{
			CertFile: crtPath,
			KeyFile:  keyPath,
			CAFile:   config.CAFile,
			Server:   false,
		})
		require.NoError(t, err)

		tlsCreds := credentials.NewTLS(tlsConfig)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}

		conn, err := grpc.Dial(l.Addr().String(), opts...)
		require.NoError(t, err)

		client := api.NewLogClient(conn)
		return conn, client, opts

	}

	var rootConn, nobodyConn *grpc.ClientConn
	rootConn, rootClient, _ = newClient(config.RootClientCertFile, config.RootClientKeyFile)
	nobodyConn, nobodyClient, _ = newClient(config.NobodyClientCertFile, config.NobodyClientKeyFile)

	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile: config.ServerCertFile,
		KeyFile:  config.ServerKeyFile,
		CAFile:   config.CAFile,
		Server:   true,
	})
	require.NoError(t, err)

	serverCreds := credentials.NewTLS(serverTLSConfig)

	dir, err := ioutil.TempDir("", "server-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)

	authorizer := auth.New(config.ACLModelFile, config.ACLPolicyFile)

	cfg = &Config{
		CommitLog:  clog,
		Authorizer: authorizer,
	}

	if fn != nil {
		fn(cfg)
	}

	server, err := NewGRPCServer(cfg, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go func() {
		server.Serve(l)
	}()

	return rootClient, nobodyClient, cfg, func() {
		server.Stop()
		rootConn.Close()
		nobodyConn.Close()
		l.Close()
	}
}

func testProduceConsume(t *testing.T, rootClient, nobodyClient api.LogClient, config *Config) {
	ctx := context.Background()

	want := &api.Record{
		Value: []byte("hello world"),
	}

	produce, err := rootClient.Produce(ctx, &api.ProduceRequest{Record: want})
	require.NoError(t, err)

	consume, err := rootClient.Consume(ctx, &api.ConsumeRequest{Offset: produce.GetOffset()})
	require.NoError(t, err)
	require.Equal(t, want.GetValue(), consume.GetRecord().GetValue())
	require.Equal(t, want.GetOffset(), consume.GetRecord().GetOffset())
}

func testConsumePastBoundary(t *testing.T, rootClient, nobodyClient api.LogClient, config *Config) {
	ctx := context.Background()

	want := &api.Record{Value: []byte("hello world")}

	produce, err := rootClient.Produce(ctx, &api.ProduceRequest{Record: want})
	require.NoError(t, err)

	consume, err := rootClient.Consume(ctx, &api.ConsumeRequest{Offset: produce.GetOffset() + 1})
	require.Nil(t, consume)
	require.Error(t, err)

	//gotErr := grpc.Code(err)
	//wantErr := grpc.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	//if gotErr != wantErr {
	//	t.Fatalf("got err: %v, want: %v", gotErr, wantErr)
	//}
}

func testProduceConsumeStream(t *testing.T, rootClient, nobodyClient api.LogClient, config *Config) {
	ctx := context.Background()

	records := []*api.Record{{
		Value:  []byte("first message"),
		Offset: 0,
	}, {
		Value:  []byte("second message"),
		Offset: 01,
	}}
	{
		stream, err := rootClient.ProduceStream(ctx)
		for offset, record := range records {
			err = stream.Send(&api.ProduceRequest{Record: record})
			require.NoError(t, err)

			res, err := stream.Recv()
			require.NoError(t, err)

			if res.Offset != uint64(offset) {
				t.Fatalf("got offset: %d, want: %d", res.Offset, offset)
			}
		}
	}
	{
		stream, err := rootClient.ConsumeStream(ctx, &api.ConsumeRequest{Offset: 0})
		require.NoError(t, err)

		for i, record := range records {
			res, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, res.GetRecord(), &api.Record{
				Value:  record.GetValue(),
				Offset: uint64(i),
			})
		}
	}
}

func testUnauthorized(t *testing.T, rootClient, nobodyClient api.LogClient, config *Config) {
	ctx := context.Background()

	want := &api.Record{
		Value: []byte("hello world"),
	}

	produce, err := nobodyClient.Produce(ctx, &api.ProduceRequest{Record: want})
	require.Nil(t, produce)
	gotCode, wantCode := status.Code(err), codes.PermissionDenied
	require.Equal(t, gotCode, wantCode)

	consume, err := nobodyClient.Consume(ctx, &api.ConsumeRequest{Offset: 0})
	require.Nil(t, consume)
	gotCode, wantCode = status.Code(err), codes.PermissionDenied
	require.Equal(t, gotCode, wantCode)
}
