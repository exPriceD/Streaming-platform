package service

import (
	"context"
	"errors"
	"fmt"
	authProto "github.com/exPriceD/Streaming-platform/pkg/proto/v1/auth"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/clients"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	customErrors "github.com/exPriceD/Streaming-platform/services/user-service/internal/errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"log/slog"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID string) (*entity.User, error)
}

type UserService struct {
	authClient *clients.AuthClient
	userRepo   UserRepository
	log        *slog.Logger
}

func NewUserService(authClient *clients.AuthClient, userRepo UserRepository, logger *slog.Logger) *UserService {
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

	resp, err := s.authenticateUser(ctx, user.ID.String())
	if err != nil {
		return "", "", "", fmt.Errorf("authentication error: %w", err)
	}

	return user.ID.String(), resp.AccessToken, resp.RefreshToken, nil
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

	resp, err := s.authenticateUser(ctx, user.ID.String())
	if err != nil {
		return "", "", "", fmt.Errorf("authentication error: %w", err)
	}

	return user.ID.String(), resp.AccessToken, resp.RefreshToken, nil
}

func (s *UserService) authenticateUser(ctx context.Context, userId string) (*authProto.AuthenticateResponse, error) {
	req := &authProto.AuthenticateRequest{UserId: userId}
	resp, err := s.authClient.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("authentication failed: %s", resp.Error.Message)
	}

	return resp, nil
}

func (s *UserService) ValidateToken(ctx context.Context, accessToken string) (bool, string, error) {
	if accessToken == "" {
		return false, "", customErrors.ErrInvalidInput
	}

	resp, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil {
		return false, "", fmt.Errorf("token validation error: %w", err)
	}
	if resp.Error != nil {
		return false, "", nil
	}

	return resp.Valid, resp.UserId, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", customErrors.ErrInvalidInput
	}

	resp, err := s.authClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("refresh token error: %w", err)
	}
	if resp.Error != nil {
		return "", "", fmt.Errorf("refresh token rejected: %s", resp.Error.Message)
	}

	return resp.AccessToken, resp.RefreshToken, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) (bool, error) {
	if refreshToken == "" {
		return false, customErrors.ErrInvalidInput
	}

	resp, err := s.authClient.Logout(ctx, refreshToken)
	if err != nil {
		return false, fmt.Errorf("logout error: %w", err)
	}
	if resp.Error != nil {
		return false, nil
	}

	return resp.Success, err
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
