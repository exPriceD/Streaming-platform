package repository

import (
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/models"
	"github.com/lib/pq"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entities.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, consent_to_data_processing, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()

	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, user.ConsentToDataProcessing, now, now)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return errors.New("the user with this email already exists")
		}
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	query := `
        SELECT id, email, password_hash, consent_to_data_processing, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var userModel models.UserModel
	err := r.db.QueryRow(query, email).
		Scan(&userModel.ID, &userModel.Email, &userModel.PasswordHash, &userModel.ConsentToDataProcessing, &userModel.CreatedAt, &userModel.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("the user was not found")
	}

	user := mapModelToEntity(&userModel)

	return user, err
}

func (r *userRepository) GetUserByUsername(username string) (*entities.User, error) {
	query := `
        SELECT id, username, email, password_hash, consent_to_data_processing, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	var userModel models.UserModel
	err := r.db.QueryRow(query, username).
		Scan(&userModel.ID, &userModel.Username, &userModel.Email, &userModel.PasswordHash, &userModel.ConsentToDataProcessing, &userModel.CreatedAt, &userModel.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("пользователь не найден")
	}

	user := mapModelToEntity(&userModel)

	return user, err
}

func mapModelToEntity(userModel *models.UserModel) *entities.User {
	return &entities.User{
		ID:                      userModel.ID,
		Username:                userModel.Username,
		Email:                   userModel.Email,
		PasswordHash:            userModel.PasswordHash,
		ConsentToDataProcessing: userModel.ConsentToDataProcessing,
		CreatedAt:               userModel.CreatedAt,
		UpdatedAt:               userModel.UpdatedAt,
	}
}
