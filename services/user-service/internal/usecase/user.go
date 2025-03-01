package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/domain"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/dto"
	customErrors "github.com/exPriceD/Streaming-platform/services/user-service/internal/errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/validation"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type UserUsecase struct {
	authClient AuthClient
	userRepo   UserRepository
	log        *slog.Logger
}

func NewUserUsecase(authClient AuthClient, userRepo UserRepository, logger *slog.Logger) *UserUsecase {
	return &UserUsecase{authClient: authClient, userRepo: userRepo, log: logger}
}

func (s *UserUsecase) RegisterUser(ctx context.Context, username, email, password, confirmPassword string, consent bool) (*RegisterResponse, error) {
	if username == "" || email == "" || password == "" {
		return nil, customErrors.ErrInvalidInput
	}

	if password != confirmPassword {
		return nil, domain.ErrPasswordsMismatch
	}
	if err := validation.ValidatePassword(password); err != nil {
		return nil, domain.WrapValidationError("password", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.WrapValidationError("password", domain.ErrHashingPassword)
	}

	cfg := domain.DefaultConfig()
	userDomain, err := domain.NewUser(username, email, string(hashedPassword), consent, cfg)
	if err != nil {
		return nil, err
	}

	userDTO := userDomain.ToDTO()
	if err := s.userRepo.CreateUser(ctx, userDTO); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.authClient.Authenticate(ctx, userDomain.ID())
	if err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	return &RegisterResponse{
		UserId:       userDomain.ID(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserUsecase) LoginUser(ctx context.Context, loginIdentifier, password string) (*LoginResponse, error) {
	if loginIdentifier == "" || password == "" {
		return nil, customErrors.ErrInvalidInput
	}

	var userDTO *dto.User
	var err error

	if utils.IsEmail(loginIdentifier) {
		userDTO, err = s.userRepo.GetUserByEmail(ctx, loginIdentifier)
	} else {
		userDTO, err = s.userRepo.GetUserByUsername(ctx, loginIdentifier)
	}

	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("login error: %w", err)
	}

	userDomain := domain.NewUserFromDTO(userDTO)
	if !userDomain.CheckPassword(password) {
		return nil, customErrors.ErrUnauthorized
	}

	accessToken, refreshToken, err := s.authClient.Authenticate(ctx, userDomain.ID())
	if err != nil {
		return nil, fmt.Errorf("authentication error: %s", err)
	}

	return &LoginResponse{
		UserId:       userDomain.ID(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserUsecase) ValidateToken(ctx context.Context, accessToken string) (bool, string, error) {
	if accessToken == "" {
		return false, "", customErrors.ErrInvalidInput
	}

	valid, userId, err := s.authClient.ValidateToken(ctx, accessToken)
	if err != nil {
		return false, "", fmt.Errorf("token validation error: %s", err)
	}
	return valid, userId, nil
}

func (s *UserUsecase) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	if refreshToken == "" {
		return nil, customErrors.ErrInvalidInput
	}

	accessToken, newRefreshToken, err := s.authClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token error: %s", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *UserUsecase) Logout(ctx context.Context, refreshToken string) (bool, error) {
	if refreshToken == "" {
		return false, customErrors.ErrInvalidInput
	}

	success, err := s.authClient.Logout(ctx, refreshToken)
	if err != nil {
		return false, fmt.Errorf("logout error: %w", err)
	}

	return success, nil
}

func (s *UserUsecase) GetUserByID(ctx context.Context, userId string) (domain.User, error) {
	if userId == "" {
		return nil, customErrors.ErrInvalidInput
	}

	userDTO, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get user error: %w", err)
	}

	return domain.NewUserFromDTO(userDTO), nil
}
