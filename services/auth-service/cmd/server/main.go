package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	logging "github.com/exPriceD/Streaming-platform/pkg/logger"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/transport/grpc"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	network         = "tcp"
	shutdownTimeout = 5 * time.Second
)

func main() {
	logger := logging.InitLogger("auth-service")

	configPath := flag.String("config", "dev", "path to config file or environment name (e.g., 'dev', 'prod', '/path/to/config.yaml')")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath) // dev, prod, test, docker
	if err != nil {
		logger.Error("❌ Couldn't load the configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("✅ Configuration loaded successfully")

	database, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		logger.Error("❌ Database connection error", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("✅ Database connection established")

	tokenRepo := repository.NewTokenRepository(database)
	jwtManager := token.NewJWTManager(cfg.JWT, logger)
	logger.Info("🔧 Repositories are initialized")

	authService := service.NewAuthService(tokenRepo, jwtManager, logger)
	logger.Info("🔧 Services are initialized")

	grpcHandler := grpcTransport.NewAuthHandler(authService, logger)
	logger.Info("🔧 GRPC handlers are initialized")

	grpcServer := grpcTransport.NewGRPCServer(grpcHandler, logger)
	logger.Info("🔧 gRPC Server is initialized")

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic in gRPC server", slog.Any("panic", r))
			}
		}()
		grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
		if err := grpcServer.Run(grpcAddr); err != nil {
			logger.Error("❌ gRPC Server error", slog.String("error", err.Error()))
		} else {
			logger.Info("🚀 Auth-service is running", slog.String("network", network), slog.String("address", grpcAddr))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := grpcServer.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown gRPC server", slog.String("error", err.Error()))
	}

	wg.Wait()

	if err := database.Close(); err != nil {
		logger.Error("Couldn't close the database", slog.String("error", err.Error()))
	} else {
		logger.Info("✅ The database connection is closed")
	}

	logger.Info("Server shut down gracefully")
}
