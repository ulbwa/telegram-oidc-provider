package service

import (
	"context"
	"errors"
)

var (
	// ErrTokenNotInCache is returned when token is not found in cache
	ErrTokenNotInCache = errors.New("token not found in cache")

	// ErrTelegramBotTokenMalformed is returned when Telegram bot token format is invalid
	ErrTelegramBotTokenMalformed = errors.New("invalid Telegram bot token format")

	// ErrTelegramBotTokenInvalid is returned when Telegram bot token is invalid
	ErrTelegramBotTokenInvalid = errors.New("invalid Telegram bot token")
)

type TelegramBotInfo struct {
	Id       int64
	Name     string
	Username string
}

// TelegramTokenVerifier verifies the validity of a Telegram bot token.
type TelegramTokenVerifier interface {
	Verify(ctx context.Context, token string, opts *VerifyOptions) (*TelegramBotInfo, error)
}

// TelegramTokenVerificationCache stores and retrieves verification results for bot tokens.
type TelegramTokenVerificationCache interface {
	GetTokenStatus(ctx context.Context, token string) (bool, error)
	CacheTokenValid(ctx context.Context, token string) error
	CacheTokenInvalid(ctx context.Context, token string) error
}

// VerifyOptions contains options for token verification.
type VerifyOptions struct {
	SkipCacheRead bool
}

// VerifyOption is a functional option for VerifyOptions.
type VerifyOption func(*VerifyOptions)

// WithSkipCacheRead returns an option that skips cache read.
func WithSkipCacheRead() VerifyOption {
	return func(opts *VerifyOptions) {
		opts.SkipCacheRead = true
	}
}

// NewVerifyOptions creates default VerifyOptions and applies provided options.
func NewVerifyOptions(options ...VerifyOption) *VerifyOptions {
	opts := &VerifyOptions{
		SkipCacheRead: false,
	}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}
