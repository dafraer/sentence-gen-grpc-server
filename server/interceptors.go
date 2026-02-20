package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) quotaLimitInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s.logger.Debugw("quota interceptor check started", "method", info.FullMethod)
	exceeded, err := s.srvc.DailyQuotaExceeded(ctx)
	if err != nil {
		s.logger.Errorw("quota interceptor failed to check daily quota", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exceeded {
		s.logger.Infow("quota exceeded request blocked", "method", info.FullMethod)
		return nil, status.Errorf(codes.ResourceExhausted, "daily quota limit exceeded")
	}
	s.logger.Debugw("quota interceptor passed", "method", info.FullMethod)
	return handler(ctx, req)
}
