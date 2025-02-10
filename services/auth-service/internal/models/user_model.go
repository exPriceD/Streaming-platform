package models

import (
	"github.com/google/uuid"
	"time"
)

type UserModel struct {
	ID                      uuid.UUID `db:"id"`
	Username                string    `db:"username"`
	Email                   string    `db:"email"`
	PasswordHash            string    `db:"password_hash"`
	ConsentToDataProcessing bool      `db:"consent_to_data_processing"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}
