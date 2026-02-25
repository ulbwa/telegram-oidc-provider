package config

// HTTPServerConfig represents HTTP server configuration.
type HTTPServerConfig struct {
	Address         string `yaml:"address"           validate:"required"`
	BaseUri         URL    `yaml:"base_uri"`
	TelegramAuthURI URL    `yaml:"telegram_auth_uri" validate:"required"`
}
