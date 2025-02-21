package repository

import (
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"
	"github.com/google/uuid"
)

type UserProfileRepositoryInterface interface {
	CreateUserProfile(profile models.UserProfile) error
	GetUserProfileByID(userID uuid.UUID) (*models.UserProfile, error)
	UpdateUserProfile(profile models.UserProfile) error
	UpdateLiveStatus(userID uuid.UUID, isLive bool) error
	SaveStreamKey(streamID uuid.UUID, streamKey string) error
	GetStreamKey(streamID uuid.UUID) (string, error)
	UpdateStreamKey(userID string, newStreamKey string) error
}
