package entity

import (
	"errors"
	"fmt"
)

var (
	// ErrNoConsent возвращается, когда пользователь не дал согласие на обработку данных.
	ErrNoConsent = errors.New("consent to the processing of personal data is required")

	// ErrPasswordsMismatch возвращается, когда пароли не совпадают.
	ErrPasswordsMismatch = errors.New("passwords do not match")

	// ErrHashingPassword возвращается при ошибке хеширования пароля.
	ErrHashingPassword = errors.New("failed to hash password")
)

// WrapValidationError оборачивает ошибки валидации с добавлением контекста.
func WrapValidationError(field string, err error) error {
	return fmt.Errorf("%s: %w", field, err)
}
