package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	customErrors "github.com/exPriceD/Streaming-platform/services/user-service/internal/errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/model"
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

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	now := time.Now()

	_, err := r.db.ExecContext(ctx, queryCreateUser,
		user.ID, user.Username, user.Email, user.PasswordHash, user.AvatarURL,
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

func (r *UserRepository) getUserByField(ctx context.Context, field, value string) (*entity.User, error) {
	query := fmt.Sprintf(queryGetUserByField, field)
	var user model.User

	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL,
		&user.ConsentToDataProcessing, &user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by %s: %w", field, err)
	}

	mappedUser := mapModelToEntity(&user)
	return mappedUser, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.getUserByField(ctx, "email", email)
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return r.getUserByField(ctx, "username", username)
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	now := time.Now()

	err := r.db.QueryRowContext(ctx, queryUpdateUser,
		user.Username, user.Email, user.AvatarURL, now, user.ID,
	).Scan(&user.ID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	user.UpdatedAt = now

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	var user model.User

	err := r.db.QueryRowContext(ctx, queryGetUserByID, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.AvatarURL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	mappedUser := mapModelToEntity(&user)

	return mappedUser, nil
}

func mapModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		ID:                      user.ID,
		Username:                user.Username,
		Email:                   user.Email,
		PasswordHash:            user.PasswordHash,
		AvatarURL:               user.AvatarURL,
		ConsentToDataProcessing: user.ConsentToDataProcessing,
		CreatedAt:               user.CreatedAt,
		UpdatedAt:               user.UpdatedAt,
	}
}
