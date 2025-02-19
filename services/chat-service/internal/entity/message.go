package entity

import (
	"time"

	"github.com/google/uuid"
)

// ChatMessage представляет сообщение в чате стрима
type ChatMessage struct {
	ID        uuid.UUID // Уникальный ID сообщения
	StreamID  uuid.UUID // Ссылка на streams.id
	UserID    uuid.UUID // Ссылка на users.id
	Username  string    // Дублирование из users.username
	Content   string    // Текст сообщения
	Timestamp time.Time // Время отправки
	IsDeleted bool      // Флаг удаления
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
