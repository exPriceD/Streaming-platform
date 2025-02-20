package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/service"
	"github.com/google/uuid"
)

// ChatHandler обрабатывает HTTP-запросы
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler создает новый обработчик
func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// GetMessages обрабатывает запрос на получение сообщений
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	streamID, err := uuid.Parse(r.URL.Query().Get("stream_id"))
	if err != nil {
		http.Error(w, "Некорректный stream_id", http.StatusBadRequest)
		return
	}

	messages, err := h.chatService.GetMessages(context.Background(), streamID, 20)
	if err != nil {
		http.Error(w, "Ошибка получения сообщений", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
