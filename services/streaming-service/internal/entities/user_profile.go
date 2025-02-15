package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                 uuid.UUID
	ChannelName        string
	ChannelDescription string
	StreamKey          string
	IsLive             bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
