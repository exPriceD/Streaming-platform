package token

import (
	"errors"
	"github.com/exPriceD/Streaming-platform/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type JWTManager struct {
	secretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	log                  *slog.Logger
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

func NewJWTManager(JWTConfig config.JWTConfig, log *slog.Logger) *JWTManager {
	return &JWTManager{
		secretKey:            JWTConfig.SecretKey,
		AccessTokenDuration:  JWTConfig.AccessTokenDuration,
		RefreshTokenDuration: JWTConfig.RefreshTokenDuration,
		log:                  log,
	}
}

func (m *JWTManager) GenerateTokens(userID uuid.UUID) (string, string, int64, time.Time, error) {
	accessExpiresAt := time.Now().Add(m.AccessTokenDuration)
	accessClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(m.secretKey))
	if err != nil {
		m.log.Error("Failed to generate access token", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	refreshExpiresAt := time.Now().Add(m.RefreshTokenDuration)
	refreshClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(m.secretKey))
	if err != nil {
		m.log.Error("Failed to generate refresh token", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
		return "", "", 0, time.Time{}, err
	}

	m.log.Info("Generated tokens successfully",
		slog.String("user_id", userID.String()),
		slog.Time("access_expires_at", accessExpiresAt),
		slog.Time("refresh_expires_at", refreshExpiresAt),
	)

	expiresIn := int64(m.AccessTokenDuration.Seconds())

	return accessToken, refreshToken, expiresIn, accessExpiresAt, nil
}

func (m *JWTManager) ValidateToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		m.log.Warn("Invalid token", slog.String("error", err.Error()))
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		m.log.Warn("Invalid token: signature or expiration error")
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		m.log.Warn("Malformed token")
		return nil, errors.New("incorrect data in the token")
	}

	m.log.Info("Token validated", slog.String("user_id", claims.UserID.String()))
	return claims, nil
}
