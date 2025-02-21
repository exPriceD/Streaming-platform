package config

import (
	"fmt"
	"github.com/exPriceD/Streaming-platform/pkg/db"
	"github.com/spf13/viper"
)

type Config struct {
	DBConfig db.DBConfig    `mapstructure:"db"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Services ServicesConfig `mapstructure:"services"`
}

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type HTTPConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServicesConfig struct {
	AuthService ServiceConfig `mapstructure:"auth_service"`
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
