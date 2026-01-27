package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/dafraer/sentence-gen-grpc-server/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSentenceGenServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GenerateSentence(ctx context.Context, request *pb.GenerateSentenceRequest) (*pb.GenerateSentenceResponse, error) {
	fmt.Println("GenerateSentence")
	return &pb.GenerateSentenceResponse{}, nil
}

func (s *Server) TranslateWord(ctx context.Context, request *pb.TranslateWordRequest) (*pb.TranslateWordResponse, error) {
	fmt.Println("TranslateWord")
	return &pb.TranslateWordResponse{}, nil
}

func (s *Server) GenerateDefinition(ctx context.Context, request *pb.GenerateDefinitionRequest) (*pb.GenerateDefinitionResponse, error) {
	fmt.Println("GenerateDefinition")
	return &pb.GenerateDefinitionResponse{Definition: request.DefinitionHint}, nil
}

func (s *Server) Run() {
	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	pb.RegisterSentenceGenServer(srv, &Server{})
	defer srv.Stop()
	if err := srv.Serve(l); err != nil {
		panic(err)
	}
}
