package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	customErrors "github.com/exPriceD/Streaming-platform/services/user-service/internal/errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"log/slog"
)

type AuthClient interface {
	Authenticate(ctx context.Context, userId string) (string, string, error)
	ValidateToken(ctx context.Context, accessToken string) (bool, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) (bool, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, userId string) (*entity.User, error)
}

type UserService struct {
	authClient AuthClient
	userRepo   UserRepository
	log        *slog.Logger
}

func NewUserService(authClient AuthClient, userRepo UserRepository, logger *slog.Logger) *UserService {
	return &UserService{authClient: authClient, userRepo: userRepo, log: logger}
}

func (s *UserService) RegisterUser(ctx context.Context, username, email, password, confirmPassword string, consent bool) (string, string, string, error) {
	if username == "" || email == "" || password == "" {
		return "", "", "", customErrors.ErrInvalidInput
	}

	user, err := entity.NewUser(username, email, password, confirmPassword, consent)
	if err != nil {
		return "", "", "", err
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return "", "", "", err
	}

	accessToken, refreshToken, err := s.authClient.Authenticate(ctx, user.Id.String())
	if err != nil {
		return "", "", "", fmt.Errorf("authentication error: %w", err)
	}

	return user.Id.String(), accessToken, refreshToken, nil
}

func (s *UserService) LoginUser(ctx context.Context, loginIdentifier, password string) (string, string, string, error) {
	if loginIdentifier == "" || password == "" {
		return "", "", "", customErrors.ErrInvalidInput
	}

	var user *entity.User
	var err error

	if utils.IsEmail(loginIdentifier) {
		user, err = s.userRepo.GetUserByEmail(ctx, loginIdentifier)
	} else {
		user, err = s.userRepo.GetUserByUsername(ctx, loginIdentifier)
	}

	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			return "", "", "", err
		}
		return "", "", "", fmt.Errorf("login error: %w", err)
	}

	if !user.CheckPassword(password) {
		return "", "", "", customErrors.ErrUnauthorized
	}

	accessToken, refreshToken, err := s.authClient.Authenticate(ctx, user.Id.String())
	if err != nil {
		return "", "", "", fmt.Errorf("authentication error: %s", err)
	}

	return user.Id.String(), accessToken, refreshToken, nil
}

func (s *UserService) ValidateToken(ctx context.Context, accessToken string) (bool, string, error) {
	if accessToken == "" {
		return false, "", customErrors.ErrInvalidInput
	}

	valid, userId, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil {
		return false, "", fmt.Errorf("token validation error: %s", err)
	}
	return valid, userId, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", customErrors.ErrInvalidInput
	}

	accessToken, newRefreshToken, err := s.authClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("refresh token error: %s", err)
	}

	return accessToken, newRefreshToken, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) (bool, error) {
	if refreshToken == "" {
		return false, customErrors.ErrInvalidInput
	}

	success, err := s.authClient.Logout(ctx, refreshToken)
	if err != nil {
		return false, fmt.Errorf("logout error: %w", err)
	}

	return success, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userId string) (*entity.User, error) {
	if userId == "" {
		return nil, customErrors.ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get user error: %w", err)
	}

	return user, nil
}
