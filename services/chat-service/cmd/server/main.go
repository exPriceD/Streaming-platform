package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exPriceD/Streaming-platform/config"
	"github.com/exPriceD/Streaming-platform/pkg/logger"
	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/handler"
	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/service"
	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Инициализация логгера
	log := logger.InitLogger("chat-service")
	log.Info("Starting chat service...")

	// Загрузка конфигурации
	cfg, err := config.LoadChatConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Подключение к PostgreSQL
	pgDB, err := connectPostgres(cfg.DB)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	log.Info("Successfully connected to PostgreSQL")

	// Подключение к MongoDB
	mongoDB, err := connectMongo(cfg.Mongo)
	if err != nil {
		log.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	log.Info("Successfully connected to MongoDB")

	// Подключение к Redis
	redisClient := connectRedis(cfg.Redis)
	log.Info("Successfully connected to Redis")

	// Инициализация репозитория
	repo := repository.NewChatRepository(mongoDB, pgDB)

	// Инициализация сервиса
	chatService := service.NewChatService(repo)

	// Инициализация WebSocket сервера
	wsServer := websocket.NewChatServer(cfg.WebSocket.JWTSecret, redisClient)

	// Инициализация HTTP обработчиков
	chatHandler := handler.NewChatHandler(chatService)

	// Настройка маршрутов HTTP
	http.HandleFunc("/messages", chatHandler.GetMessages)
	http.HandleFunc("/ws", wsServer.HandleConnection)

	// Запуск HTTP сервера
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: nil,
	}

	go func() {
		log.Info("Starting HTTP server", "host", cfg.Server.Host, "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start HTTP server", "error", err)
			os.Exit(1)
		}
	}()

	// Запуск рассылки сообщений через WebSocket
	go wsServer.StartBroadcast(context.Background())

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", "error", err)
	}

	log.Info("Server exited properly")
}

// connectPostgres подключается к PostgreSQL
func connectPostgres(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	return db, nil
}

// connectMongo подключается к MongoDB
func connectMongo(cfg config.MongoConfig) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(cfg.Database), nil
}

// connectRedis подключается к Redis
func connectRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
