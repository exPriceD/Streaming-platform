package main

import (
	"context"
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	logging "github.com/exPriceD/Streaming-platform/pkg/logger"
	_ "github.com/exPriceD/Streaming-platform/services/user-service/docs"
	cl "github.com/exPriceD/Streaming-platform/services/user-service/internal/clients"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/grpc"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	shutdownTimeout = 5 * time.Second
)

// @title User Service API
// @version 1.0
// @description API для управления пользователями в стриминговой платформе
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token (JWT) для авторизации
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name refreshToken
// @description Refresh-токен в куках для аутентификации
func main() {
	logger := logging.InitLogger("user-service")

	cfg, err := config.LoadConfig("dev") // dev, prod, test
	if err != nil {
		logger.Error("❌ Couldn't load the configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("✅ Configuration loaded successfully")

	authClientAddr := fmt.Sprintf("%s:%d", cfg.Services.AuthService.Host, cfg.Services.AuthService.Port)
	clients, err := cl.NewClients(authClientAddr)
	if err != nil {
		logger.Error("❌ Clients initialization error", slog.String("error", err.Error()))
		os.Exit(1)
	} else {
		logger.Info("✅ Clients are initialized")
	}

	database, err := db.NewPostgresConnection(cfg.DBConfig)
	if err != nil {
		logger.Error("❌ Database connection error", slog.String("error", err.Error()))
		os.Exit(1)
		return
	}

	userRepo := repository.NewUserRepository(database)
	logger.Info("🔧 Repositories are initialized")

	userService := service.NewUserService(clients.Auth, userRepo, logger)
	logger.Info("🔧 Services are initialized")

	httpHandler := httpTransport.NewHandler(userService, logger)
	logger.Info("🔧 HTTP handlers are initialized")

	httpRouter := httpTransport.NewRouter(httpHandler, logger, cfg.CORS)
	logger.Info("🔧 HTTP Router is initialized")

	grpcHandler := grpcTransport.NewHandler(userService, logger)
	logger.Info("🔧 GRPC handlers are initialized")

	grpcServer := grpcTransport.NewGRPCServer(grpcHandler, logger)
	logger.Info("🔧 gRPC Server is initialized")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic in HTTP server", slog.Any("panic", r))
			}
		}()
		httpAddr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
		if err := httpRouter.Run(httpAddr); err != nil {
			logger.Error("❌ HTTP Server error", slog.String("error", err.Error()))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic in HTTP server", slog.Any("panic", r))
			}
		}()
		grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
		if err := grpcServer.Run(grpcAddr); err != nil {
			logger.Error("❌ gRPC Server error", slog.String("error", err.Error()))
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

	if err := httpRouter.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown HTTP server", slog.String("error", err.Error()))
	}

	wg.Wait()

	if err := database.Close(); err != nil {
		logger.Error("Couldn't close the database", slog.String("error", err.Error()))
	} else {
		logger.Info("✅ The database connection is closed")
	}

	if err := clients.Auth.Close(); err != nil {
		logger.Error("❌ Failed to close AuthClient connection", slog.String("error", err.Error()))
	} else {
		logger.Info("✅ AuthClient connection closed")
	}

	logger.Info("Server shut down gracefully")
}
