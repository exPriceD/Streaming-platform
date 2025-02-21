package repository

import (
	"context"
	"errors"
	"time"

	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/entity"
	"github.com/exPriceD/Streaming-platform/services/chat-service/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// ChatRepositoryImpl реализует интерфейс ChatRepository
type ChatRepositoryImpl struct {
	mongoCollection *mongo.Collection
	pgDB            *gorm.DB
}

// NewChatRepository создает новый репозиторий чата
func NewChatRepository(mongoDB *mongo.Database, pgDB *gorm.DB) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{
		mongoCollection: mongoDB.Collection("messages"),
		pgDB:            pgDB,
	}
}

// SaveMessage сохраняет сообщение в MongoDB
func (r *ChatRepositoryImpl) SaveMessage(ctx context.Context, msg *entity.ChatMessage) error {
	_, err := r.mongoCollection.InsertOne(ctx, msg)
	return err
}

// GetMessages получает последние сообщения из MongoDB по streamID
func (r *ChatRepositoryImpl) GetMessages(ctx context.Context, streamID uuid.UUID, limit int) ([]*entity.ChatMessage, error) {
	var messages []*entity.ChatMessage
	cur, err := r.mongoCollection.Find(ctx, bson.M{"stream_id": streamID}, options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var msg entity.ChatMessage
		if err := cur.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

// DeleteMessage удаляет сообщение по ID
func (r *ChatRepositoryImpl) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	res, err := r.mongoCollection.DeleteOne(ctx, bson.M{"id": messageID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("сообщение не найдено")
	}
	return nil
}

// CreateRoom создает комнату в PostgreSQL
func (r *ChatRepositoryImpl) CreateRoom(ctx context.Context, streamID uuid.UUID, title string) (*entity.ChatRoom, error) {
	room := &model.ChatRoom{
		ID:          uuid.New(),
		StreamID:    streamID,
		StreamTitle: title,
		CreatedAt:   time.Now(),
		IsActive:    true,
	}
	if err := r.pgDB.WithContext(ctx).Create(room).Error; err != nil {
		return nil, err
	}
	return room.ToEntity(), nil
}

// GetRoom получает комнату по streamID
func (r *ChatRepositoryImpl) GetRoom(ctx context.Context, streamID uuid.UUID) (*model.ChatRoom, error) {
	var room model.ChatRoom
	if err := r.pgDB.WithContext(ctx).Where("stream_id = ? AND is_active = ?", streamID, true).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// CloseRoom закрывает комнату (делает неактивной)
func (r *ChatRepositoryImpl) CloseRoom(ctx context.Context, streamID uuid.UUID) error {
	res := r.pgDB.WithContext(ctx).Model(&model.ChatRoom{}).Where("stream_id = ?", streamID).Update("is_active", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("комната не найдена")
	}
	return nil
}

// BanUser блокирует пользователя в PostgreSQL
func (r *ChatRepositoryImpl) BanUser(ctx context.Context, streamID, userID uuid.UUID) error {
	ban := &model.ChatBan{
		ID:       uuid.New(),
		StreamID: streamID,
		UserID:   userID,
		BannedAt: time.Now(),
	}
	return r.pgDB.WithContext(ctx).Create(ban).Error
}

// UnbanUser удаляет блокировку пользователя
func (r *ChatRepositoryImpl) UnbanUser(ctx context.Context, streamID, userID uuid.UUID) error {
	res := r.pgDB.WithContext(ctx).Where("stream_id = ? AND user_id = ?", streamID, userID).Delete(&model.ChatBan{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("пользователь не найден в бан-листе")
	}
	return nil
}
