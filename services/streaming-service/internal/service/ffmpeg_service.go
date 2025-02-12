package service

import (
	"fmt"
	"os/exec"
	"sync"
)

// FFmpegService управляет процессами FFmpeg для трансляций.
type FFmpegService struct {
	// Хранит запущенные процессы FFmpeg.
	// Ключ — это streamID, значение — указатель на *exec.Cmd.
	processes map[string]*exec.Cmd
	mu        sync.Mutex // Мьютекс для потокобезопасного доступа к мапе.
}

// NewFFmpegService создает новый экземпляр FFmpegService.
func NewFFmpegService() *FFmpegService {
	return &FFmpegService{
		processes: make(map[string]*exec.Cmd),
	}
}

// StartStream запускает процесс FFmpeg для конкретного стрима.
// Параметры:
//   - streamID: уникальный идентификатор стрима (используется для отслеживания процесса)
//   - inputURL: URL входного потока (например, "rtmp://localhost/live/abc123")
func (s *FFmpegService) StartStream(streamID string, inputURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если для данного стрима уже запущен процесс, возвращаем ошибку.
	if _, exists := s.processes[streamID]; exists {
		return fmt.Errorf("стрим %s уже запущен", streamID)
	}

	// Команда FFmpeg для трансляции.
	// Здесь можно настроить нужные параметры кодирования, выходной формат и путь к плейлисту.
	cmd := exec.Command("ffmpeg",
		"-re",          // Читаем входной поток в режиме реального времени.
		"-i", inputURL, // Используем полученный inputURL (например, rtmp://localhost/live/abc123).
		"-c:v", "libx264", "-preset", "veryfast", "-b:v", "1500k",
		"-c:a", "aac", "-b:a", "128k",
		"-f", "hls",
		"-hls_time", "4",
		"-hls_playlist_type", "event",
		fmt.Sprintf("/path/to/hls/%s/index.m3u8", streamID), // Выходной путь для HLS сегментов.
	)

	// Запускаем процесс FFmpeg.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка запуска FFmpeg для стрима %s: %w", streamID, err)
	}

	// Сохраняем процесс в мапе.
	s.processes[streamID] = cmd
	return nil
}

// StopStream завершает процесс FFmpeg для указанного стрима.
func (s *FFmpegService) StopStream(streamID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cmd, exists := s.processes[streamID]
	if !exists {
		return fmt.Errorf("процесс для стрима %s не найден", streamID)
	}

	// Останавливаем процесс FFmpeg.
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("ошибка остановки FFmpeg для стрима %s: %w", streamID, err)
	}

	// Удаляем процесс из мапы.
	delete(s.processes, streamID)
	return nil
}
