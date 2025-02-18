package entity

import (
	"sync"

	"github.com/google/uuid"
)

// ChatRoom управляет подключениями пользователей для стрима
type ChatRoom struct {
	ID          uuid.UUID                     // Уникальный ID комнаты
	StreamID    uuid.UUID                     // Ссылка на streams.id
	Connections map[uuid.UUID]*UserConnection // Активные подключения
	mu          sync.RWMutex                  // Для конкурентного доступа
}

// NewChatRoom создает новую комнату для стрима
func NewChatRoom(streamID uuid.UUID) *ChatRoom {
	return &ChatRoom{
		ID:          uuid.New(),
		StreamID:    streamID,
		Connections: make(map[uuid.UUID]*UserConnection),
	}
}

// AddConnection добавляет пользователя в комнату
func (cr *ChatRoom) AddConnection(uc *UserConnection) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.Connections[uc.UserID] = uc
	uc.StreamID = cr.StreamID
}

// RemoveConnection удаляет подключение
func (cr *ChatRoom) RemoveConnection(userID uuid.UUID) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	if conn, exists := cr.Connections[userID]; exists {
		conn.Close()
		delete(cr.Connections, userID)
	}
}

// Broadcast отправляет сообщение всем участникам
func (cr *ChatRoom) Broadcast(message []byte) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	for _, conn := range cr.Connections {
		select {
		case conn.SendChan <- message:
		default:
			conn.Close()
		}
	}
}
