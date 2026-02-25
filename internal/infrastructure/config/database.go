package config

// DatabaseConfig represents database configuration.
type DatabaseConfig struct {
	DSN URL `yaml:"dsn" validate:"required"` // Database connection string
}
