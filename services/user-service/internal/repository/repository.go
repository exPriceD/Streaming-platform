package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/dto"
	customErrors "github.com/exPriceD/Streaming-platform/services/user-service/internal/errors"
	"github.com/lib/pq"
	"time"
)

const (
	queryCreateUser = `
		INSERT INTO users (id, username, email, password_hash, avatar_url, consent_to_data_processing, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	queryGetUserByField = `
		SELECT id, username, email, password_hash, avatar_url, consent_to_data_processing, created_at, updated_at
		FROM users
		WHERE %s = $1
	`
	queryGetUserByID = `
		SELECT id, username, email, avatar_url
		FROM users
		WHERE id = $1
	`
	queryUpdateUser = `
		UPDATE users
		SET username = $1, email = $2, avatar_url = $3, updated_at = $4
		WHERE id = $5
		RETURNING id
	`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *dto.User) error {
	now := time.Now()

	_, err := r.db.ExecContext(ctx, queryCreateUser,
		user.Id, user.Username, user.Email, user.PasswordHash, user.AvatarURL,
		user.ConsentToDataProcessing, now, now,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return customErrors.ErrUserAlreadyExist
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *UserRepository) getUserByField(ctx context.Context, field, value string) (*dto.User, error) {
	query := fmt.Sprintf(queryGetUserByField, field)
	var user dto.User

	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL,
		&user.ConsentToDataProcessing, &user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by %s: %w", field, err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*dto.User, error) {
	return r.getUserByField(ctx, "email", email)
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*dto.User, error) {
	return r.getUserByField(ctx, "username", username)
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *dto.User) (*dto.User, error) {
	now := time.Now()

	err := r.db.QueryRowContext(ctx, queryUpdateUser,
		user.Username, user.Email, user.AvatarURL, now, user.Id,
	).Scan(&user.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	user.UpdatedAt = now

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userId string) (*dto.User, error) {
	var user dto.User

	err := r.db.QueryRowContext(ctx, queryGetUserByID, userId).Scan(
		&user.Id, &user.Username, &user.Email, &user.AvatarURL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by Id: %w", err)
	}

	return &user, nil
}
