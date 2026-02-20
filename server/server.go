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
		s.logger.Errorw("generate sentence rpc failed: nil request", "error", errors.New("nil request"))
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}
	s.logger.Infow("generate sentence rpc request received", "word", request.Word, "word_language", request.WordLanguage, "translation_language", request.TranslationLanguage, "include_audio", request.IncludeAudio)

	result, err := s.srvc.GenerateSentence(ctx, &service.GenerateSentenceRequest{
		Word:                request.Word,
		WordLanguage:        request.WordLanguage,
		TranslationLanguage: request.TranslationLanguage,
		TranslationHint:     request.TranslationHint,
		IncludeAudio:        request.IncludeAudio,
		VoiceGender:         service.Gender(request.VoiceGender),
	})
	if err != nil {
		s.logger.Errorw("generate sentence rpc failed", "error", err)
		return nil, formatError(err)
	}
	resp := &pb.GenerateSentenceResponse{
		OriginalSentence:   result.OriginalSentence,
		TranslatedSentence: result.TranslatedSentence,
		Audio: &pb.Audio{
			Data: result.Audio,
		},
	}
	s.logger.Infow("generate sentence rpc completed", "has_audio", len(result.Audio) > 0)
	return resp, nil
}

func (s *Server) Translate(ctx context.Context, request *pb.TranslateRequest) (*pb.TranslateResponse, error) {
	if request == nil {
		s.logger.Errorw("translate rpc failed: nil request", "error", errors.New("nil request"))
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}
	s.logger.Infow("translate rpc request received", "word", request.Word, "from_language", request.FromLanguage, "to_language", request.ToLanguage, "include_audio", request.IncludeAudio)

	result, err := s.srvc.Translate(ctx, &service.TranslateRequest{
		Word:            request.Word,
		FromLanguage:    request.FromLanguage,
		ToLanguage:      request.ToLanguage,
		TranslationHint: request.TranslationHint,
		IncludeAudio:    request.IncludeAudio,
		VoiceGender:     service.Gender(request.VoiceGender),
	})
	if err != nil {
		s.logger.Errorw("translate rpc failed", "error", err)
		return nil, formatError(err)
	}
	resp := &pb.TranslateResponse{
		Translation: result.Translation,
		Audio: &pb.Audio{
			Data: result.Audio,
		},
	}
	s.logger.Infow("translate rpc completed", "has_audio", len(result.Audio) > 0)
	return resp, nil
}

func (s *Server) GenerateDefinition(ctx context.Context, request *pb.GenerateDefinitionRequest) (*pb.GenerateDefinitionResponse, error) {
	if request == nil {
		s.logger.Errorw("generate definition rpc failed: nil request", "error", errors.New("nil request"))
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}
	s.logger.Infow("generate definition rpc request received", "word", request.Word, "language", request.Language, "include_audio", request.IncludeAudio)

	result, err := s.srvc.GenerateDefinition(ctx, &service.GenerateDefinitionRequest{
		Word:           request.Word,
		Language:       request.Language,
		DefinitionHint: request.DefinitionHint,
		IncludeAudio:   request.IncludeAudio,
		VoiceGender:    service.Gender(request.VoiceGender),
	})
	if err != nil {
		s.logger.Errorw("generate definition rpc failed", "error", err)
		return nil, formatError(err)
	}
	resp := &pb.GenerateDefinitionResponse{
		Definition: result.Definition,
		Audio: &pb.Audio{
			Data: result.Audio,
		},
	}
	s.logger.Infow("generate definition rpc completed", "has_audio", len(result.Audio) > 0)

	return resp, nil
}

func (s *Server) Run(ctx context.Context, addr string) error {
	s.logger.Infow("starting grpc server", "address", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Errorw("failed to start tcp listener", "error", err)
		return err
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(s.quotaLimitInterceptor),
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterSentenceGenServer(srv, s)

	//Create a channel to listen for errors
	ch := make(chan error)

	//Run the server in a separate goroutine
	go func() {
		defer close(ch)
		if err := srv.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Errorw("grpc server serve failed", "error", err)
			ch <- err
			return
		}
		s.logger.Debugw("grpc server serve loop exited")
		ch <- nil
	}()
	s.logger.Infow("Server is running")

	//Wait for the context to be done or for an error to occur and shutdown the server
	select {
	case <-ctx.Done():
		s.logger.Infow("grpc server context canceled, stopping server")
		srv.Stop()
		err := <-ch
		if err != nil {
			s.logger.Errorw("grpc server stopped with error", "error", err)
			return err
		}
	case err := <-ch:
		if err != nil {
			s.logger.Errorw("grpc server exited with error", "error", err)
		}
		return err
	}
	s.logger.Infow("grpc server stopped")
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
