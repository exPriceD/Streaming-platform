package model

import (
	"time"

	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/entity"
	"github.com/google/uuid"
)

type ChatRoom struct {
	ID          uuid.UUID `bson:"_id"`       // UUID комнаты
	StreamID    uuid.UUID `bson:"stream_id"` // streams.id
	StreamTitle string    `bson:"stream_title"`
	CreatedAt   time.Time `bson:"created_at"`
	IsActive    bool      `bson:"is_active"`
}

// ToEntity конвертирует в бизнес-сущность
func (cr *ChatRoom) ToEntity() *entity.ChatRoom {
	return &entity.ChatRoom{
		ID:       cr.ID,
		StreamID: cr.StreamID,
	}
}
