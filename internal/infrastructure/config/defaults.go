package config

import "time"

const (
	defaultTelegramAuthURI              = "https://oauth.telegram.org/auth"
	defaultTelegramAuthDataTTL          = 5 * time.Minute
	defaultTelegramTokenVerificationTTL = 5 * time.Minute
	defaultTelegramReplayGuardTTL       = 5 * time.Minute
)

var defaultConfig = Config{
	HTTPServer: HTTPServerConfig{
		TelegramAuthURI: MustParseURL(defaultTelegramAuthURI),
	},
	Database: DatabaseConfig{
		DSN: MustParseURL("postgres://localhost:5432/dbname?sslmode=disable"),
	},
	Redis: RedisConfig{
		Addr: "127.0.0.1:6379",
		DB:   0,
	},
	Security: SecurityConfig{
		Telegram: TelegramSecurityConfig{
			AuthDataTTLSeconds: defaultTelegramAuthDataTTL,
			TokenVerificationCache: TelegramTokenVerificationCacheConfig{
				TTL: defaultTelegramTokenVerificationTTL,
			},
			ReplayGuard: TelegramReplayGuardConfig{
				TTL: defaultTelegramReplayGuardTTL,
			},
		},
	},
}
