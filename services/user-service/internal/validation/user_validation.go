package validation

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	MinUsernameLength = 3
	MaxUsernameLength = 30
	MinPasswordLength = 6
)

var (
	// ErrInvalidUsernameLength возвращается, когда длина имени пользователя недопустима.
	ErrInvalidUsernameLength = fmt.Errorf("username must be between %d and %d characters long", MinUsernameLength, MaxUsernameLength)

	// ErrUsernameNotAlphanumeric возвращается, когда имя пользователя содержит недопустимые символы.
	ErrUsernameNotAlphanumeric = errors.New("username must contain only alphanumeric characters and underscores")

	// ErrInvalidEmail возвращается, когда формат email недопустим.
	ErrInvalidEmail = errors.New("email format is invalid")

	// ErrPasswordTooShort возвращается, когда пароль слишком короткий.
	ErrPasswordTooShort = fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
)

// ValidateUsername - проверка корректности username
func ValidateUsername(username string) error {
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return ErrInvalidUsernameLength
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return ErrUsernameNotAlphanumeric
	}
	return nil
}

// ValidateEmail - проверка корректности email
func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword - проверка пароля
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	return nil
}
