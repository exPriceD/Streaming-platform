package service

import (
	"errors"
	"log"
	"time"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/repository"
)

// StreamService управляет логикой работы со стримами
type StreamService struct {
	streamRepo repository.StreamRepository
}

// NewStreamService создает новый сервис управления стримами
func NewStreamService(streamRepo repository.StreamRepository) *StreamService {
	return &StreamService{
		streamRepo: streamRepo,
	}
}

// StartStream запускает стрим пользователя
func (s *StreamService) StartStream(userID, title, description string) (*models.Stream, error) {
	stream := &models.Stream{
		StreamID:    generateUUID(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "LIVE",
		StartTime:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.streamRepo.CreateStream(stream)
	if err != nil {
		log.Printf("Ошибка при создании стрима: %v", err)
		return nil, errors.New("не удалось создать стрим")
	}

	return stream, nil
}

// StopStream завершает активный стрим
func (s *StreamService) StopStream(streamID string) error {
	stream, err := s.streamRepo.GetStreamByID(streamID)
	if err != nil {
		return err
	}
	if stream == nil {
		return errors.New("стрим не найден")
	}

	stream.Status = "OFFLINE"
	stream.EndTime = time.Now()

	err = s.streamRepo.UpdateStream(stream)
	if err != nil {
		return errors.New("не удалось обновить статус стрима")
	}

	return nil
}

// GetStream получает информацию о стриме по ID
func (s *StreamService) GetStream(streamID string) (*models.Stream, error) {
	return s.streamRepo.GetStreamByID(streamID)
}

// generateUUID — простая заглушка для генерации UUID
func generateUUID() string {
	return time.Now().Format("20060102150405") // В реальном коде лучше использовать библиотеку "github.com/google/uuid"
}
