package repository

import (
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"
	"github.com/google/uuid"
)

type UserProfileRepositoryInterface interface {
	CreateUserProfile(profile models.UserProfile) error
	GetUserProfileByID(id uuid.UUID) (*models.UserProfile, error)
	UpdateUserProfile(profile models.UserProfile) error
	UpdateLiveStatus(userID uuid.UUID, isLive bool) error
}
