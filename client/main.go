package main

import (
	pb "github.com/dafraer/sentence-gen-grpc-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer client.Close()
	helloWorldClient := pb.NewSentenceGenClient(client)
}
