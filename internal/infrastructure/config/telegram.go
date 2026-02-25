package config

import "time"

// TelegramTokenVerificationCacheConfig holds token verification cache settings.
type TelegramTokenVerificationCacheConfig struct {
	Prefix string        `yaml:"prefix" validate:"required"`
	TTL    time.Duration `yaml:"ttl"    validate:"required,gt=0"`
	Secret string        `yaml:"secret" validate:"required"`
}

// TelegramReplayGuardConfig holds replay guard cache settings.
type TelegramReplayGuardConfig struct {
	Prefix string        `yaml:"prefix" validate:"required"`
	TTL    time.Duration `yaml:"ttl"    validate:"required,gt=0"`
}

// TelegramSecurityConfig holds Telegram-related security settings.
type TelegramSecurityConfig struct {
	AuthDataTTLSeconds     time.Duration                        `yaml:"auth_data_ttl"            validate:"required,gt=0"`
	TokenVerificationCache TelegramTokenVerificationCacheConfig `yaml:"token_verification_cache" validate:"required"`
	ReplayGuard            TelegramReplayGuardConfig            `yaml:"replay_guard"             validate:"required"`
}
