package repository

import (
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/models"
	"log"
)

var (
	ErrTokenNotFound = errors.New("refresh token not found")
)

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) SaveRefreshToken(token *entity.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at)
        VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (token) DO NOTHING
    `
	_, err := r.db.Exec(query, token.UserID, token.Token, token.ExpiresAt, false, token.CreatedAt)
	return err
}

func (r *tokenRepository) GetRefreshToken(tokenStr string) (*entity.RefreshToken, error) {
	query := `
        SELECT id, user_id, token, expires_at, revoked, created_at
        FROM refresh_tokens
        WHERE token = $1
    `
	var token entity.RefreshToken
	err := r.db.QueryRow(query, tokenStr).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.Revoked, &token.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("the token was not found")
	}

	return &token, nil
}

func (r *tokenRepository) RevokeRefreshToken(tokenStr string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `
        UPDATE refresh_tokens SET revoked = true WHERE token = $1
    `
	result, err := tx.Exec(query, tokenStr)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return ErrTokenNotFound
	}

	return tx.Commit()
}

func (r *tokenRepository) DeleteExpiredRefreshTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	result, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Outdated refresh_token has been removed - %d pcs.", rowsAffected)
	return nil
}

func mapTokenModelToEntity(tokenModel *models.RefreshTokenModel) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        tokenModel.ID,
		UserID:    tokenModel.UserID,
		Token:     tokenModel.Token,
		ExpiresAt: tokenModel.ExpiresAt,
		Revoked:   tokenModel.Revoked,
		CreatedAt: tokenModel.CreatedAt,
	}
}
