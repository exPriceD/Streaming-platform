package repository

import (
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/model"
	"log/slog"
)

var (
	ErrTokenNotFound = errors.New("refresh token not found")
)

type tokenRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewTokenRepository(db *sql.DB, log *slog.Logger) TokenRepository {
	return &tokenRepository{
		db:  db,
		log: log,
	}
}

func (r *tokenRepository) SaveRefreshToken(token *entity.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at)
        VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (token) DO NOTHING
    `
	_, err := r.db.Exec(query, token.UserID, token.Token, token.ExpiresAt, false, token.CreatedAt)
	if err != nil {
		r.log.Error("Failed to save refresh token", slog.String("error", err.Error()), slog.String("user_id", token.UserID.String()))
	}

	r.log.Info("Refresh token is saved", slog.String("user_id", token.UserID.String()))
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
		r.log.Warn("Refresh token not found", slog.String("token", tokenStr))
		return nil, ErrTokenNotFound
	}

	if err != nil {
		r.log.Error("Failed to get refresh token", slog.String("token", tokenStr), slog.String("error", err.Error()))
		return nil, err
	}

	r.log.Info("Refresh token retrieved", slog.String("user_id", token.UserID.String()))
	return &token, nil
}

func (r *tokenRepository) RevokeRefreshToken(tokenStr string) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.log.Error("Failed to begin transaction", slog.String("error", err.Error()))
		return err
	}

	query := `
        UPDATE refresh_tokens SET revoked = true WHERE token = $1
    `
	result, err := tx.Exec(query, tokenStr)
	if err != nil {
		r.log.Error("Failed to revoke refresh token", slog.String("token", tokenStr), slog.String("error", err.Error()))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		_ = tx.Rollback()
		r.log.Warn("Attempted to revoke non-existing refresh token", slog.String("token", tokenStr))
		return ErrTokenNotFound
	}

	r.log.Info("Refresh token revoked", slog.String("token", tokenStr))
	return tx.Commit()
}

func (r *tokenRepository) DeleteExpiredRefreshTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	result, err := r.db.Exec(query)
	if err != nil {
		r.log.Error("Failed to delete expired refresh tokens", slog.String("error", err.Error()))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		r.log.Info("Deleted expired refresh tokens", slog.Int64("count", rowsAffected))
	} else {
		r.log.Warn("No expired refresh tokens found for deletion")
	}
	return nil
}

func mapTokenModelToEntity(tokenModel *model.RefreshTokenModel) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        tokenModel.ID,
		UserID:    tokenModel.UserID,
		Token:     tokenModel.Token,
		ExpiresAt: tokenModel.ExpiresAt,
		Revoked:   tokenModel.Revoked,
		CreatedAt: tokenModel.CreatedAt,
	}
}
