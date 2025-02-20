package websocket

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/exPriceD/Streaming-platform/pkg/logger"
	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	log = logger.InitLogger("websocket")
)

// ChatServer управляет подключениями пользователей
type ChatServer struct {
	clients   map[uuid.UUID]*websocket.Conn // Хранение соединений по ID пользователя
	mu        sync.Mutex
	broadcast chan *entity.ChatMessage
	jwtSecret string
}

// NewChatServer создает новый WebSocket-сервер
func NewChatServer(jwtSecret string) *ChatServer {
	return &ChatServer{
		clients:   make(map[uuid.UUID]*websocket.Conn),
		broadcast: make(chan *entity.ChatMessage),
		jwtSecret: jwtSecret,
	}
}

// HandleConnection обрабатывает новое подключение
func (s *ChatServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Аутентификация через JWT
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		log.Warn("Connection attempt without token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Валидация токена
	claims, err := s.validateToken(tokenString)
	if err != nil {
		log.Warn("Invalid token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Обновление соединения до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	// Сохранение соединения
	userID := uuid.MustParse(claims["user_id"].(string))
	s.mu.Lock()
	s.clients[userID] = conn
	s.mu.Unlock()

	log.Info("New WebSocket connection", "user_id", userID)

	// Обработка входящих сообщений
	for {
		var msg entity.ChatMessage
		if err := conn.ReadJSON(&msg); err != nil {
			s.mu.Lock()
			delete(s.clients, userID)
			s.mu.Unlock()
			log.Info("WebSocket connection closed", "user_id", userID)
			break
		}

		// Валидация и обработка сообщения
		if err := s.processMessage(userID, &msg); err != nil {
			log.Warn("Message processing failed", "error", err)
			continue
		}

		s.broadcast <- &msg
	}
}

// validateToken проверяет JWT токен
func (s *ChatServer) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// processMessage обрабатывает входящее сообщение
func (s *ChatServer) processMessage(userID uuid.UUID, msg *entity.ChatMessage) error {
	// Установка идентификатора пользователя
	msg.UserID = userID
	msg.Timestamp = time.Now().UTC()

	// Валидация сообщения
	if !msg.Validate() {
		return errors.New("invalid message format")
	}

	// Проверка на спам
	if s.isSpam(userID) {
		return errors.New("spam detected")
	}

	return nil
}

// isSpam проверяет, не отправляет ли пользователь слишком много сообщений
func (s *ChatServer) isSpam(userID uuid.UUID) bool {
	// Реализация защиты от спама
	return false
}
