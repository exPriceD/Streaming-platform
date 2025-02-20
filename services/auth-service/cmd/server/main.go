package main

import (
	"database/sql"
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	logging "github.com/exPriceD/Streaming-platform/pkg/logger"
	"github.com/exPriceD/Streaming-platform/pkg/proto/v1/auth"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/transport/grpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

var (
	network = "tcp"
)

func main() {
	logger := logging.InitLogger("auth-service")

	cfg, err := config.LoadConfig("dev") // dev, prod, test
	if err != nil {
		logger.Error("❌ Couldn't load the configuration", slog.String("error", err.Error()))
	}
	logger.Info("✅ Configuration loaded successfully")

	database, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		logger.Error("❌ Database connection error", slog.String("error", err.Error()))
		return
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			logger.Error("Couldn't close the database", slog.String("error", err.Error()))
		} else {
			logger.Info("✅ The database connection is closed")
		}
	}(database)

	tokenRepo := repository.NewTokenRepository(database, logger)
	jwtManager := token.NewJWTManager(cfg.JWT, logger)

	authService := service.NewAuthService(tokenRepo, jwtManager, logger)

	logger.Info("🔧 Repositories and services are initialized")

	addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	lis, err := net.Listen(network, addr)
	if err != nil {
		logger.Error("❌ Couldn't start the server", slog.String("error", err.Error()))
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, handler.NewAuthHandler(authService, logger))

	logger.Info("🚀 Auth-service is running", slog.String("network", network), slog.String("address", addr))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("❌ Server error", slog.String("error", err.Error()))
	}
}
