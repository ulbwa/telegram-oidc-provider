package entity

import "time"

// Bot represents a Telegram bot.
type Bot struct {
	Id        int64
	Name      string
	ClientId  *string
	Username  string
	Token     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewBot(id int64, name string, username string, token string) (*Bot, error) {
	if err := validateBotName(name); err != nil {
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

func (b *Bot) SetUsername(username string) error {
	if err := validateBotUsername(username); err != nil {
		return err
	}
	if b.Username == username {
		return nil
	}
	b.Username = username
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

func (b *Bot) SetClientId(clientId *string) error {
	if clientId != nil {
		if err := validateClientId(*clientId); err != nil {
			return err
		}
	}
	if b.ClientId == clientId {
		return nil
	}
	b.ClientId = clientId
	b.Touch()
	return nil
}
