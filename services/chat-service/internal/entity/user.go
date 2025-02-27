package entity

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// UserConnection представляет активное WebSocket-подключение пользователя
type UserConnection struct {
	ID       uuid.UUID       // Уникальный ID подключения
	UserID   uuid.UUID       // Ссылка на users.id
	Username string          // Дублирование из таблицы users
	Conn     *websocket.Conn // WebSocket соединение
	SendChan chan []byte     // Канал для исходящих сообщений
	StreamID uuid.UUID       // Идентификатор текущего стрима
	mu       sync.Mutex      // Для потокобезопасной работы
}

// NewUserConnection создает новое подключение пользователя
func NewUserConnection(userID uuid.UUID, username string, conn *websocket.Conn) *UserConnection {
	return &UserConnection{
		ID:       uuid.New(),
		UserID:   userID,
		Username: username,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
	}
}

// SendMessage безопасно отправляет сообщение
func (uc *UserConnection) SendMessage(msg []byte) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.Conn.WriteMessage(websocket.TextMessage, msg)
}

// Close аккуратно закрывает соединение
func (uc *UserConnection) Close() {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	close(uc.SendChan)
	uc.Conn.Close()
}
