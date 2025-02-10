package repository

import (
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"
	"time"
)

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) SaveRefreshToken(token *entities.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, revoked)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.Exec(query, token.UserID, token.Token, token.ExpiresAt, false)
	return err
}

func (r *tokenRepository) GetRefreshToken(tokenStr string) (*entities.RefreshToken, error) {
	query := `
        SELECT id, user_id, token, expires_at, revoked
        FROM refresh_tokens
        WHERE token = $1
    `
	var token entities.RefreshToken
	err := r.db.QueryRow(query, tokenStr).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.Revoked)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("the token was not found")
	}

	if token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("the token is invalid")
	}

	return &token, nil
}

func (r *tokenRepository) RevokeRefreshToken(tokenStr string) error {
	query := `
        UPDATE refresh_tokens SET revoked = true WHERE token = $1
    `
	_, err := r.db.Exec(query, tokenStr)
	return err
}
