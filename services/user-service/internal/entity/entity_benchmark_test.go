package entity_test

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"testing"
)

// BenchmarkNewUser измеряет производительность создания пользователя.
func BenchmarkNewUser(b *testing.B) {
	cfg := defaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := entity.NewUser("testuser", "test@example.com", "password123", "password123", true, cfg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
