package service

import (
	"errors"
	"fmt"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	"github.com/google/uuid"
	"time"
)

var (
	ErrTokenExpired        = errors.New("the token expired")
	ErrTokenInvalid        = errors.New("invalid token")
	ErrRefreshTokenRevoked = errors.New("refresh token has been revoked")
)

type AuthService struct {
	tokenRepo  repository.TokenRepository
	jwtManager *token.JWTManager
}

func NewAuthService(tokenRepo repository.TokenRepository, jwtManager *token.JWTManager) *AuthService {
	service := &AuthService{
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}

	go service.startTokenCleanupRoutine()

	return service
}

func (s *AuthService) Authenticate(userID uuid.UUID) (string, string, int64, time.Time, error) {
	accessToken, refreshToken, expiresIn, expiresAt, err := s.GenerateTokens(userID)
	if err != nil {
		return "", "", 0, time.Time{}, err
	}

	return accessToken, refreshToken, expiresIn, expiresAt, nil
}

func (s *AuthService) RefreshTokens(refreshTokenStr string) (string, string, int64, time.Time, error) {
	refreshToken, err := s.tokenRepo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", 0, time.Time{}, err
	}

	if refreshToken.Revoked {
		return "", "", 0, time.Time{}, ErrRefreshTokenRevoked
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return "", "", 0, time.Time{}, ErrTokenExpired
	}

	accessToken, newRefreshToken, expiresIn, expiresAt, err := s.GenerateTokens(refreshToken.UserID)
	if err != nil {
		return "", "", 0, time.Time{}, err
	}

	err = s.tokenRepo.RevokeRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", 0, time.Time{}, err
	}
	return accessToken, newRefreshToken, expiresIn, expiresAt, nil
}

func (s *AuthService) GenerateTokens(userID uuid.UUID) (string, string, int64, time.Time, error) {
	accessToken, refreshToken, expiresIn, expiresAt, err := s.jwtManager.GenerateTokens(userID)
	if err != nil {
		return "", "", 0, time.Time{}, err
	}

	refreshTokenEntity := &entity.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenDuration),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.SaveRefreshToken(refreshTokenEntity); err != nil {
		return "", "", 0, time.Time{}, err
	}

	return accessToken, refreshToken, expiresIn, expiresAt, nil
}

func (s *AuthService) ValidateAccessToken(accessToken string) (uuid.UUID, error) {
	claims, err := s.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		if err.Error() == "token is expired" {
			return uuid.Nil, ErrTokenExpired
		}
		return uuid.Nil, ErrTokenInvalid
	}

	return claims.UserID, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	storedToken, err := s.tokenRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	if storedToken.Revoked {
		return nil
	}

	err = s.tokenRepo.RevokeRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) startTokenCleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.tokenRepo.DeleteExpiredRefreshTokens(); err != nil {
				fmt.Println("refresh_tokens cleanup error")
			}
		}
	}
}
