package config

// SecurityBotTokenConfig represents bot token security settings.
type SecurityBotTokenConfig struct {
	EncryptionKey string `yaml:"encryption_key" validate:"required"`
}

// SecurityConfig represents application security configuration.
type SecurityConfig struct {
	BotToken SecurityBotTokenConfig `yaml:"bot_token" validate:"required"`
	Telegram TelegramSecurityConfig `yaml:"telegram"  validate:"required"`
}
