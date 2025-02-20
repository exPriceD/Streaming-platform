package model

import (
	"time"

	"github.com/google/uuid"
)

// ChatBan хранит информацию о забаненных пользователях
type ChatBan struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	StreamID uuid.UUID `gorm:"type:uuid;index"`
	UserID   uuid.UUID `gorm:"type:uuid;index"`
	BannedAt time.Time
}
