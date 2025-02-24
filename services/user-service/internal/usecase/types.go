package usecase

import (
	"context"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/dto"
)

// User представляет интерфейс пользователя в доменной модели.
type User interface {
	ID() string
	CheckPassword(password string) bool
	ToDTO() *dto.User
}

type AuthClient interface {
	Authenticate(ctx context.Context, userId string) (string, string, error)
	ValidateToken(ctx context.Context, accessToken string) (bool, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) (bool, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *dto.User) error
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	GetUserByUsername(ctx context.Context, username string) (*dto.User, error)
	GetUserByID(ctx context.Context, userId string) (*dto.User, error)
}

// RegisterResponse представляет результат регистрации пользователя.
type RegisterResponse struct {
	UserId       string
	AccessToken  string
	RefreshToken string
}

// LoginResponse представляет результат входа пользователя.
type LoginResponse struct {
	UserId       string
	AccessToken  string
	RefreshToken string
}

// RefreshTokenResponse представляет результат обновления токена.
type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}
