package repository

import "github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"

type TokenRepository interface {
	SaveRefreshToken(token *entities.RefreshToken) error
	GetRefreshToken(tokenStr string) (*entities.RefreshToken, error)
	RevokeRefreshToken(tokenStr string) error
	DeleteExpiredRefreshTokens() error
}

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
}
