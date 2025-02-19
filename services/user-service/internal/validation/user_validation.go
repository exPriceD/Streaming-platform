package validation

import (
	"regexp"
)

// ValidateUsername - проверка корректности username
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return false
	}
	return true
}

// ValidateEmail - проверка корректности email
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ValidatePassword - проверка пароля
func ValidatePassword(password string) bool {
	return len(password) >= 6
}
