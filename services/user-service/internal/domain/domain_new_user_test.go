package domain_test

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/domain"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

const pathToAvatarDir = "assets/avatars/"

// defaultConfig возвращает стандартную конфигурацию для тестов.
func defaultConfig() domain.Config {
	return domain.Config{
		BcryptCost:           10,
		DefaultAvatar:        "default.png",
		IDGenerator:          func() uuid.UUID { return uuid.MustParse("550e8400-e29b-41d4-a716-446655440000") },
		TimeNow:              func() time.Time { return time.Unix(1234567890, 0) },
		SelectRandomAvatarFn: func() (string, error) { return "test.png", nil },
	}
}

// assertNewUser создаёт пользователя и проверяет базовые ожидания.
func assertNewUser(t *testing.T, username, email, password, confirmPassword string, consent bool, cfg domain.Config, wantErr error) domain.User {
	t.Helper()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	userDomain, err := domain.NewUser(username, email, string(hashedPassword), consent, cfg)
	if wantErr != nil {
		assert.Error(t, err)
		assert.True(t, errors.Is(err, wantErr), "expected error %v, got %v", wantErr, err)
		assert.Nil(t, userDomain)
		return nil
	}
	assert.NoError(t, err)
	assert.NotNil(t, userDomain)
	return userDomain
}

// TestNewUserSuccess проверяет успешное создание пользователя.
func TestNewUserSuccess(t *testing.T) {
	tests := []struct {
		name              string
		username          string
		email             string
		password          string
		avatarFails       bool
		wantAvatarURL     string
		wantUsername      string
		wantEmail         string
		wantConsent       bool
		wantCreatedAt     time.Time
		wantUpdatedAt     time.Time
		wantPasswordValid bool
	}{
		{
			name:              "BasicSuccess",
			username:          "testuser",
			email:             "test@example.com",
			password:          "password123",
			avatarFails:       false,
			wantAvatarURL:     pathToAvatarDir + "test.png",
			wantUsername:      "testuser",
			wantEmail:         "test@example.com",
			wantConsent:       true,
			wantCreatedAt:     time.Unix(1234567890, 0),
			wantUpdatedAt:     time.Unix(1234567890, 0),
			wantPasswordValid: true,
		},
		{
			name:              "MinUsernameLength",
			username:          "abc", // 3 символа
			email:             "test@example.com",
			password:          "password123",
			avatarFails:       false,
			wantAvatarURL:     pathToAvatarDir + "test.png",
			wantUsername:      "abc",
			wantEmail:         "test@example.com",
			wantConsent:       true,
			wantCreatedAt:     time.Unix(1234567890, 0),
			wantUpdatedAt:     time.Unix(1234567890, 0),
			wantPasswordValid: true,
		},
		{
			name:              "MaxUsernameLength",
			username:          "abcdefghijklmnopqrstuvwxyzabcd", // 30 символов
			email:             "test@example.com",
			password:          "password123",
			avatarFails:       false,
			wantAvatarURL:     pathToAvatarDir + "test.png",
			wantUsername:      "abcdefghijklmnopqrstuvwxyzabcd",
			wantEmail:         "test@example.com",
			wantConsent:       true,
			wantCreatedAt:     time.Unix(1234567890, 0),
			wantUpdatedAt:     time.Unix(1234567890, 0),
			wantPasswordValid: true,
		},
		{
			name:              "MinPasswordLength",
			username:          "testuser",
			email:             "test@example.com",
			password:          "pass12", // 6 символов
			avatarFails:       false,
			wantAvatarURL:     pathToAvatarDir + "test.png",
			wantUsername:      "testuser",
			wantEmail:         "test@example.com",
			wantConsent:       true,
			wantCreatedAt:     time.Unix(1234567890, 0),
			wantUpdatedAt:     time.Unix(1234567890, 0),
			wantPasswordValid: true,
		},
		{
			name:              "AvatarFails",
			username:          "testuser",
			email:             "test@example.com",
			password:          "password123",
			avatarFails:       true,
			wantAvatarURL:     pathToAvatarDir + "default.png",
			wantUsername:      "testuser",
			wantEmail:         "test@example.com",
			wantConsent:       true,
			wantCreatedAt:     time.Unix(1234567890, 0),
			wantUpdatedAt:     time.Unix(1234567890, 0),
			wantPasswordValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig()
			if tt.avatarFails {
				cfg.SelectRandomAvatarFn = func() (string, error) { return "", errors.New("avatar failed") }
			}

			user := assertNewUser(t, tt.username, tt.email, tt.password, tt.password, true, cfg, nil)
			assert.Equal(t, tt.wantUsername, user.Username())
			assert.Equal(t, tt.wantEmail, user.Email())
			assert.Equal(t, tt.wantAvatarURL, user.AvatarURL())
			assert.Equal(t, tt.wantConsent, user.ConsentToDataProcessing())
			assert.Equal(t, tt.wantCreatedAt, user.CreatedAt())
			assert.Equal(t, tt.wantUpdatedAt, user.UpdatedAt())
			assert.Equal(t, tt.wantPasswordValid, user.CheckPassword(tt.password))
			assert.Equal(t, cfg.IDGenerator().String(), user.ID())
		})
	}
}

// TestNewUserErrors проверяет сценарии с ошибками при создании пользователя.
func TestNewUserErrors(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		email           string
		password        string
		confirmPassword string
		consent         bool
		wantErr         error
	}{
		{
			name:            "NoConsent",
			username:        "testuser",
			email:           "test@example.com",
			password:        "password123",
			confirmPassword: "password123",
			consent:         false,
			wantErr:         domain.ErrNoConsent,
		},
		{
			name:            "InvalidEmail",
			username:        "testuser",
			email:           "invalid-email",
			password:        "password123",
			confirmPassword: "password123",
			consent:         true,
			wantErr:         validation.ErrInvalidEmail,
		},
		{
			name:            "InvalidUsernameTooShort",
			username:        "ab",
			email:           "test@example.com",
			password:        "password123",
			confirmPassword: "password123",
			consent:         true,
			wantErr:         validation.ErrInvalidUsernameLength,
		},
		{
			name:            "InvalidUsernameNotAlphanumeric",
			username:        "test#user",
			email:           "test@example.com",
			password:        "password123",
			confirmPassword: "password123",
			consent:         true,
			wantErr:         validation.ErrUsernameNotAlphanumeric,
		},
		{
			name:            "PasswordsMismatch",
			username:        "testuser",
			email:           "test@example.com",
			password:        "password123",
			confirmPassword: "password456",
			consent:         true,
			wantErr:         domain.ErrPasswordsMismatch,
		},
		{
			name:            "WeakPassword",
			username:        "testuser",
			email:           "test@example.com",
			password:        "pass",
			confirmPassword: "pass",
			consent:         true,
			wantErr:         validation.ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assertNewUser(t, tt.username, tt.email, tt.password, tt.confirmPassword, tt.consent, defaultConfig(), tt.wantErr)
		})
	}
}
