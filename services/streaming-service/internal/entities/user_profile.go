package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID        uuid.UUID
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
