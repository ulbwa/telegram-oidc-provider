package cache

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
)

type RedisTokenVerificationCache struct {
	redis       *redis.Client
	cachePrefix string
	cacheExp    time.Duration
	cacheSecret string
}

var _ service.TelegramTokenVerificationCache = (*RedisTokenVerificationCache)(nil)

func NewRedisTokenVerificationCache(redisClient *redis.Client, cachePrefix string, cacheExp time.Duration, cacheSecret string) (*RedisTokenVerificationCache, error) {
	if redisClient == nil {
		return nil, errors.New("redis client cannot be nil")
	}
	return &RedisTokenVerificationCache{
		redis:       redisClient,
		cachePrefix: cachePrefix,
		cacheExp:    cacheExp,
		cacheSecret: cacheSecret,
	}, nil
}

func (c *RedisTokenVerificationCache) getCacheKey(token string) string {
	h := hmac.New(sha256.New, []byte(c.cacheSecret))
	h.Write([]byte(token))
	hashString := hex.EncodeToString(h.Sum(nil))
	return c.cachePrefix + hashString
}

func (c *RedisTokenVerificationCache) GetTokenStatus(ctx context.Context, token string) (bool, error) {
	key := c.getCacheKey(token)
	log := zerolog.Ctx(ctx).With().Str("service", "redisTokenVerificationCache").Str("cache_key", key).Logger()
	val, err := c.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Debug().Msg("token not found in cache")
		return false, service.ErrTokenNotInCache
	}
	if err != nil {
		log.Err(err).Msg("failed to get token status from cache")
		return false, err
	}

	isValid := val == "1"
	log.Debug().Bool("is_valid", isValid).Msg("token status retrieved from cache")
	return isValid, nil
}

func (c *RedisTokenVerificationCache) CacheTokenValid(ctx context.Context, token string) error {
	key := c.getCacheKey(token)
	log := zerolog.Ctx(ctx).With().Str("service", "redisTokenVerificationCache").Str("cache_key", key).Logger()
	log.Debug().Msg("caching valid token")
	return c.redis.Set(ctx, key, "1", c.cacheExp).Err()
}

func (c *RedisTokenVerificationCache) CacheTokenInvalid(ctx context.Context, token string) error {
	key := c.getCacheKey(token)
	log := zerolog.Ctx(ctx).With().Str("service", "redisTokenVerificationCache").Str("cache_key", key).Logger()
	log.Debug().Msg("caching invalid token")
	return c.redis.Set(ctx, key, "0", c.cacheExp).Err()
}
