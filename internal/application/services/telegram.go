package services

import (
	"context"
	"net/url"
	"time"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

// TelegramUserData represents the user data received from Telegram Mini Apps or Widgets.
type TelegramUserData struct {
	Id           int64
	FirstName    string
	LastName     *string
	Username     *string
	LanguageCode *string
	PhotoUrl     *url.URL
	IsPremium    *bool
}

// TelegramAuthData contains Telegram authentication data with signature verification.
type TelegramAuthData struct {
	Raw      string // Query string without "hash" parameter
	Hash     string // HMAC-SHA256 signature signed with bot token
	User     *TelegramUserData
	AuthDate time.Time
}

func (d *TelegramAuthData) IsExpired(ttl time.Duration) bool {
	return d.AuthDate.Add(ttl).Before(time.Now())
}

// TelegramHashVerifier verifies HMAC-SHA256 signatures of Telegram authentication data.
type TelegramHashVerifier interface {
	Verify(query string, botToken string) error
}

// TelegramWidgetDataParser parses authentication data from Telegram Login Widget.
type TelegramWidgetDataParser interface {
	Parse(query string) (*TelegramAuthData, error)
}

// TelegramInitDataParser parses authentication data from Telegram Mini Apps initData.
type TelegramInitDataParser interface {
	Parse(query string) (*TelegramAuthData, error)
}

// TelegramReplayGuard prevents replay attacks by tracking used authentication hashes.
type TelegramReplayGuard interface {
	// CheckAndMarkUsed verifies the hash hasn't been used and marks it as used atomically.
	// Returns error if hash was already used or if operation fails.
	CheckAndMarkUsed(ctx context.Context, hash string, ttl time.Duration) error
}

type TelegramUserFactory interface {
	CreateUser(data *TelegramUserData) (*domain.User, error)
}

type TelegramBotData struct {
	Id       int64
	Name     string
	Username string
}

// TelegramBotVerifier validates bot tokens by querying Telegram Bot API.
type TelegramBotVerifier interface {
	Verify(ctx context.Context, botToken string) (*TelegramBotData, error)
}
