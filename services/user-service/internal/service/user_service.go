package service

import (
	"context"
	"errors"
	"fmt"
	authProto "github.com/exPriceD/Streaming-platform/pkg/proto/v1/auth"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/clients"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"log/slog"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByID(ctx context.Context, userID string) (*entity.User, error)
}

type UserService struct {
	authClient *clients.AuthClient
	userRepo   UserRepository
	log        *slog.Logger
}

func NewUserService(authClient *clients.AuthClient, userRepo UserRepository, logger *slog.Logger) *UserService {
	return &UserService{authClient: authClient, userRepo: userRepo, log: logger}
}

func (s *UserService) RegisterUser(ctx context.Context, username, email, password, confirmPassword string, consent bool) (string, string, string, *entity.User, error) {
	user, err := entity.NewUser(username, email, password, confirmPassword, consent)
	if err != nil {
		return "", "", "", nil, err
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", "", "", nil, err
	}

	resp, err := s.authenticateUser(ctx, user.ID.String())
	if err != nil {
		return "", "", "", nil, err
	}

	return user.ID.String(), resp.AccessToken, resp.RefreshToken, user, nil
}

func (s *UserService) LoginUser(ctx context.Context, loginIdentifier, password string) (string, string, string, error) {
	var user *entity.User
	var err error

	if utils.IsEmail(loginIdentifier) {
		user, err = s.userRepo.GetUserByEmail(ctx, loginIdentifier)
	} else {
		user, err = s.userRepo.GetUserByUsername(ctx, loginIdentifier)
	}

	if err != nil {
		return "", "", "", err
	}

	if !user.CheckPassword(password) {
		return "", "", "", errors.New("incorrect password")
	}

	resp, err := s.authenticateUser(ctx, user.ID.String())
	if err != nil {
		return "", "", "", err
	}
	return user.ID.String(), resp.AccessToken, resp.RefreshToken, nil
}

func (s *UserService) authenticateUser(ctx context.Context, userID string) (*authProto.AuthenticateResponse, error) {
	req := &authProto.AuthenticateRequest{UserId: userID}
	resp, err := s.authClient.Authenticate(ctx, req)
	if err != nil || resp.Error != nil {
		return nil, err
	}
	return resp, nil
}

func (s *UserService) ValidateToken(ctx context.Context, accessToken string) (bool, error) {
	validateResp, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil || validateResp.Error != nil {
		return false, err
	}

	if validateResp.Error != nil {
		return false, fmt.Errorf("invalid token: %v", validateResp.Error)
	}

	if validateResp.Valid {
		return true, nil
	}

	return false, nil
}
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	refreshResp, err := s.authClient.RefreshToken(ctx, refreshToken)
	if err != nil || refreshResp.Error != nil {
		return "", "", err
	}

	if refreshResp.Error != nil {
		return "", "", fmt.Errorf("invalid refresh token: %v", refreshResp.Error)
	}
	return refreshResp.AccessToken, refreshResp.RefreshToken, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) (bool, error) {
	logoutResp, err := s.authClient.Logout(ctx, refreshToken)
	if err != nil {
		return false, err
	}
	return logoutResp.Success, err
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}
