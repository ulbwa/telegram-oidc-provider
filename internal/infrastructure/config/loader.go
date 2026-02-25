package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var configValidator = validator.New()

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := defaultConfig
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	if err := configValidator.Struct(c); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &c, nil
}
