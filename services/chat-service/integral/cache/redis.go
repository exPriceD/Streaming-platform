package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/exPriceD/Streaming-platform/services/chat-service/integral/entity"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisCache управляет кешированием сообщений и банов
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache создает новый объект кеша
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: client}
}

// SaveMessages кэширует сообщения чата
func (r *RedisCache) SaveMessages(ctx context.Context, streamID uuid.UUID, messages []*entity.ChatMessage) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("chat:%s:messages", streamID)
	return r.client.Set(ctx, key, data, time.Minute*10).Err()
}

// GetMessages получает сообщения из кеша
func (r *RedisCache) GetMessages(ctx context.Context, streamID uuid.UUID) ([]*entity.ChatMessage, error) {
	key := fmt.Sprintf("chat:%s:messages", streamID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var messages []*entity.ChatMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// BanUser кеширует информацию о бане пользователя
func (r *RedisCache) BanUser(ctx context.Context, streamID, userID uuid.UUID) error {
	key := fmt.Sprintf("chat:%s:banned", streamID)
	return r.client.SAdd(ctx, key, userID.String()).Err()
}

// IsUserBanned проверяет, забанен ли пользователь
func (r *RedisCache) IsUserBanned(ctx context.Context, streamID, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("chat:%s:banned", streamID)
	return r.client.SIsMember(ctx, key, userID.String()).Result()
}
