package main

import (
	"database/sql"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/exPriceD/Streaming-platform/pkg/logger"
	cl "github.com/exPriceD/Streaming-platform/services/user-service/internal/clients"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/config"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/service"
	router "github.com/exPriceD/Streaming-platform/services/user-service/internal/transport/http"
	"log/slog"
)

func main() {
	log := logger.InitLogger("user-service")
	clients, err := cl.NewClients("localhost:50051")
	if err != nil {
		log.Error("❌ Clients initialization error", slog.String("error", err.Error()))
	} else {
		log.Info("✅ Clients are initialized")
	}

	cfg, err := config.LoadConfig("dev") // dev, prod, test
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
	log.Info("🔧 Repositories are initialized")

	userService := service.NewUserService(clients.Auth, userRepo, log)
	log.Info("🔧 Services are initialized")

	handler := router.NewHandler(userService)
	log.Info("🔧 Handlers are initialized")

	r := router.NewRouter(handler, log, userService)

	if err := r.Run(":8080"); err != nil {
		log.Error("❌ Server error", slog.String("error", err.Error()))
	}
}
