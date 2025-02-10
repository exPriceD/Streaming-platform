package repository

import (
	"database/sql"
	"errors"
	"github.com/exPriceD/Streaming-platform/services/auth-service/internal/entities"
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
		INSERT INTO users (id, email, password_hash, consent_to_data_processing, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()

	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.ConsentToDataProcessing, now, now)

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
        SELECT id, email, password_hash, consent_to_data_processing
        FROM users
        WHERE email = $1
    `

	var user entities.User
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.ConsentToDataProcessing)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("the user was not found")
	}

	return &user, err
}

func (r *userRepository) GetUserByUsername(username string) (*entities.User, error) {
	query := `
        SELECT id, username, email, password_hash, consent_to_data_processing
        FROM users
        WHERE username = $1
    `

	var user entities.User
	err := r.db.QueryRow(query, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.ConsentToDataProcessing)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("пользователь не найден")
	}

	return &user, err
}
