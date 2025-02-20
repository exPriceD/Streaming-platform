package main

import (
	"log"

	"github.com/exPriceD/Streaming-platform/config"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/handler"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/service"
	auth "github.com/exPriceD/Streaming-platform/services/streaming-service/internal/token"
	"github.com/labstack/echo/v4"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadStreamingConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаем PostgreSQL
	db, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Создаём JWT менеджер
	jwtManager := auth.NewJWTManager(cfg)

	// Создаём Echo сервер
	e := echo.New()

	// Инициализируем сервис
	streamRepo := repository.NewStreamRepository(db)
	userRepo := repository.NewUserProfileRepository(db)
	ffmpegService := service.NewFFmpegService()
	baseStreamURL := "http://localhost:8080/hls/"

	streamService := service.NewStreamService(streamRepo, ffmpegService, userRepo, baseStreamURL)

	// Регистрируем обработчики
	handler.NewStreamHandler(e, streamService, jwtManager)

	// Запускаем сервер
	address := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Streaming Service запущен на %s", address)
	if err := e.Start(address); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
