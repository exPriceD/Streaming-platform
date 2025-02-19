package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type ChatServiceConfig struct {
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
}

func LoadChatConfig() (*ChatServiceConfig, error) {
	data, err := LoadYAML()
	if err != nil {
		return nil, err
	}

	chatServiceRaw, ok := data["chat_service"]
	if !ok {
		return nil, fmt.Errorf("конфигурация chat_service не найдена")
	}

	chatServiceData, err := yaml.Marshal(chatServiceRaw)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга chat_service: %v", err)
	}

	var chatConfig ChatServiceConfig
	if err := yaml.Unmarshal(chatServiceData, &chatConfig); err != nil {
		return nil, fmt.Errorf("ошибка парсинга chat_service: %v", err)
	}

	return &chatConfig, nil
}
