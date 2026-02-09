package entities

import (
	"net"
	"time"
)

// UserBotLogin represents a user's login session with a specific bot.
type UserBotLogin struct {
	UserId      int64
	BotId       int64
	IP          net.IP
	UserAgent   *string
	Language    *string
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// NewUserBotLogin creates a new UserBotLogin instance with validation.
func NewUserBotLogin(userId, botId int64, ip net.IP, userAgent, language *string) (*UserBotLogin, error) {
	if err := validateLoginUserId(userId); err != nil {
		return nil, err
	}
	if err := validateLoginBotId(botId); err != nil {
		return nil, err
	}
	if err := validateIP(ip); err != nil {
		return nil, err
	}

	now := time.Now()
	return &UserBotLogin{
		UserId:      userId,
		BotId:       botId,
		IP:          ip,
		UserAgent:   userAgent,
		Language:    language,
		LastLoginAt: now,
		CreatedAt:   now,
	}, nil
}

func (ubl *UserBotLogin) ModifiedAt() time.Time {
	if ubl.UpdatedAt != nil {
		return *ubl.UpdatedAt
	}
	return ubl.CreatedAt
}

func (ubl *UserBotLogin) Touch() {
	now := time.Now()
	ubl.UpdatedAt = &now
}

func (ubl *UserBotLogin) SetIP(ip net.IP) error {
	if ubl.IP.Equal(ip) {
		return nil
	}
	if err := validateIP(ip); err != nil {
		return err
	}
	ubl.IP = ip
	ubl.Touch()
	return nil
}

func (ubl *UserBotLogin) SetUserAgent(userAgent *string) {
	if (ubl.UserAgent == nil && userAgent == nil) || (ubl.UserAgent != nil && userAgent != nil && *ubl.UserAgent == *userAgent) {
		return
	}
	ubl.UserAgent = userAgent
	ubl.Touch()
}

func (ubl *UserBotLogin) SetLanguage(language *string) {
	if (ubl.Language == nil && language == nil) || (ubl.Language != nil && language != nil && *ubl.Language == *language) {
		return
	}
	ubl.Language = language
	ubl.Touch()
}

func (ubl *UserBotLogin) UpdateLastLogin() {
	ubl.LastLoginAt = time.Now()
	ubl.Touch()
}
