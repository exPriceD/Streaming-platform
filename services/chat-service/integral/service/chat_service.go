package service

import (
	"context"

	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/entity"
	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/repository"
	"github.com/google/uuid"
)

// ChatService реализует бизнес-логику чата
type ChatService struct {
	repo repository.ChatRepository
}

// NewChatService создает новый сервис
func NewChatService(repo repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

// GetMessages получает сообщения
func (s *ChatService) GetMessages(ctx context.Context, streamID uuid.UUID, limit int) ([]*entity.ChatMessage, error) {
	return s.repo.GetMessages(ctx, streamID, limit)
}

// SendMessage сохраняет сообщение
func (s *ChatService) SendMessage(ctx context.Context, msg *entity.ChatMessage) error {
	return s.repo.SaveMessage(ctx, msg)
}
