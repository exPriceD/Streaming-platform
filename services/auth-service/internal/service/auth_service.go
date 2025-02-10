package service

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/token"
	"time"
)

type AuthService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager *token.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwtManager *token.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
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

func (s *AuthService) generateTokens(userID string) (string, string, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshTokenStr, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	refreshToken := &entities.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenDuration),
		Revoked:   false,
	}

	if err := s.tokenRepo.SaveRefreshToken(refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenStr, nil
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (*token.UserClaims, error) {
	return s.jwtManager.ValidateToken(tokenStr)
}

func (s *AuthService) RefreshTokens(refreshTokenStr string) (string, error) {
	refreshToken, err := s.tokenRepo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return "", err
	}

	return s.jwtManager.GenerateAccessToken(refreshToken.UserID)
}
