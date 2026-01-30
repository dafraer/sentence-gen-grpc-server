package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/dafraer/sentence-gen-grpc-server/proto"
	"github.com/dafraer/sentence-gen-grpc-server/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSentenceGenServer
	srvc   *service.Service
	logger *zap.SugaredLogger
}

func NewServer(srvc *service.Service, logger *zap.SugaredLogger) *Server {
	return &Server{srvc: srvc, logger: logger}
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

func (s *Server) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}
	defer func(l net.Listener) {
		if err := l.Close(); err != nil {
			s.logger.Errorw("Failed to close listener", "error", err)
			panic(err)
		}
	}(l)

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	pb.RegisterSentenceGenServer(srv, s)

	//Create a channel to listen for errors
	ch := make(chan error)

	//Run the server in a separate goroutine
	go func() {
		defer close(ch)
		if err := srv.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			ch <- err
			return
		}
		ch <- nil
	}()
	s.logger.Infow("Service is running")

	//Wait for the context to be done or for an error to occur and shutdown the server
	select {
	case <-ctx.Done():
		srv.Stop()
		err := <-ch
		if err != nil {
			return err
		}
	case err := <-ch:
		return err
	}
	return nil
}
