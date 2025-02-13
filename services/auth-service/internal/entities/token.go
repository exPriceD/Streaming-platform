package entities

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        int
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}
