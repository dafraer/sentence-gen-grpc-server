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

func (s *Server) Translate(ctx context.Context, request *pb.TranslateRequest) (*pb.TranslateResponse, error) {
	fmt.Println("TranslateWord")
	return &pb.TranslateResponse{}, nil
}

func (s *Server) GenerateDefinition(ctx context.Context, request *pb.GenerateDefinitionRequest) (*pb.GenerateDefinitionResponse, error) {
	fmt.Println("GenerateDefinition")
	return &pb.GenerateDefinitionResponse{Definition: request.DefinitionHint}, nil
}

func (s *Server) Run(ctx context.Context, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

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
	s.logger.Infow("Server is running")

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
