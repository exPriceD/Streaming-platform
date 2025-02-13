package service

import (
	"errors"
	"fmt"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	"time"
)

var (
	ErrTokenExpired = errors.New("the token expired")
	ErrTokenInvalid = errors.New("invalid token")
)

type AuthService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager *token.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwtManager *token.JWTManager) *AuthService {
	service := &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}

	go service.startTokenCleanupRoutine()

	return service
}

func (s *AuthService) Register(username, email, password string, consent bool) (*entities.User, string, string, error) {
	user, err := entities.NewUser(username, email, password, consent)
	if err != nil {
		return nil, "", "", err
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, err
}

func (s *AuthService) Authenticate(identifier, password string, isEmail bool) (*entities.User, string, string, error) {
	var user *entities.User
	var err error

	if isEmail {
		user, err = s.userRepo.GetUserByEmail(identifier)
	} else {
		user, err = s.userRepo.GetUserByUsername(identifier)
	}

	if err != nil {
		return nil, "", "", errors.New("the user was not found")
	}

	if !user.CheckPassword(password) {
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshTokenStr string) (string, string, error) {
	refreshToken, err := s.tokenRepo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", err
	}

	if refreshToken.Revoked {
		return "", "", errors.New("refresh token has been revoked")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return "", "", errors.New("refresh token expired")
	}

	accessToken, newRefreshToken, err := s.generateTokens(refreshToken.UserID)
	if err != nil {
		return "", "", err
	}

	_ = s.tokenRepo.RevokeRefreshToken(refreshTokenStr)

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) generateTokens(userID string) (string, string, error) {
	accessToken, refreshToken, err := s.jwtManager.GenerateTokens(userID)
	if err != nil {
		return "", "", err
	}

	refreshTokenEntity := &entities.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenDuration),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.SaveRefreshToken(refreshTokenEntity); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateAccessToken(accessToken string) (string, error) {
	claims, err := s.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		if err.Error() == "token is expired" {
			return "", ErrTokenExpired
		}
		return "", ErrTokenInvalid
	}

	return claims.UserID, nil
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
