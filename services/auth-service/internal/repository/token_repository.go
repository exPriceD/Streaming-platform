package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/model"
	"log/slog"
)

var (
	ErrTokenNotFound = errors.New("refresh token not found")
)

type TokenRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at)
        VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (token) DO NOTHING
    `
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.Token, token.ExpiresAt, false, token.CreatedAt)
	if err != nil {
		return err
	}

	return err
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, tokenStr string) (*entity.RefreshToken, error) {
	query := `
        SELECT id, user_id, token, expires_at, revoked, created_at
        FROM refresh_tokens
        WHERE token = $1
    `
	var token model.RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.Revoked, &token.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTokenNotFound
	}

	if err != nil {
		return nil, err
	}
	return mapTokenModelToEntity(&token), nil
}

func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `
        UPDATE refresh_tokens SET revoked = true WHERE token = $1
    `
	result, err := tx.ExecContext(ctx, query, tokenStr)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return ErrTokenNotFound
	}

	return tx.Commit()
}

func (r *TokenRepository) DeleteExpiredRefreshTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func mapTokenModelToEntity(tokenModel *model.RefreshToken) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        tokenModel.ID,
		UserID:    tokenModel.UserID,
		Token:     tokenModel.Token,
		ExpiresAt: tokenModel.ExpiresAt,
		Revoked:   tokenModel.Revoked,
		CreatedAt: tokenModel.CreatedAt,
	}
}
