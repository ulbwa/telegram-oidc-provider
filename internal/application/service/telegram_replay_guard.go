package service

import (
	"context"
	"errors"
	"time"
)

// ErrReplayDetected is returned when a replay attack is detected
var ErrReplayDetected = errors.New("replay detected")

// TelegramReplayGuard prevents replay attacks by tracking used authentication hashes.
type TelegramReplayGuard interface {
	CheckAndMarkUsed(ctx context.Context, hash string, ttl time.Duration) error
}
