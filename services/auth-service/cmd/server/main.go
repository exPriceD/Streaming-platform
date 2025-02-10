package main

import (
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/handler"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	database, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	userRepo := repository.NewUserRepository(database)
	tokenRepo := repository.NewTokenRepository(database)
	jwtManager := token.NewJWTManager()

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Couldn't start the server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler.NewAuthHandler(authService))

	log.Println("Auth-service is running on the port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
