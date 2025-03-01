package domain

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/dto"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/utils"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User interface {
	ID() string
	Username() string
	Email() string
	AvatarURL() string
	ConsentToDataProcessing() bool
	CreatedAt() time.Time
	UpdatedAt() time.Time
	CheckPassword(password string) bool
	ToDTO() *dto.User
}

// user реализует интерфейс User.
type user struct {
	id                      uuid.UUID
	username                string
	email                   string
	passwordHash            string
	avatarURL               string
	consentToDataProcessing bool
	createdAt               time.Time
	updatedAt               time.Time
}

// Config определяет конфигурацию для создания и обновления пользователя.
type Config struct {
	BcryptCost           int
	DefaultAvatar        string
	IDGenerator          func() uuid.UUID
	TimeNow              func() time.Time
	SelectRandomAvatarFn func() (string, error)
}

// DefaultConfig возвращает стандартный конфиг.
func DefaultConfig() Config {
	return Config{
		BcryptCost:           bcrypt.DefaultCost,
		DefaultAvatar:        "person-1.png",
		IDGenerator:          uuid.New,
		TimeNow:              time.Now,
		SelectRandomAvatarFn: utils.SelectRandomAvatar,
	}
}

// NewUser создаёт нового пользователя с заданной конфигурацией.
func NewUser(username, email, passwordHash string, consent bool, cfg Config) (User, error) {
	if !consent {
		return nil, ErrNoConsent
	}

	if err := validation.ValidateEmail(email); err != nil {
		return nil, WrapValidationError("email", err)
	}

	if err := validation.ValidateUsername(username); err != nil {
		return nil, WrapValidationError("username", err)
	}

	avatar, err := cfg.SelectRandomAvatarFn()
	if err != nil {
		avatar = cfg.DefaultAvatar
	}

	avatarURL := utils.GetAvatarPath(avatar)

	id := cfg.IDGenerator()
	now := cfg.TimeNow()

	return &user{
		id:                      id,
		username:                username,
		email:                   email,
		passwordHash:            passwordHash,
		avatarURL:               avatarURL,
		consentToDataProcessing: consent,
		createdAt:               now,
		updatedAt:               now,
	}, nil
}

// NewUserFromDTO создаёт пользователя из DTO.
func NewUserFromDTO(d *dto.User) User {
	return &user{
		id:                      d.Id,
		username:                d.Username,
		email:                   d.Email,
		passwordHash:            d.PasswordHash,
		avatarURL:               d.AvatarURL,
		consentToDataProcessing: d.ConsentToDataProcessing,
		createdAt:               d.CreatedAt,
		updatedAt:               d.UpdatedAt,
	}
}

func (u *user) ID() string {
	return u.id.String()
}

func (u *user) Username() string {
	return u.username
}

func (u *user) Email() string {
	return u.email
}

func (u *user) AvatarURL() string {
	return u.avatarURL
}

func (u *user) ConsentToDataProcessing() bool {
	return u.consentToDataProcessing
}

func (u *user) CreatedAt() time.Time {
	return u.createdAt
}

func (u *user) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *user) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)) == nil
}

func (u *user) ToDTO() *dto.User {
	return &dto.User{
		Id:                      u.id,
		Username:                u.username,
		Email:                   u.email,
		PasswordHash:            u.passwordHash,
		AvatarURL:               u.avatarURL,
		ConsentToDataProcessing: u.consentToDataProcessing,
		CreatedAt:               u.createdAt,
		UpdatedAt:               u.updatedAt,
	}
}
