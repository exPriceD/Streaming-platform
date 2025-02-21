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
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
}

type UserService struct {
	authClient *clients.AuthClient
	userRepo   UserRepository
	log        *slog.Logger
}

func NewUserService(authClient *clients.AuthClient, userRepo UserRepository, logger *slog.Logger) *UserService {
	return &UserService{authClient: authClient, userRepo: userRepo, log: logger}
}

func (s *UserService) RegisterUser(username, email, password, confirmPassword string, consent bool) (string, string, string, *entity.User, error) {
	user, err := entity.NewUser(username, email, password, confirmPassword, consent)
	if err != nil {
		return "", "", "", nil, err
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return "", "", "", nil, err
	}

	resp, err := s.authenticateUser(user.ID.String())
	if err != nil {
		return "", "", "", nil, err
	}

	return user.ID.String(), resp.AccessToken, resp.RefreshToken, user, nil
}

func (s *UserService) LoginUser(loginIdentifier, password string) (string, string, string, error) {
	var user *entity.User
	var err error

	if utils.IsEmail(loginIdentifier) {
		user, err = s.userRepo.GetUserByEmail(loginIdentifier)
	} else {
		user, err = s.userRepo.GetUserByUsername(loginIdentifier)
	}

	if err != nil {
		return "", "", "", err
	}

	if !user.CheckPassword(password) {
		return "", "", "", errors.New("incorrect password")
	}

	resp, err := s.authenticateUser(user.ID.String())
	if err != nil {
		return "", "", "", err
	}
	return user.ID.String(), resp.AccessToken, resp.RefreshToken, nil
}

func (s *UserService) authenticateUser(userID string) (*authProto.AuthenticateResponse, error) {
	req := &authProto.AuthenticateRequest{UserId: userID}
	resp, err := s.authClient.Authenticate(context.Background(), req)
	if err != nil || resp.Error != nil {
		return nil, err
	}
	return resp, nil
}

func (s *UserService) ValidateToken(accessToken string) (bool, error) {
	validateResp, err := s.authClient.ValidateToken(context.Background(), accessToken)
	if err != nil || validateResp.Error != nil {
		return false, err
	}

	if validateResp.Error != nil {
		return false, fmt.Errorf("invalid token: %w", validateResp.Error)
	}

	if validateResp.Valid {
		return true, nil
	}

	return false, nil
}
func (s *UserService) RefreshToken(refreshToken string) (string, string, error) {
	refreshResp, err := s.authClient.RefreshToken(context.Background(), refreshToken)
	if err != nil || refreshResp.Error != nil {
		return "", "", err
	}

	if refreshResp.Error != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", refreshResp.Error)
	}
	return refreshResp.AccessToken, refreshResp.RefreshToken, nil
}

func (s *UserService) Logout(refreshToken string) (bool, error) {
	logoutResp, err := s.authClient.Logout(context.Background(), refreshToken)
	if err != nil {
		return false, err
	}
	return logoutResp.Success, err
}
func (s *UserService) GetUserByID(userID string) (*entity.User, error) {
	return nil, nil
}
