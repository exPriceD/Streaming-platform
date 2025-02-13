package config

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

// StreamingServiceConfig — общая конфигурация streaming-service
type StreamingServiceConfig struct {
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
	Server ServerConfig `yaml:"server"`
}

// LoadStreamingConfig загружает конфигурацию streaming-service
func LoadStreamingConfig() (*StreamingServiceConfig, error) {
	data, err := LoadYAML()
	if err != nil {
		return nil, err
	}

	streamingServiceRaw, ok := data["streaming_service"]
	if !ok {
		return nil, fmt.Errorf("конфигурация streaming_service не найдена")
	}

	streamingServiceData, err := yaml.Marshal(streamingServiceRaw)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга streaming_service: %v", err)
	}

	var streamingConfig StreamingServiceConfig
	if err := yaml.Unmarshal(streamingServiceData, &streamingConfig); err != nil {
		return nil, fmt.Errorf("ошибка парсинга streaming_service: %v", err)
	}

	// Преобразуем продолжительность токенов в time.Duration
	streamingConfig.JWT.AccessTokenDuration *= time.Minute
	streamingConfig.JWT.RefreshTokenDuration *= time.Hour

	return &streamingConfig, nil
}
