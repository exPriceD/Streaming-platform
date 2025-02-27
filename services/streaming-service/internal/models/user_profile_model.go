package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                 uuid.UUID `db:"id"`                  // Уникальный идентификатор профиля пользователя
	ChannelName        string    `db:"channel_name"`        // Название канала пользователя
	ChannelDescription string    `db:"channel_description"` // Описание канала пользователя
	StreamKey          string    `db:"stream_key"`          // Ключ стрима
	IsLive             bool      `db:"is_live"`             // Флаг, указывающий на то, что пользователь в данный момент стримит
	CreatedAt          time.Time `db:"created_at"`          // Дата создания записи в БД
	UpdatedAt          time.Time `db:"updated_at"`          // Дата последнего обновления
}
