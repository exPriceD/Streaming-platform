package main

import (
	"context"
	_ "github.com/exPriceD/Streaming-platform/services/user-service/docs"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/app"
	"github.com/labstack/gommon/log"
	"log/slog"
	"os"
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
	application, err := app.NewFromEnv(context.Background())
	if err != nil {
		log.Printf("Failed to create application: %v", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		application.Logger().Error("Application run error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
