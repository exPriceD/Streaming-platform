package main

import (
	"database/sql"
	"fmt"
	"github.com/exPriceD/Streaming-platform/config"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/exPriceD/Streaming-platform/pkg/logger"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/handler"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	pb "github.com/exPriceD/Streaming-platform/services/auth-service/proto"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

var (
	network = "tcp"
)

func main() {
	log := logger.InitLogger("auth-service")

	cfg, err := config.LoadAuthConfig()
	if err != nil {
		log.Error("❌ Couldn't load the configuration", slog.String("error", err.Error()))
	}
	log.Info("✅ Configuration loaded successfully")

	database, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Error("❌ Database connection error", slog.String("error", err.Error()))
		return
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Error("Couldn't close the database", slog.String("error", err.Error()))
		} else {
			log.Info("✅ The database connection is closed")
		}
	}(database)

	userRepo := repository.NewUserRepository(database)
	tokenRepo := repository.NewTokenRepository(database)
	jwtManager := token.NewJWTManager(cfg.JWT)

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)

	log.Info("🔧 Repositories and services are initialized")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen(network, addr)
	if err != nil {
		log.Error("❌ Couldn't start the server", slog.String("error", err.Error()))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler.NewAuthHandler(authService))

	log.Info("🚀 Auth-service is running", slog.String("network", network), slog.String("address", addr))
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("❌ Server error", slog.String("error", err.Error()))
	}
}
