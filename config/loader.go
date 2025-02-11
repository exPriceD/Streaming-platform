package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func LoadYAML() (map[string]interface{}, error) {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("couldn't READ the configuration file: %v", err)
	}
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("couldn't PARSE the configuration file: %v", err)
	}
	return config, nil
}
