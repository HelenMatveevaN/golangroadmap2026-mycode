package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config описывает структуру нашего файла настроек config.yaml
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Site struct {
		Title  string `yaml:"title"`
		Author string `yaml:"author"`
	} `yaml:"site"`
}

// LoadConfig читает файл с диска и парсит его из YAML в структуру Go
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config YAML: %w", err)
	}

	return &cfg, nil
}