package config

import (
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/spf13/viper"
)

type Config struct {
	DBConfig db.DBConfig `mapstructure:"db"`
	GRPC     GRPC        `mapstructure:"grpc"`
	HTTP     HTTP        `mapstructure:"http"`
	Services Services    `mapstructure:"services"`
}

type GRPC struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type HTTP struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Service struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Services struct {
	AuthService Service `mapstructure:"auth_service"`
}

func LoadConfig(env string) (*Config, error) {
	v := viper.New()

	v.AddConfigPath("internal/config")
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")

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
