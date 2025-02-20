package repository

import (
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"

	"github.com/google/uuid"
)

// StreamRepositoryInterface определяет методы работы с потоками в БД
type StreamRepositoryInterface interface {
	CreateStream(stream models.Stream) error
	GetStreamByID(id uuid.UUID) (*models.Stream, error)
	UpdateStream(stream models.Stream) error
	DeleteStream(id uuid.UUID) error
}
