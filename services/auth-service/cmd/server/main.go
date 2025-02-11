package main

import (
	"database/sql"
	"fmt"
	"github.com/exPriceD/Streaming-platform/config"
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
	cfg, err := config.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Couldn't load the configuration: %v", err)
	}

	database, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("Couldn't close the database: %v", err)
		}
	}(database)

	userRepo := repository.NewUserRepository(database)
	tokenRepo := repository.NewTokenRepository(database)
	jwtManager := token.NewJWTManager(cfg.JWT)

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
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
