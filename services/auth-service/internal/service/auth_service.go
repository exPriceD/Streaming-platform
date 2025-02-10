package service

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, email, password string, consent bool) (*entities.User, error) {
	user, err := entities.NewUser(username, email, password, consent)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Authenticate(identifier, password string, isEmail bool) (*entities.User, error) {
	var user *entities.User
	var err error

	if isEmail {
		user, err = s.userRepo.GetUserByEmail(identifier)
	} else {
		user, err = s.userRepo.GetUserByUsername(identifier)
	}

	if err != nil {
		return nil, errors.New("the user was not found")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
