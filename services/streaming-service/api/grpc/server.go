package grpc

import (
	"log"
	"net"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/proto"
	"google.golang.org/grpc"
)

// StartGRPCServer запускает gRPC сервер
func StartGRPCServer(streamService *service.StreamService, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterStreamingServiceServer(grpcServer, NewStreamHandler(streamService))

	log.Printf("gRPC сервер запущен на %s", addr)
	return grpcServer.Serve(listener)
}
