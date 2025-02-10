package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type JWTManager struct {
	secretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func NewJWTManager() *JWTManager {
	accessTokenDuration, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	refreshTokenDuration, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))

	return &JWTManager{
		secretKey:            os.Getenv("JWT_SECRET"),
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JWTManager) GenerateAccessToken(userID string) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.AccessTokenDuration)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) GenerateRefreshToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.RefreshTokenDuration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) ValidateToken(accessToken string) (*UserClaims, error) {
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
