package entity

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id                      uuid.UUID
	Username                string
	Email                   string
	PasswordHash            string
	AvatarURL               string
	ConsentToDataProcessing bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func NewUser(username, email, password, confirmPassword string, consent bool) (*User, error) {
	if !consent {
		return nil, errors.New("consent to the processing of personal data is required")
	}

	if !validation.ValidateEmail(email) {
		return nil, errors.New("incorrect email")
	}

	if !validation.ValidateUsername(username) {
		return nil, errors.New("username must be alphanumeric and between 3 and 30 characters long")
	}

	if password != confirmPassword {
		return nil, errors.New("passwords do not match")
	}

	if !validation.ValidatePassword(password) {
		return nil, errors.New("password must be at least 6 characters long")
	}

	hashedPassword, err := hashPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	avatar, err := utils.SelectRandomAvatar()
	if err != nil {
		return nil, err
	}

	avatarURL := utils.GetAvatarPath(avatar)

	now := time.Now()

	return &User{
		Id:                      uuid.New(),
		Username:                username,
		Email:                   email,
		PasswordHash:            hashedPassword,
		AvatarURL:               avatarURL,
		ConsentToDataProcessing: consent,
		CreatedAt:               now,
		UpdatedAt:               now,
	}, nil
}

func hashPassword(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}
