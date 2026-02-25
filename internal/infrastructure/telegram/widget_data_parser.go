package telegram

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
)

type DefaultTelegramWidgetDataParser struct{}

var _ service.TelegramWidgetDataParser = (*DefaultTelegramWidgetDataParser)(nil)

func NewTelegramWidgetDataParser() *DefaultTelegramWidgetDataParser {
	return &DefaultTelegramWidgetDataParser{}
}

func (p *DefaultTelegramWidgetDataParser) Parse(params map[string]any) (*service.TelegramAuthData, error) {
	var output service.TelegramAuthData

	hash, ok := params["hash"].(string)
	if !ok || hash == "" {
		return nil, fmt.Errorf("invalid 'hash' parameter: %w", service.ErrInvalidTelegramAuthData)
	}
	output.Hash = hash

	// Build raw string without hash (verifier will sort keys for verification)
	var rawParts []string
	for key, value := range params {
		if key != "hash" {
			valueStr := fmt.Sprintf("%v", value)
			rawParts = append(rawParts, url.QueryEscape(key)+"="+url.QueryEscape(valueStr))
		}
	}
	output.Raw = strings.Join(rawParts, "&")

	authDateInt64, err := parseIntField(params, "auth_date")
	if err != nil {
		return nil, fmt.Errorf("invalid 'auth_date': %w", err)
	}
	output.AuthDate = time.Unix(authDateInt64, 0)

	user := &service.TelegramUserData{}

	userId, err := parseIntField(params, "id")
	if err != nil {
		return nil, fmt.Errorf("invalid 'id': %w", err)
	}
	user.Id = userId

	firstName, ok := params["first_name"].(string)
	if !ok || firstName == "" {
		return nil, fmt.Errorf("invalid 'first_name': %w", service.ErrInvalidTelegramAuthData)
	}
	user.FirstName = firstName

	if lastName, ok := params["last_name"].(string); ok && lastName != "" {
		user.LastName = &lastName
	}

	if username, ok := params["username"].(string); ok && username != "" {
		user.Username = &username
	}

	if photoUrlStr, ok := params["photo_url"].(string); ok && photoUrlStr != "" {
		photoUrl, err := url.Parse(photoUrlStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'photo_url': %w", service.ErrInvalidTelegramAuthData)
		}
		user.PhotoUrl = photoUrl
	}

	output.User = user

	return &output, nil
}

func parseIntField(params map[string]any, field string) (int64, error) {
	value, ok := params[field]
	if !ok {
		return 0, fmt.Errorf("missing field: %w", service.ErrInvalidTelegramAuthData)
	}

	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse as int64: %w", service.ErrInvalidTelegramAuthData)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported type: %w", service.ErrInvalidTelegramAuthData)
	}
}
