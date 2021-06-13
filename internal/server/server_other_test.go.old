package server

import (
	"context"
	api "github.com/rezaAmiri123/proglog/api/v1"
	"github.com/rezaAmiri123/proglog/internal/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io/ioutil"
	"testing"
)

// TODO should do something to fix authorize in grpc handlers
//func TestAllServer(t *testing.T){
//	server, config, teardown := setupTestServer(t, nil)
//	defer teardown()
//	//testProduceConsumeServer(t, server, config)
//	testConsumePastBoundaryServer(t, server, config)
//	//testProduceConsumeStream(t, client, config)
//}

func setupTestServer(t *testing.T, fn func(config *Config)) (
	server *grpcServer,
	cfg *Config,
	teardown func(),
) {
	t.Helper()

	dir, err := ioutil.TempDir("", "server-test")
	require.NoError(t, err)

	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)

	cfg = &Config{
		CommitLog: clog,
	}
	server, err = newgrpcServer(cfg)
	require.NoError(t, err)

	return server, cfg, func() {
		clog.Remove()
	}
}

func testProduceConsumeServer(t *testing.T, server *grpcServer, config *Config) {
	ctx := context.Background()
	want := &api.Record{
		Value: []byte("hello world"),
	}
	produce, err := server.Produce(
		ctx,
		&api.ProduceRequest{
			Record: want,
		},
	)
	require.NoError(t, err)

	consume, err := server.Consume(ctx, &api.ConsumeRequest{
		Offset: produce.Offset,
	})
	require.NoError(t, err)
	require.Equal(t, want.Value, consume.Record.Value)
	require.Equal(t, want.Offset, consume.Record.Offset)
}

func testConsumePastBoundaryServer(t *testing.T, server *grpcServer, config *Config) {
	ctx := context.Background()
	produce, err := server.Produce(ctx, &api.ProduceRequest{
		Record: &api.Record{
			Value: []byte("hello world"),
		},
	})
	require.NoError(t, err)

	consume, err := server.Consume(ctx, &api.ConsumeRequest{
		Offset: produce.Offset + 1,
	})
	if consume != nil {
		t.Fatal("consume not nil")
	}
	got := grpc.Code(err)
	want := grpc.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	if got != want {
		t.Fatalf("got err: %v, want: %v", got, want)
	}
}

// TODO this test must be fixed
//func testProduceConsumeStreamServer(t *testing.T, server *grpcServer, config *Config) {
//	ctx := context.Background()
//	records := []*api.Record{{
//		Value:  []byte("first message"),
//		Offset: 0,
//	}, {
//		Value:  []byte("second message"),
//		Offset: 1,
//	}}
//	{
//		err := server.ProduceStream(ctx)
//		require.NoError(t, err)
//
//		for offset, record := range records {
//			err = stream.Send(&api.ProduceRequest{
//				Record: record,
//			})
//			require.NoError(t, err)
//			res, err := stream.Recv()
//			require.NoError(t, err)
//			if res.Offset != uint64(offset) {
//				t.Fatalf(
//					"got offset: %d, want: %d",
//					res.Offset,
//					offset,
//				)
//			}
//		}
//	}
//	{
//		stream, err := client.ConsumeStream(
//			ctx,
//			&api.ConsumeRequest{Offset: 0},
//		)
//		require.NoError(t, err)
//		for i, record := range records {
//			res, err := stream.Recv()
//			require.NoError(t, err)
//			require.Equal(t, res.Record, &api.Record{
//				Value:  record.Value,
//				Offset: uint64(i),
//			})
//		}
//	}
//}
//
