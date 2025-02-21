package auth

import (
	"fmt"
	"time"

	"github.com/exPriceD/Streaming-platform/config"
	"github.com/golang-jwt/jwt/v4"
)

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secretKey string
	expiry    time.Duration
}

// NewJWTManager создаёт новый JWTManager
func NewJWTManager(JWTConfig config.JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey: JWTConfig.SecretKey,
		expiry:    JWTConfig.AccessTokenDuration, // Используем время жизни Access-токена
	}
}

// GenerateToken создаёт JWT токен для пользователя
func (jm *JWTManager) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(jm.expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

// VerifyToken проверяет валидность JWT токена
func (jm *JWTManager) VerifyToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jm.secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("неверный токен: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("неверный токен")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("неверный формат user_id")
	}

	return userID, nil
}
