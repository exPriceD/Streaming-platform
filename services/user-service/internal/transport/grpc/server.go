package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	pb.UnimplementedUserServiceServer // Встраиваем базовую реализацию из generated кода
	logger                            *slog.Logger
	server                            *grpc.Server
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("Failed to listen", slog.String("error", err.Error()))
		return err
	}
	s.logger.Info("gRPC server is running", slog.String("address", addr))
	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down gRPC server")
	s.server.GracefulStop()
	return nil
}
