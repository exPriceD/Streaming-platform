package domain_test

import (
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/domain"
	"testing"
)

// BenchmarkNewUser измеряет производительность создания пользователя.
func BenchmarkNewUser(b *testing.B) {
	cfg := defaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := domain.NewUser("testuser", "test@example.com", "hash", true, cfg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
