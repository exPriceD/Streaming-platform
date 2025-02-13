package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshTokenModel struct {
	ID        int       `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}
