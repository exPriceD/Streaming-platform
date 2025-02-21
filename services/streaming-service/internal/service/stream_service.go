package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/models"
	"github.com/exPriceD/Streaming-platform/services/streaming-service/internal/repository"
	"github.com/google/uuid"
)

// StreamService управляет логикой работы со стримами.
type StreamService struct {
	streamRepo      repository.StreamRepositoryInterface
	ffmpegService   *FFmpegService
	userProfileRepo repository.UserProfileRepositoryInterface
	// Базовый URL RTMP-сервера, например: "rtmp://localhost/live"
	rtmpServerURL string
	db            *sql.DB
}

// NewStreamService создает новый экземпляр StreamService с необходимыми зависимостями.
func NewStreamService(
	streamRepo repository.StreamRepositoryInterface,
	ffmpegService *FFmpegService,
	userProfileRepo repository.UserProfileRepositoryInterface,
	rtmpServerURL string,
) *StreamService {
	return &StreamService{
		streamRepo:      streamRepo,
		ffmpegService:   ffmpegService,
		userProfileRepo: userProfileRepo,
		rtmpServerURL:   rtmpServerURL,
	}
}

// StartStream запускает новый стрим для пользователя.
// Он создает запись в БД, получает профиль пользователя для извлечения stream_key,
// формирует inputURL для FFmpeg и запускает процесс транскодинга.
func (s *StreamService) StartStream(userID, title, description string) (*models.Stream, error) {
	// Создаем новую запись стрима.
	stream := &models.Stream{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "LIVE",
		StartTime:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Сохраняем стрим в БД.
	err := s.streamRepo.CreateStream(*stream)
	if err != nil {
		log.Printf("Ошибка при создании стрима: %v", err)
		return nil, errors.New("не удалось создать стрим")
	}

	// Преобразуем userID в uuid и получаем профиль пользователя.
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("некорректный userID")
	}

	profile, err := s.userProfileRepo.GetUserProfileByID(uid)
	if err != nil || profile == nil {
		return nil, errors.New("профиль пользователя не найден")
	}

	// Формируем inputURL, используя базовый RTMP URL и stream_key из профиля.
	// Например, если rtmpServerURL = "rtmp://localhost/live" и stream_key = "abc123",
	// то inputURL будет "rtmp://localhost/live/abc123".
	inputURL := fmt.Sprintf("%s/%s", s.rtmpServerURL, profile.StreamKey)

	// Запускаем процесс FFmpeg для транскодинга.
	err = s.ffmpegService.StartStream(stream.ID.String(), inputURL)
	if err != nil {
		log.Printf("Ошибка при запуске FFmpeg: %v", err)
		return nil, errors.New("не удалось запустить процесс трансляции")
	}

	return stream, nil
}

// StopStream завершает активный стрим.
// Он обновляет статус стрима и останавливает соответствующий процесс FFmpeg.
func (s *StreamService) StopStream(streamID string) error {
	id, err := uuid.Parse(streamID)
	if err != nil {
		return errors.New("некорректный UUID стрима")
	}

	stream, err := s.streamRepo.GetStreamByID(id)
	if err != nil {
		return err
	}
	if stream == nil {
		return errors.New("стрим не найден")
	}

	stream.Status = "OFFLINE"
	stream.EndTime = time.Now()
	stream.UpdatedAt = time.Now()

	err = s.streamRepo.UpdateStream(*stream)
	if err != nil {
		return errors.New("не удалось обновить статус стрима")
	}

	// Останавливаем процесс FFmpeg.
	err = s.ffmpegService.StopStream(streamID)
	if err != nil {
		log.Printf("Ошибка при остановке FFmpeg: %v", err)
	}

	return nil
}

// GetStream возвращает информацию о стриме по его UUID.
func (s *StreamService) GetStream(streamID string) (*models.Stream, error) {
	id, err := uuid.Parse(streamID)
	if err != nil {
		return nil, errors.New("некорректный UUID стрима")
	}
	return s.streamRepo.GetStreamByID(id)
}

// UpdateStreamStatus обновляет статус стрима в базе данных.
// Например, если RTMP-сервер отправляет событие on_publish, статус будет "LIVE",
// а при on_publish_done – "OFFLINE".
func (s *StreamService) UpdateStreamStatus(streamID, status string) error {
	_, err := s.db.Exec("UPDATE streams SET status = $1 WHERE id = $2", status, streamID)
	if err != nil {
		log.Printf("Ошибка обновления статуса стрима (ID: %s) на %s: %v", streamID, status, err)
		return errors.New("не удалось обновить статус стрима")
	}
	return nil
}
