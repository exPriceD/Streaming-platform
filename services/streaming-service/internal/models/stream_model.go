package models

import (
	"time"

	"github.com/google/uuid"
)

type Stream struct {
	ID          uuid.UUID `db:"id"`          // Уникальный идентификатор стрима (UUID)
	UserID      string    `db:"user_id"`     // ID пользователя, создавшего стрим
	Title       string    `db:"title"`       // Название стрима
	Description string    `db:"description"` // Описание стрима
	Thumbnail   string    `db:"thumbnail"`   // URL миниатюры стрима
	Status      string    `db:"status"`      // Текущий статус стрима (LIVE, OFFLINE)
	StartTime   time.Time `db:"start_time"`  // Время начала стрима
	EndTime     time.Time `db:"end_time"`    // Время окончания стрима (если завершен)
	CreatedAt   time.Time `db:"created_at"`  // Дата создания записи в БД
	UpdatedAt   time.Time `db:"updated_at"`  // Дата последнего обновления
}
