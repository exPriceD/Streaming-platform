package domain_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCheckPassword проверяет метод CheckPassword.
func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		input    string
		want     bool
	}{
		{
			name:     "CorrectPassword",
			password: "password123",
			input:    "password123",
			want:     true,
		},
		{
			name:     "IncorrectPassword",
			password: "password123",
			input:    "wrongpass",
			want:     false,
		},
		{
			name:     "EmptyInput",
			password: "password123",
			input:    "",
			want:     false,
		},
		{
			name:     "CaseSensitive",
			password: "Password123",
			input:    "password123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user := assertNewUser(t, "testuser", "test@example.com", tt.password, tt.password, true, defaultConfig(), nil)
			got := user.CheckPassword(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
