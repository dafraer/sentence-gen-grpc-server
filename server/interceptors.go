package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) quotaLimitInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	exceeded, err := s.srvc.DailyQuotaExceeded(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exceeded {
		return nil, status.Errorf(codes.ResourceExhausted, "daily quota limit exceeded")
	}
	return handler(ctx, req)
}
