package repository

import "github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
}
