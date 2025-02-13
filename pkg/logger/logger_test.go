package logger_test

import (
	"github.com/exPriceD/Streaming-platform/pkg/logger"
	"log/slog"
	"os"
	"testing"
	"time"
)

// TestLogger проверяет создание логов
func TestLogger(t *testing.T) {
	log := logger.InitLogger("test-service")

	log.Info("🚀 Логгер тестируется!", slog.String("env", "testing"))
	log.Warn("⚠️  Предупреждение!")
	log.Error("❌ Ошибка!", slog.String("trace_id", "abc-1234"))

	time.Sleep(500 * time.Millisecond)

	logDir := "logs"
	logFilePath := logDir + "/test-service.log"
	jsonLogFilePath := logDir + "/test-service.json"

	checkFileExists(t, logFilePath)
	checkFileExists(t, jsonLogFilePath)

	// Можно дописать проверку содержимого логов
}

// checkFileExists проверяет, существует ли файл
func checkFileExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		t.Logf("✅ Файл найден: %s", filePath)
	} else if os.IsNotExist(err) {
		t.Errorf("❌ Файл не найден: %s", filePath)
	} else {
		t.Errorf("⚠️  Ошибка при проверке файла %s: %v", filePath, err)
	}
}

// checkLogContent проверяет, содержится ли строка в файле
func checkLogContent(t *testing.T, filePath, expected string) {
	// прочитать файл
	// проверить содержимое через contains
}

// contains проверяет, есть ли подстрока в строке
func contains(str, substr string) {
	// проверка что строка содержит подстроку
}
