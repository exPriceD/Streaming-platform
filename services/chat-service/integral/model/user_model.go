package model

import (
	"time"

	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/entity"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `bson:"user_id"`  // Первичный ключ
	Username     string    `bson:"username"` // Уникальное имя
	Email        string    `bson:"email"`    // Уникальный email
	PasswordHash string    `bson:"password_hash"`
	AvatarURL    string    `bson:"avatar_url"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

// ToEntity конвертирует в сущность для бизнес-логики
func (u *User) ToEntity() *entity.UserConnection {
	return &entity.UserConnection{
		UserID:   u.ID,
		Username: u.Username,
	}
}
