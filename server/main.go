package main

import (
	"net"

	pb "github.com/dafraer/sentence-gen-grpc-server/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSentenceGenServer
}

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	pb.RegisterSentenceGenServer(srv, &server{})
	defer srv.Stop()
	if err := srv.Serve(l); err != nil {
		panic(err)
	}
}
