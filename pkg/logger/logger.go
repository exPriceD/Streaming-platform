package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newHandlers []slog.Handler
	for _, h := range m.handlers {
		newHandlers = append(newHandlers, h.WithAttrs(attrs))
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	var newHandlers []slog.Handler
	for _, h := range m.handlers {
		newHandlers = append(newHandlers, h.WithGroup(name))
	}
	return &multiHandler{handlers: newHandlers}
}

// MultiHandler создаёт кастомный обработчик, записывающий логи в несколько мест
func MultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

// globalLogger - мапа с логгерами для каждого сервиса
var globalLoggers = make(map[string]*slog.Logger)
var mu sync.Mutex

func formatTime() string {
	return time.Now().Format("02.01.2006 15:04:05")
}

// colorize добавляет цвета в логи (только для консоли)
func colorize(level slog.Level, message string) string {
	switch level {
	case slog.LevelInfo:
		return fmt.Sprintf("\033[32m%s\033[0m", message) // 🟢 Зелёный для INFO
	case slog.LevelWarn:
		return fmt.Sprintf("\033[33m%s\033[0m", message) // 🟠 Жёлтый для WARN
	case slog.LevelError:
		return fmt.Sprintf("\033[31m%s\033[0m", message) // 🔴 Красный для ERROR
	default:
		return message
	}
}

// InitLogger создаёт новый логгер для конкретного сервиса
func InitLogger(serviceName string) *slog.Logger {
	mu.Lock()
	defer mu.Unlock()

	if logger, exists := globalLoggers[serviceName]; exists {
		return logger
	}

	logDir := filepath.Join("services", serviceName, "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Ошибка создания папки логов: %v", err))
	}

	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.log", serviceName))
	jsonLogFilePath := filepath.Join(logDir, fmt.Sprintf("%s.json", serviceName))

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Ошибка создания текстового лог-файла: %v", err))
	}

	jsonFile, err := os.OpenFile(jsonLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Ошибка создания JSON лог-файла: %v", err))
	}

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.String("time", formatTime())
			}
			if a.Key == "msg" {
				return slog.String("msg", colorize(slog.LevelInfo, a.Value.String()))
			}
			return a
		},
		AddSource: true,
	})

	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	jsonHandler := slog.NewJSONHandler(jsonFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(MultiHandler(consoleHandler, fileHandler, jsonHandler))
	globalLoggers[serviceName] = logger

	return logger
}
