package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main application configuration.
type Config struct {
	HTTPServer HTTPServerConfig `yaml:"http_server"`
	Logger     LoggerConfig     `yaml:"logger"`
}

// HTTPServerConfig represents HTTP server configuration.
type HTTPServerConfig struct {
	Address string `yaml:"address"`
}

// LoggerConfig represents logging configuration.
type LoggerConfig struct {
	Level      string           `yaml:"level"`       // One of: trace, debug, info, warn, error, fatal, panic
	TimeFormat string           `yaml:"time_format"` // One of: unix, unixms, unixmicro, rfc3339, rfc3339nano
	Console    ConsoleLogConfig `yaml:"console"`
	Files      []FileLogConfig  `yaml:"files"`
}

// ConsoleLogConfig represents console output configuration.
type ConsoleLogConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Level    string `yaml:"level"`     // If empty, uses global level
	MaxLevel string `yaml:"max_level"` // Optional, for filtering
	Colored  bool   `yaml:"colored"`
	Pretty   bool   `yaml:"pretty"` // If true, pretty-print logs instead of JSON
}

// FileLogConfig represents file output configuration.
type FileLogConfig struct {
	Path     string       `yaml:"path"`
	Level    string       `yaml:"level"`     // If empty, uses global level
	MaxLevel string       `yaml:"max_level"` // Optional, for filtering
	Rotate   RotateConfig `yaml:"rotate"`
}

// RotateConfig represents log rotation configuration.
type RotateConfig struct {
	Enabled    bool `yaml:"enabled"`
	MaxSize    int  `yaml:"max_size"` // In megabytes
	MaxAge     int  `yaml:"max_age"`  // In days
	MaxBackups int  `yaml:"max_backups"`
	Compress   bool `yaml:"compress"` // Compress old files
}

// ReadConfig reads configuration from YAML file.
func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
