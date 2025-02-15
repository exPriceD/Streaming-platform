package repository

import (
	"database/sql"
	"errors"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"
	"github.com/google/uuid"
)

// UserProfileRepository реализует доступ к профилям пользователей в БД.
type UserProfileRepository struct {
	db *sql.DB
}

// NewUserProfileRepository создаёт новый репозиторий профилей пользователей.
func NewUserProfileRepository(db *sql.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

// CreateUserProfile добавляет новый профиль в БД.
func (r *UserProfileRepository) CreateUserProfile(profile models.UserProfile) error {
	query := `INSERT INTO user_profiles (id, channel_name, channel_description, stream_key, is_live, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, profile.ID, profile.ChannelName, profile.ChannelDescription, profile.StreamKey, profile.IsLive, profile.CreatedAt, profile.UpdatedAt)
	return err
}

// GetUserProfileByID получает профиль пользователя по ID.
func (r *UserProfileRepository) GetUserProfileByID(id uuid.UUID) (*models.UserProfile, error) {
	query := `SELECT id, channel_name, channel_description, stream_key, is_live, created_at, updated_at FROM user_profiles WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var profile models.UserProfile
	err := row.Scan(&profile.ID, &profile.ChannelName, &profile.ChannelDescription, &profile.StreamKey, &profile.IsLive, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user profile not found")
		}
		return nil, err
	}

	return &profile, nil
}

// UpdateUserProfile обновляет данные профиля пользователя.
func (r *UserProfileRepository) UpdateUserProfile(profile models.UserProfile) error {
	query := `UPDATE user_profiles
			  SET channel_name=$1, channel_description=$2, stream_key=$3, is_live=$4, updated_at=$5
			  WHERE id=$6`
	_, err := r.db.Exec(query, profile.ChannelName, profile.ChannelDescription, profile.StreamKey, profile.IsLive, profile.UpdatedAt, profile.ID)
	return err
}

// UpdateLiveStatus обновляет статус стрима пользователя.
func (r *UserProfileRepository) UpdateLiveStatus(userID uuid.UUID, isLive bool) error {
	query := `UPDATE user_profiles SET is_live=$1, updated_at=NOW() WHERE id=$2`
	_, err := r.db.Exec(query, isLive, userID)
	return err
}
