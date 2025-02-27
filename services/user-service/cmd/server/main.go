package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	database "github.com/exPriceD/Streaming-platform/pkg/db"
	logging "github.com/exPriceD/Streaming-platform/pkg/logger"
	_ "github.com/exPriceD/Streaming-platform/services/user-service/docs"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/clients"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/service"
	grpcTransport "github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/grpc"
	httpTransport "github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	logger      *slog.Logger
	cfg         *config.Config
	clients     *clients.Clients
	db          *sql.DB
	httpRouter  *httpTransport.Router
	grpcServer  *grpcTransport.Server
	shutdownWg  sync.WaitGroup
	shutdownCtx context.Context
	shutdownFn  context.CancelFunc
}

func NewApp(ctx context.Context) (*App, error) {
	logger := logging.InitLogger("user-service")

	cfg, err := config.LoadConfig("dev")
	if err != nil {
		logger.Error("Failed to load config", slog.String("error", err.Error()))
		return nil, fmt.Errorf("load config: %w", err)
	}
	logger.Info("Configuration loaded", slog.String("env", "dev"))

	clientsCfg := clients.Config{
		AuthServiceAddr: fmt.Sprintf("%s:%d", cfg.Services.AuthService.Host, cfg.Services.AuthService.Port),
		DialTimeout:     cfg.GRPC.DialTimeout,
	}

	cl, err := clients.NewClients(clientsCfg)
	if err != nil {
		logger.Error("Failed to initialize clients", slog.String("error", err.Error()))
		return nil, fmt.Errorf("initialize clients: %w", err)
	}
	logger.Info("Clients initialized")

	db, err := database.NewPostgresConnection(cfg.DBConfig)
	if err != nil {
		if shutdownErr := cl.Shutdown(ctx); shutdownErr != nil {
			logger.Error("Failed to shutdown clients during cleanup", slog.String("error", shutdownErr.Error()))
			return nil, fmt.Errorf("shutdown clients: %w; original error: %v", shutdownErr, err)
		}
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	logger.Info("Database connection established")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(cl.Auth, userRepo, logger)
	logger.Info("Repositories and services initialized")

	httpHandler := httpTransport.NewHandler(userService, logger)
	httpRouter := httpTransport.NewRouter(httpHandler, logger, cfg.CORS)
	logger.Info("HTTP router initialized")

	grpcHandler := grpcTransport.NewHandler(userService, logger)
	grpcServer := grpcTransport.NewGRPCServer(grpcHandler, logger)
	logger.Info("gRPC server initialized")

	shutdownCtx, shutdownFn := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)

	return &App{
		logger:      logger,
		cfg:         cfg,
		clients:     cl,
		db:          db,
		httpRouter:  httpRouter,
		grpcServer:  grpcServer,
		shutdownWg:  sync.WaitGroup{},
		shutdownCtx: shutdownCtx,
		shutdownFn:  shutdownFn,
	}, nil
}

func (a *App) Run() error {
	a.shutdownWg.Add(2)

	go a.runHTTPServer()
	go a.runGRPCServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	a.logger.Info("Received shutdown signal")
	return a.Shutdown()
}

func (a *App) runHTTPServer() {
	defer a.shutdownWg.Done()
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("Panic in HTTP server", slog.Any("panic", r))
		}
	}()

	httpAddr := fmt.Sprintf("%s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port)
	if err := a.httpRouter.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error("Failed to start HTTP server", slog.String("address", httpAddr), slog.String("error", err.Error()))
	}
}

func (a *App) runGRPCServer() {
	defer a.shutdownWg.Done()
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("Panic in gRPC server", slog.Any("panic", r))
		}
	}()

	grpcAddr := fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port)
	if err := a.grpcServer.Run(grpcAddr); err != nil {
		a.logger.Error("Failed to start gRPC server", slog.String("address", grpcAddr), slog.String("error", err.Error()))
	}
}

func (a *App) Shutdown() error {
	defer a.shutdownFn()

	errChan := make(chan error, 3)

	a.shutdownWg.Add(3)
	go func() {
		defer a.shutdownWg.Done()
		if err := a.httpRouter.Shutdown(a.shutdownCtx); err != nil {
			errChan <- fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}()
	go func() {
		defer a.shutdownWg.Done()
		if err := a.grpcServer.Shutdown(a.shutdownCtx); err != nil {
			errChan <- fmt.Errorf("shutdown gRPC server: %w", err)
		}
	}()
	go func() {
		defer a.shutdownWg.Done()
		if err := a.clients.Shutdown(a.shutdownCtx); err != nil {
			errChan <- fmt.Errorf("shutdown clients: %w", err)
		}
	}()

	go func() {
		a.shutdownWg.Wait()
		close(errChan)
	}()

	var shutdownErr error
	for err := range errChan {
		if shutdownErr == nil {
			shutdownErr = err
		} else {
			shutdownErr = fmt.Errorf("%v; %w", shutdownErr, err)
		}
		a.logger.Error("Shutdown error", slog.String("error", err.Error()))
	}

	if err := a.db.Close(); err != nil {
		a.logger.Error("Failed to close database", slog.String("error", err.Error()))
		if shutdownErr == nil {
			shutdownErr = fmt.Errorf("close database: %w", err)
		} else {
			shutdownErr = fmt.Errorf("%v; close database: %w", shutdownErr, err)
		}
	} else {
		a.logger.Info("Database connection closed")
	}

	if shutdownErr != nil {
		return shutdownErr
	}
	a.logger.Info("Server shut down gracefully")
	return nil
}

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
	app, err := NewApp(context.Background())
	if err != nil {
		slog.Error("Failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		app.logger.Error("Application run error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
