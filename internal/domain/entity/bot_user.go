package entity

import (
	"fmt"
	"net/netip"
	"time"
)

type BotUser struct {
	BotId       int64
	UserId      int64
	User        User
	IP          netip.Addr
	UserAgent   *string
	Language    *string
	LastLoginAt time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewBotUser(botId, userId int64, user *User, ip netip.Addr, userAgent *string, language *string) (*BotUser, error) {
	if err := validateBotId(botId); err != nil {
		return nil, err
	}
	if err := validateUserId(userId); err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil: %w", ErrInvariantCheckFailed)
	}
	if *user == (User{}) {
		return nil, fmt.Errorf("user cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if err := validateIP(ip); err != nil {
		return nil, err
	}

	return &BotUser{
		BotId:       botId,
		UserId:      userId,
		User:        *user,
		IP:          ip,
		UserAgent:   userAgent,
		Language:    language,
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}, nil
}

func (u *BotUser) ModifiedAt() time.Time {
	if u.UpdatedAt != nil {
		return *u.UpdatedAt
	}
	return u.CreatedAt
}

func (u *BotUser) Touch() {
	now := time.Now()
	u.UpdatedAt = &now
}

func (u *BotUser) SetUser(user *User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil: %w", ErrInvariantCheckFailed)
	}
	if *user == (User{}) {
		return fmt.Errorf("user cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if u.User == *user {
		return nil
	}
	u.User = *user
	u.Touch()
	return nil
}

func (u *BotUser) SetIP(ip netip.Addr) error {
	if u.IP == ip {
		return nil
	}
	if err := validateIP(ip); err != nil {
		return err
	}
	u.IP = ip
	u.Touch()
	return nil
}

func (u *BotUser) SetUserAgent(userAgent *string) {
	if (u.UserAgent == nil && userAgent == nil) || (u.UserAgent != nil && userAgent != nil && *u.UserAgent == *userAgent) {
		return
	}
	u.UserAgent = userAgent
	u.Touch()
}

func (u *BotUser) SetLanguage(language *string) {
	if (u.Language == nil && language == nil) || (u.Language != nil && language != nil && *u.Language == *language) {
		return
	}
	u.Language = language
	u.Touch()
}

func (u *BotUser) UpdateLastLogin() {
	u.LastLoginAt = time.Now()
	u.Touch()
}
