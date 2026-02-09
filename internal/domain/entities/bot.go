package domain

import (
	"strings"
	"time"
)

// Bot represents a Telegram bot.
type Bot struct {
	Id        int64
	Name      string
	ClientId  string
	Username  string
	Token     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// NewBot creates a new Bot instance with validation.
func NewBot(id int64, name string, clientId string, username string, token string) (*Bot, error) {
	if err := validateBotId(id); err != nil {
		return nil, err
	}
	if err := validateBotName(name); err != nil {
		return nil, err
	}
	if err := validateClientId(clientId); err != nil {
		return nil, err
	}
	if err := validateBotUsername(username); err != nil {
		return nil, err
	}
	if err := validateBotToken(token); err != nil {
		return nil, err
	}

	return &Bot{
		Id:        id,
		Name:      strings.TrimSpace(name),
		ClientId:  strings.TrimSpace(clientId),
		Username:  strings.TrimSpace(username),
		Token:     strings.TrimSpace(token),
		CreatedAt: time.Now(),
	}, nil
}

func (b *Bot) ModifiedAt() time.Time {
	if b.UpdatedAt != nil {
		return *b.UpdatedAt
	}
	return b.CreatedAt
}

func (b *Bot) SetName(name string) error {
	if err := validateBotName(name); err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	if b.Name == name {
		return nil
	}
	currentTime := time.Now()
	b.Name = name
	b.UpdatedAt = &currentTime
	return nil
}

func (b *Bot) SetToken(token string) error {
	if err := validateBotToken(token); err != nil {
		return err
	}
	token = strings.TrimSpace(token)
	if b.Token == token {
		return nil
	}
	currentTime := time.Now()
	b.Token = token
	b.UpdatedAt = &currentTime
	return nil
}
