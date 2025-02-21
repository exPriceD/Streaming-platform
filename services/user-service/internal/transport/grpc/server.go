package grpcTransport

import (
	"context"
	pb "github.com/exPriceD/Streaming-platform/pkg/proto/v1/user"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	server  *grpc.Server
	handler *Handler
	logger  *slog.Logger
}

func NewGRPCServer(handler *Handler, logger *slog.Logger) *Server {
	srv := &Server{
		server:  grpc.NewServer(),
		handler: handler,
		logger:  logger,
	}
	pb.RegisterUserServiceServer(srv.server, srv)
	return srv
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
