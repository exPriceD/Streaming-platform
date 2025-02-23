package model

import (
	"time"

	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/entity"
	"github.com/google/uuid"
)

type ChatMessage struct {
	ID        uuid.UUID `bson:"_id"`       // ObjectID для MongoDB
	StreamID  uuid.UUID `bson:"stream_id"` // streams.id
	UserID    uuid.UUID `bson:"user_id"`   // users.id
	Username  string    `bson:"username"`  // users.username
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"sent_at"`
	IsDeleted bool      `bson:"is_deleted"`
	ModReason string    `bson:"mod_reason,omitempty"`
}

// ToEntity конвертирует в бизнес-сущность
func (cm *ChatMessage) ToEntity() *entity.ChatMessage {
	return &entity.ChatMessage{
		ID:        cm.ID,
		StreamID:  cm.StreamID,
		UserID:    cm.UserID,
		Username:  cm.Username,
		Content:   cm.Content,
		Timestamp: cm.Timestamp,
		IsDeleted: cm.IsDeleted,
	}
}
