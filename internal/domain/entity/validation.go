package entity

import (
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
)

var ErrInvariantCheckFailed = errors.New("invariant check failed")

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if strings.TrimSpace(username) != username {
		return fmt.Errorf("username contains leading or trailing whitespace: %w", ErrInvariantCheckFailed)
	}
	for _, r := range username {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return fmt.Errorf("username contains invalid characters: %w", ErrInvariantCheckFailed)
		}
	}
	return nil
}

func validateBotId(id int64) error {
	if id <= 0 {
		return fmt.Errorf("bot id must be positive: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateBotName(name string) error {
	if name == "" {
		return fmt.Errorf("bot name cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if strings.TrimSpace(name) != name {
		return fmt.Errorf("bot name contains leading or trailing whitespace: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateBotToken(token string) error {
	if token == "" {
		return fmt.Errorf("bot token cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if strings.TrimSpace(token) != token {
		return fmt.Errorf("bot token contains leading or trailing whitespace: %w", ErrInvariantCheckFailed)
	}

	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return fmt.Errorf("bot token must have format <bot_id>:<token_string>: %w", ErrInvariantCheckFailed)
	}

	botId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("bot token: bot_id must be a valid integer: %w", ErrInvariantCheckFailed)
	}
	if botId <= 0 {
		return fmt.Errorf("bot token: bot_id must be positive: %w", ErrInvariantCheckFailed)
	}

	tokenString := parts[1]
	if tokenString == "" {
		return fmt.Errorf("bot token: token_string cannot be empty: %w", ErrInvariantCheckFailed)
	}

	return nil
}

func validateBotUsername(username string) error {
	if err := validateUsername(username); err != nil {
		return fmt.Errorf("invalid bot username: %w", err)
	}
	return nil
}

func validateClientId(clientId string) error {
	if clientId == "" {
		return fmt.Errorf("client id cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if strings.TrimSpace(clientId) != clientId {
		return fmt.Errorf("client id contains leading or trailing whitespace: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateUserId(id int64) error {
	if id <= 0 {
		return fmt.Errorf("user id must be positive: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateUserFirstName(firstName string) error {
	if firstName == "" {
		return fmt.Errorf("first name cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if strings.TrimSpace(firstName) != firstName {
		return fmt.Errorf("first name contains leading or trailing whitespace: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateUserLastName(lastName string) error {
	if lastName == "" {
		return fmt.Errorf("last name cannot be empty: %w", ErrInvariantCheckFailed)
	}
	if strings.TrimSpace(lastName) != lastName {
		return fmt.Errorf("last name contains leading or trailing whitespace: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateUserUsername(username string) error {
	if err := validateUsername(username); err != nil {
		return fmt.Errorf("invalid user username: %w", err)
	}
	return nil
}

func validateIP(ip netip.Addr) error {
	if !ip.IsValid() || ip.IsUnspecified() {
		return fmt.Errorf("invalid IP address: %w", ErrInvariantCheckFailed)
	}
	return nil
}

func validateUrl(url *url.URL) error {
	if url == nil {
		return fmt.Errorf("url cannot be nil: %w", ErrInvariantCheckFailed)
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return fmt.Errorf("url scheme must be http or https: %w", ErrInvariantCheckFailed)
	}
	if url.Host == "" {
		return fmt.Errorf("url host cannot be empty: %w", ErrInvariantCheckFailed)
	}
	return nil
}
