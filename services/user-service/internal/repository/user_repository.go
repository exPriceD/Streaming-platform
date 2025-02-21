package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/user-service/internal/model"
	"github.com/lib/pq"
	"log"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, avatar_url, consent_to_data_processing, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()

	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, user.AvatarURL, user.ConsentToDataProcessing, now, now)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return errors.New("the user with this email already exists")
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	query := `
        SELECT id, username, email, password_hash, avatar_url, consent_to_data_processing, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user model.User
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL, &user.ConsentToDataProcessing, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("the user was not found")
	}

	mappedUser := mapModelToEntity(&user)

	return mappedUser, err
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	query := `
        SELECT id, username, email, password_hash, avatar_url, consent_to_data_processing, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	var user model.User
	err := r.db.QueryRow(query, username).
		Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL,
			&user.ConsentToDataProcessing, &user.CreatedAt, &user.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("the user was not found")
	}

	mappedUser := mapModelToEntity(&user)

	return mappedUser, err
}

func (r *UserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	query := `UPDATE users SET username = $1, email = $2, avatar_url = $3, updated_at = $4 WHERE id = $5 RETURNING id`

	now := time.Now()

	err := r.db.QueryRow(query, user.Username, user.Email, user.AvatarURL, now, user.ID).Scan(&user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	query := `SELECT id, username, email, avatar_url FROM users WHERE id = $1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
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
