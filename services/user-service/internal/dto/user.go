package dto

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id                      uuid.UUID `db:"id"`
	Username                string    `db:"username"`
	Email                   string    `db:"email"`
	PasswordHash            string    `db:"password_hash"`
	AvatarURL               string    `db:"avatar_url"`
	ConsentToDataProcessing bool      `db:"consent_to_data_processing"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}
