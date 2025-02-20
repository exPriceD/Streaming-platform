package entity

import (
	"time"

	"github.com/google/uuid"
)

// ChatMessage представляет сообщение в чате стрима
type ChatMessage struct {
	ID        uuid.UUID `json:"id"`         // Уникальный ID сообщения
	StreamID  uuid.UUID `json:"stream_id"`  // Ссылка на streams.id
	UserID    uuid.UUID `json:"user_id"`    // Ссылка на users.id
	Username  string    `json:"username"`   // Дублирование из users.username
	Content   string    `json:"content"`    // Текст сообщения
	Timestamp time.Time `json:"timestamp"`  // Время отправки
	IsDeleted bool      `json:"is_deleted"` // Флаг удаления
}

// NewChatMessage создает новое сообщение
func NewChatMessage(streamID, userID uuid.UUID, username, content string) *ChatMessage {
	return &ChatMessage{
		ID:        uuid.New(),
		StreamID:  streamID,
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now().UTC(),
	}
}

// Validate проверяет валидность сообщения
func (m *ChatMessage) Validate() bool {
	return len(m.Content) > 0 && len(m.Username) > 0 && m.StreamID != uuid.Nil
}
