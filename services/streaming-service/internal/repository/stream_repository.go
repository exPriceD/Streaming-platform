package repository

import (
	"database/sql"
	"errors"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"
	"github.com/google/uuid"
)

// StreamRepository управляет доступом к данным о стримах в БД.
type StreamRepository struct {
	db *sql.DB
}

// NewStreamRepository создаёт новый репозиторий потоков.
func NewStreamRepository(db *sql.DB) *StreamRepository {
	return &StreamRepository{db: db}
}

// CreateStream добавляет новый стрим в БД.
func (r *StreamRepository) CreateStream(stream models.Stream) error {
	query := `INSERT INTO streams (id, title, user_id, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, stream.ID, stream.Title, stream.UserID, stream.Status, stream.CreatedAt, stream.UpdatedAt)
	return err
}

// GetStreamByID получает стрим по ID.
func (r *StreamRepository) GetStreamByID(id uuid.UUID) (*models.Stream, error) {
	query := `SELECT id, title, user_id, status, created_at, updated_at FROM streams WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var stream models.Stream
	err := row.Scan(&stream.ID, &stream.Title, &stream.UserID, &stream.Status, &stream.CreatedAt, &stream.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	return &stream, nil
}

// UpdateStream обновляет данные стрима.
func (r *StreamRepository) UpdateStream(stream models.Stream) error {
	query := `UPDATE streams SET title=$1, status=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.Exec(query, stream.Title, stream.Status, stream.UpdatedAt, stream.ID)
	return err
}

// DeleteStream удаляет стрим по ID.
func (r *StreamRepository) DeleteStream(id uuid.UUID) error {
	query := `DELETE FROM streams WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
