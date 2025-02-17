package repository

import "github.com/exPriceD/Streaming-platform/services/auth-service/internal/entity"

type TokenRepository interface {
	SaveRefreshToken(token *entity.RefreshToken) error
	GetRefreshToken(tokenStr string) (*entity.RefreshToken, error)
	RevokeRefreshToken(tokenStr string) error
	DeleteExpiredRefreshTokens() error
}
