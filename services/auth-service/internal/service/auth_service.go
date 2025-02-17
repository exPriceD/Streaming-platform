package service

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	"github.com/google/uuid"
	"log/slog"
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
	log        *slog.Logger
}

func NewAuthService(tokenRepo repository.TokenRepository, jwtManager *token.JWTManager, log *slog.Logger) *AuthService {
	service := &AuthService{
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		log:        log,
	}

	go service.startTokenCleanupRoutine()

	return service
}

func (s *AuthService) Authenticate(userID uuid.UUID) (string, string, int64, time.Time, error) {
	accessToken, refreshToken, expiresIn, expiresAt, err := s.GenerateTokens(userID)
	if err != nil {
		s.log.Error("Failed to authenticate user", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	s.log.Info("User authenticated", slog.String("user_id", userID.String()))
	return accessToken, refreshToken, expiresIn, expiresAt, nil
}

func (s *AuthService) RefreshTokens(refreshTokenStr string) (string, string, int64, time.Time, error) {
	refreshToken, err := s.tokenRepo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		s.log.Warn("Failed to refresh tokens", slog.String("refresh_token", refreshTokenStr), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	if refreshToken.Revoked {
		s.log.Warn("Refresh token has been revoked", slog.String("user_id", refreshToken.UserID.String()))
		return "", "", 0, time.Time{}, ErrRefreshTokenRevoked
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		s.log.Warn("Refresh token expired", slog.String("user_id", refreshToken.UserID.String()))
		return "", "", 0, time.Time{}, ErrTokenExpired
	}

	accessToken, newRefreshToken, expiresIn, expiresAt, err := s.GenerateTokens(refreshToken.UserID)
	if err != nil {
		s.log.Error("Failed to generate new tokens", slog.String("user_id", refreshToken.UserID.String()), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	err = s.tokenRepo.RevokeRefreshToken(refreshTokenStr)
	if err != nil {
		s.log.Error("Failed to revoke old refresh token", slog.String("token", refreshTokenStr), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	s.log.Info("Tokens refreshed", slog.String("user_id", refreshToken.UserID.String()))
	return accessToken, newRefreshToken, expiresIn, expiresAt, nil
}

func (s *AuthService) GenerateTokens(userID uuid.UUID) (string, string, int64, time.Time, error) {
	accessToken, refreshToken, expiresIn, expiresAt, err := s.jwtManager.GenerateTokens(userID)
	if err != nil {
		s.log.Error("Failed to generate tokens", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
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
		s.log.Error("Failed to save refresh token", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	s.log.Info("Tokens generated", slog.String("user_id", userID.String()))
	return accessToken, refreshToken, expiresIn, expiresAt, nil
}

func (s *AuthService) ValidateToken(accessToken string) (uuid.UUID, error) {
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		if err.Error() == "token is expired" {
			s.log.Warn("Access token expired")
			return uuid.Nil, ErrTokenExpired
		}
		s.log.Warn("Invalid access token")
		return uuid.Nil, ErrTokenInvalid
	}

	s.log.Info("Access token validated", slog.String("user_id", claims.UserID.String()))
	return claims.UserID, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	storedToken, err := s.tokenRepo.GetRefreshToken(refreshToken)
	if err != nil {
		s.log.Warn("Failed to logout", slog.String("error", err.Error()))
		return err
	}

	if storedToken.Revoked {
		s.log.Warn("User already logged out", slog.String("user_id", storedToken.UserID.String()))
		return nil
	}

	err = s.tokenRepo.RevokeRefreshToken(refreshToken)
	if err != nil {
		s.log.Error("Failed to revoke refresh token", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("User logged out", slog.String("user_id", storedToken.UserID.String()))
	return nil
}

func (s *AuthService) startTokenCleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.tokenRepo.DeleteExpiredRefreshTokens(); err != nil {
				s.log.Error("Failed to clean up expired refresh tokens", slog.String("error", err.Error()))
			} else {
				s.log.Info("Expired refresh tokens cleaned up")
			}
		}
	}
}
