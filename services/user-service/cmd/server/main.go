package main

import (
	"database/sql"
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	logging "github.com/exPriceD/Streaming-platform/pkg/logger"
	cl "github.com/exPriceD/Streaming-platform/services/user-service/internal/clients"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/service"
	router "github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http"
	"log/slog"
)

func main() {
	logger := logging.InitLogger("user-service")

	cfg, err := config.LoadConfig("dev") // dev, prod, test
	if err != nil {
		logger.Error("❌ Couldn't load the configuration", slog.String("error", err.Error()))
	}
	logger.Info("✅ Configuration loaded successfully")

	authClientAddr := fmt.Sprintf("%s:%d", cfg.Services.AuthService.Host, cfg.Services.AuthService.Port)
	clients, err := cl.NewClients(authClientAddr)
	if err != nil {
		logger.Error("❌ Clients initialization error", slog.String("error", err.Error()))
	} else {
		logger.Info("✅ Clients are initialized")
	}

	database, err := db.NewPostgresConnection(cfg.DBConfig)
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

	userRepo := repository.NewUserRepository(database)
	logger.Info("🔧 Repositories are initialized")

	userService := service.NewUserService(clients.Auth, userRepo, logger)
	logger.Info("🔧 Services are initialized")

	handler := router.NewHandler(userService, logger)
	logger.Info("🔧 Handlers are initialized")

	r := router.NewRouter(handler)

	httpServerAddr := fmt.Sprintf("%s:%d", cfg.Services.AuthService.Host, cfg.Services.AuthService.Port)
	if err := r.Run(httpServerAddr); err != nil {
		logger.Error("❌ Server error", slog.String("error", err.Error()))
	}
}
