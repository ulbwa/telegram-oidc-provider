package service

import (
	"errors"
	"net/url"
	"time"
)

// ErrInvalidTelegramAuthData is returned when Telegram authentication data is invalid
var ErrInvalidTelegramAuthData = errors.New("invalid Telegram authentication data")

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

// TelegramWidgetDataParser parses Telegram Widget authentication data.
type TelegramWidgetDataParser interface {
	Parse(params map[string]any) (*TelegramAuthData, error)
}

// TelegramMiniAppDataParser parses Telegram Mini App authentication data.
type TelegramMiniAppDataParser interface {
	Parse(params map[string]any) (*TelegramAuthData, error)
}
