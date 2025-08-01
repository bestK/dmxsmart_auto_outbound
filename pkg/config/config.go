package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration structure
type ConfigStruct struct {
	Account     string   `yaml:"account"`
	Password    string   `yaml:"password"`
	AccessToken string   `yaml:"access_token"`
	WarehouseID string   `yaml:"warehouse_id"`
	CustomerIDs []string `yaml:"customer_ids"`
	OcrEndpoint string   `yaml:"ocr_endpoint"`
	Debug       bool     `yaml:"debug"`
}

var Config *ConfigStruct
var ConfigPath string

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*ConfigStruct, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	ConfigPath = filename

	config := &ConfigStruct{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	Config = config

	return config, nil
}

// SaveConfig saves configuration to a YAML file
func SaveConfig() error {
	data, err := yaml.Marshal(Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
