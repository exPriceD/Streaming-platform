package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// StreamingServiceConfig — общая конфигурация streaming-service
type StreamingServiceConfig struct {
	DB     DBConfig     `yaml:"db"`
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

	return &streamingConfig, nil
}
