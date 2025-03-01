package entity

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	id                      uuid.UUID
	username                string
	email                   string
	passwordHash            string
	avatarURL               string
	consentToDataProcessing bool
	createdAt               time.Time
	updatedAt               time.Time
}

type Config struct {
	BcryptCost           int
	DefaultAvatar        string
	IDGenerator          func() uuid.UUID
	TimeNow              func() time.Time
	SelectRandomAvatarFn func() (string, error)
}

func DefaultConfig() Config {
	return Config{
		BcryptCost:           bcrypt.DefaultCost,
		DefaultAvatar:        "person-1.png",
		IDGenerator:          uuid.New,
		TimeNow:              time.Now,
		SelectRandomAvatarFn: utils.SelectRandomAvatar,
	}
}

func NewUser(username, email, password, confirmPassword string, consent bool, cfg Config) (*User, error) {
	if !consent {
		return nil, ErrNoConsent
	}

	if err := validation.ValidateEmail(email); err != nil {
		return nil, WrapValidationError("email", err)
	}

	if err := validation.ValidateUsername(username); err != nil {
		return nil, WrapValidationError("username", err)
	}

	if password != confirmPassword {
		return nil, ErrPasswordsMismatch
	}

	if err := validation.ValidatePassword(password); err != nil {
		return nil, WrapValidationError("password", err)
	}

	hashedPassword, err := hashPassword(password, cfg.BcryptCost)
	if err != nil {
		return nil, WrapValidationError("password", ErrHashingPassword)
	}

	avatar, err := cfg.SelectRandomAvatarFn()
	if err != nil {
		avatar = cfg.DefaultAvatar
	}

	avatarURL := utils.GetAvatarPath(avatar)

	id := cfg.IDGenerator()
	now := cfg.TimeNow()

	return &User{
		id:                      id,
		username:                username,
		email:                   email,
		passwordHash:            hashedPassword,
		avatarURL:               avatarURL,
		consentToDataProcessing: consent,
		createdAt:               now,
		updatedAt:               now,
	}, nil
}

func (u *User) ID() uuid.UUID                 { return u.id }
func (u *User) Username() string              { return u.username }
func (u *User) Email() string                 { return u.email }
func (u *User) AvatarURL() string             { return u.avatarURL }
func (u *User) ConsentToDataProcessing() bool { return u.consentToDataProcessing }
func (u *User) CreatedAt() time.Time          { return u.createdAt }
func (u *User) UpdatedAt() time.Time          { return u.updatedAt }

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)) == nil
}

func hashPassword(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", ErrHashingPassword
	}
	return string(bytes), nil
}
