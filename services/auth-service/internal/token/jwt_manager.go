package token

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type JWTManager struct {
	secretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

func NewJWTManager(JWTConfig config.JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey:            JWTConfig.SecretKey,
		AccessTokenDuration:  JWTConfig.AccessTokenDuration,
		RefreshTokenDuration: JWTConfig.RefreshTokenDuration,
	}
}

func (m *JWTManager) GenerateTokens(userID uuid.UUID) (string, string, error) {
	accessClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.AccessTokenDuration)),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(m.secretKey))
	if err != nil {
		return "", "", err
	}

	refreshClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.RefreshTokenDuration)),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(m.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *JWTManager) ValidateAccessToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("incorrect data in the token")
	}

	return claims, nil
}

func (m *JWTManager) ValidateRefreshToken(refreshToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("недействительный refresh_token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("некорректные данные в refresh_token")
	}

	return claims, nil
}
