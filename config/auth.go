package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"time"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type JWTConfig struct {
	SecretKey            string        `yaml:"secret_key"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AuthServiceConfig struct {
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
	Server ServerConfig `yaml:"server"`
}

func LoadAuthConfig() (*AuthServiceConfig, error) {
	data, err := LoadYAML()
	if err != nil {
		return nil, err
	}
	authServiceRaw, ok := data["auth_service"]
	if !ok {
		return nil, fmt.Errorf("the configuration for auth_service was not found")
	}

	authServiceData, err := yaml.Marshal(authServiceRaw)
	if err != nil {
		return nil, fmt.Errorf("auth_service marshalling error: %v", err)
	}

	var authServiceConfig AuthServiceConfig
	if err := yaml.Unmarshal(authServiceData, &authServiceConfig); err != nil {
		return nil, fmt.Errorf("auth_service unmarshalling error: %v", err)
	}

	return &authServiceConfig, nil
}
