package entities

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID                      uuid.UUID
	Username                string
	Email                   string
	PasswordHash            string
	ConsentToDataProcessing bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func NewUser(username, email, password string, consent bool) (*User, error) {
	if !consent {
		return nil, errors.New("consent to the processing of personal data is required")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		ID:                      uuid.New(),
		Username:                username,
		Email:                   email,
		PasswordHash:            hashedPassword,
		ConsentToDataProcessing: consent,
		CreatedAt:               now,
		UpdatedAt:               now,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
