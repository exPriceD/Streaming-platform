package config

import (
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	DB   db.DBConfig `mapstructure:"db"`
	JWT  JWTConfig   `mapstructure:"jwt"`
	GRPC GRPCConfig  `mapstructure:"grpc"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SecretKey            string        `mapstructure:"secret_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

func LoadConfig(configInput string) (*Config, error) {
	v := viper.New()

	if filepath.IsAbs(configInput) || (len(configInput) > 5 && configInput[len(configInput)-5:] == ".yaml") {
		if _, err := os.Stat(configInput); os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", configInput)
		}
		v.SetConfigFile(configInput)
	} else {
		v.AddConfigPath("internal/config")
		v.SetConfigName(fmt.Sprintf("config.%s", configInput))
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading the config: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config parsing error: %v", err)
	}

	return &cfg, nil
}
