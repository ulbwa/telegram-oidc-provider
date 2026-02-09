package domain

import (
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
		Name:      name,
		ClientId:  clientId,
		Username:  username,
		Token:     token,
		CreatedAt: time.Now(),
	}, nil
}

func (b *Bot) ModifiedAt() time.Time {
	if b.UpdatedAt != nil {
		return *b.UpdatedAt
	}
	return b.CreatedAt
}

func (b *Bot) Touch() {
	now := time.Now()
	b.UpdatedAt = &now
}

func (b *Bot) SetName(name string) error {
	if err := validateBotName(name); err != nil {
		return err
	}
	if b.Name == name {
		return nil
	}
	b.Name = name
	b.Touch()
	return nil
}

func (b *Bot) SetToken(token string) error {
	if err := validateBotToken(token); err != nil {
		return err
	}
	if b.Token == token {
		return nil
	}
	b.Token = token
	b.Touch()
	return nil
}
