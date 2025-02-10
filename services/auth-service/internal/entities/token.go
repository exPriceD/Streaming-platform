package entities

import "time"

type RefreshToken struct {
	ID        int
	UserID    string
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}
