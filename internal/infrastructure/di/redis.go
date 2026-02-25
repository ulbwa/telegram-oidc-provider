package di

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
)

// redisConnection wraps redis.Client and provides graceful shutdown.
type redisConnection struct {
	client *redis.Client
}

// Shutdown closes Redis connection.
func (r *redisConnection) Shutdown(ctx context.Context) error {
	return r.client.Close()
}

func provideRedis(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (*redisConnection, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			return nil, err
		}

		return &redisConnection{client: client}, nil
	})

	do.Provide(injector, func(i do.Injector) (*redis.Client, error) {
		conn, err := do.Invoke[*redisConnection](i)
		if err != nil {
			return nil, err
		}
		return conn.client, nil
	})
}
