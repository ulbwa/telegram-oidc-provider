package telegram

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
)

type RedisTelegramReplayGuard struct {
	redis  *redis.Client
	prefix string
}

var _ service.TelegramReplayGuard = (*RedisTelegramReplayGuard)(nil)

func NewRedisTelegramReplayGuard(redisClient *redis.Client, prefix string) (*RedisTelegramReplayGuard, error) {
	if redisClient == nil {
		return nil, errors.New("redis client cannot be nil")
	}
	return &RedisTelegramReplayGuard{
		redis:  redisClient,
		prefix: prefix,
	}, nil
}
func (g *RedisTelegramReplayGuard) getKey(hash string) string {
	return g.prefix + hash
}

func (g *RedisTelegramReplayGuard) CheckAndMarkUsed(ctx context.Context, hash string, ttl time.Duration) error {
	key := g.getKey(hash)
	log := zerolog.Ctx(ctx).With().Str("service", "redisTelegramReplayGuard").Str("key", key).Logger()
	set, err := g.redis.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		log.Err(err).Msg("failed to set key in redis")
		return err
	}

	// Если значение уже существовало, то это повторное использование (replay)
	if !set {
		log.Debug().Msg("replay detected: hash already used")
		return service.ErrReplayDetected
	}

	// Успешно установлено, значение не существовало
	return nil
}
