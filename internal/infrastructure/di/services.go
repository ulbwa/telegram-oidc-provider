package di

import (
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/cache"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/telegram"
)

func provideServices(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (service.TelegramTokenVerificationCache, error) {
		redisClient, err := do.Invoke[*redis.Client](i)
		if err != nil {
			return nil, err
		}

		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		cacheCfg := cfg.Security.Telegram.TokenVerificationCache
		return cache.NewRedisTokenVerificationCache(
			redisClient,
			cacheCfg.Prefix,
			cacheCfg.TTL,
			cacheCfg.Secret,
		)
	})

	do.Provide(injector, func(i do.Injector) (service.TelegramTokenVerifier, error) {
		tokenCache, err := do.Invoke[service.TelegramTokenVerificationCache](i)
		if err != nil {
			return nil, err
		}

		return telegram.NewTelegramTokenVerifier(tokenCache), nil
	})

	do.Provide(injector, func(i do.Injector) (service.TelegramWidgetDataParser, error) {
		return telegram.NewTelegramWidgetDataParser(), nil
	})

	do.Provide(injector, func(i do.Injector) (service.TelegramMiniAppDataParser, error) {
		return telegram.NewTelegramWidgetDataParser(), nil
	})

	do.Provide(injector, func(i do.Injector) (service.TelegramAuthHashVerifier, error) {
		return telegram.NewTelegramAuthHashVerifier(), nil
	})

	do.Provide(injector, func(i do.Injector) (service.TelegramReplayGuard, error) {
		redisClient, err := do.Invoke[*redis.Client](i)
		if err != nil {
			return nil, err
		}

		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		replayCfg := cfg.Security.Telegram.ReplayGuard
		return telegram.NewRedisTelegramReplayGuard(redisClient, replayCfg.Prefix)
	})
}
