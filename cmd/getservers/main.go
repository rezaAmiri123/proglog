package main

import (
	"context"
	"flag"
	"fmt"
	api "github.com/rezaAmiri123/proglog/api/v1"
	"google.golang.org/grpc"
	"log"
)

func main(){
	addr := flag.String("addr", ":8400", "service address")
	flag.Parse()
	conn,err := grpc.Dial(*addr,grpc.WithInsecure())
	if err != nil{
		log.Fatal(err)
	}
	client := api.NewLogClient(conn)
	ctx := context.Background()
	res,err := client.GetServers(ctx, &api.GetServersRequest{})
	if err !=nil{
		log.Fatal(err)
	}
	fmt.Println("Servers:")
	for _, server := range res.Servers{
		fmt.Printf("\t- %v\n", server)
	}
}
