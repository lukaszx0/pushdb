package main

import (
	"google.golang.org/grpc"

	"context"
	"io"
	"log"

	pb "github.com/lukaszx0/pushdb/proto"
)

var (
	serverAddr = "localhost:50051"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPushdbServiceClient(conn)
	stream, err := client.Watch(context.Background(), &pb.RegisterWatchRequest{KeyName: "test"})
	if err != nil {
		log.Println("stream err: %v", err)
	}

	for {
		watch, err := stream.Recv()
		if err == io.EOF {
			log.Printf("EOF")
			break
		}
		if err != nil {
			log.Println("client=%v err=%v", client, err)
			break
		}
		log.Println(watch)
	}
}
