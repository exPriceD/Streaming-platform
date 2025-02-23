package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/entity"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ChatServer управляет подключениями пользователей
type ChatServer struct {
	clients   map[*websocket.Conn]uuid.UUID
	mu        sync.Mutex
	broadcast chan *entity.ChatMessage
}

// NewChatServer создает новый WebSocket-сервер
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[*websocket.Conn]uuid.UUID),
		broadcast: make(chan *entity.ChatMessage),
	}
}

// HandleConnection обрабатывает новое подключение
func (s *ChatServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка WebSocket:", err)
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.clients[conn] = uuid.New()
	s.mu.Unlock()

	for {
		var msg entity.ChatMessage
		if err := conn.ReadJSON(&msg); err != nil {
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			break
		}

		s.broadcast <- &msg
	}
}

// StartBroadcast запускает рассылку сообщений
func (s *ChatServer) StartBroadcast(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.broadcast:
			data, _ := json.Marshal(msg)
			s.mu.Lock()
			for client := range s.clients {
				client.WriteMessage(websocket.TextMessage, data)
			}
			s.mu.Unlock()
		}
	}
}
