package repository

import (
	"context"

	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/entity"

	"github.com/google/uuid"
)

// ChatRepository определяет контракт для работы с данными чата
type ChatRepository interface {
	// Сообщения
	SaveMessage(ctx context.Context, msg *entity.ChatMessage) error
	GetMessages(ctx context.Context, streamID uuid.UUID, limit int) ([]*entity.ChatMessage, error)
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error

	// Комнаты
	CreateRoom(ctx context.Context, streamID uuid.UUID, title string) (*entity.ChatRoom, error)
	GetRoom(ctx context.Context, streamID uuid.UUID) (*entity.ChatRoom, error)
	CloseRoom(ctx context.Context, streamID uuid.UUID) error

	// Модерация
	BanUser(ctx context.Context, streamID, userID uuid.UUID) error
	UnbanUser(ctx context.Context, streamID, userID uuid.UUID) error
}
