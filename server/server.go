package server

import (
	"context"
	"errors"
	"net"

	pb "github.com/dafraer/sentence-gen-grpc-server/proto"
	"github.com/dafraer/sentence-gen-grpc-server/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	result, err := s.srvc.GenerateSentence(ctx, &service.GenerateSentenceRequest{
		Word:                request.Word,
		WordLanguage:        request.WordLanguage,
		TranslationLanguage: request.TranslationLanguage,
		TranslationHint:     request.TranslationHint,
		IncludeAudio:        request.IncludeAudio,
		VoiceGender:         service.Gender(request.VoiceGender),
	})
	if err != nil {
		return nil, formatError(err)
	}
	resp := &pb.GenerateSentenceResponse{
		OriginalSentence:   result.OriginalSentence,
		TranslatedSentence: result.TranslatedSentence,
		Audio: &pb.Audio{
			Data: result.Audio,
		},
	}
	return resp, nil
}

func (s *Server) Translate(ctx context.Context, request *pb.TranslateRequest) (*pb.TranslateResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	result, err := s.srvc.Translate(ctx, &service.TranslateRequest{
		Word:            request.Word,
		FromLanguage:    request.FromLanguage,
		ToLanguage:      request.ToLanguage,
		TranslationHint: request.TranslationHint,
		IncludeAudio:    request.IncludeAudio,
		VoiceGender:     service.Gender(request.VoiceGender),
	})
	if err != nil {
		return nil, formatError(err)
	}
	resp := &pb.TranslateResponse{
		Translation: result.Translation,
		Audio: &pb.Audio{
			Data: result.Audio,
		},
	}
	return resp, nil
}

func (s *Server) GenerateDefinition(ctx context.Context, request *pb.GenerateDefinitionRequest) (*pb.GenerateDefinitionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	result, err := s.srvc.GenerateDefinition(ctx, &service.GenerateDefinitionRequest{
		Word:           request.Word,
		Language:       request.Language,
		DefinitionHint: request.DefinitionHint,
		IncludeAudio:   request.IncludeAudio,
	})
	if err != nil {
		return nil, formatError(err)
	}
	resp := &pb.GenerateDefinitionResponse{
		Definition: result.Definition,
		Audio: &pb.Audio{
			Data: result.Audio,
		},
	}

	return resp, nil
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

func formatError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidRequest) || errors.Is(err, service.ErrInvalidResponse):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
