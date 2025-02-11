package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Stream struct {
	ID        uuid.UUID
	Title     string
	UserID    uuid.UUID
	Status    string // active, ended, scheduled
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Возможные статусы стрима
const (
	StreamStatusActive    = "active"
	StreamStatusEnded     = "ended"
	StreamStatusScheduled = "scheduled"
)

// NewStream создаёт новый объект стрима с проверками
func NewStream(title string, userID uuid.UUID) (*Stream, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &Stream{
		ID:        uuid.New(),
		Title:     title,
		UserID:    userID,
		Status:    StreamStatusScheduled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateStatus обновляет статус стрима
func (s *Stream) UpdateStatus(status string) error {
	if status != StreamStatusActive && status != StreamStatusEnded && status != StreamStatusScheduled {
		return errors.New("invalid stream status")
	}
	s.Status = status
	s.UpdatedAt = time.Now()
	return nil
}
