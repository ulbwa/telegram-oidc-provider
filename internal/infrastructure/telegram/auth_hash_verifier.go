package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"

	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
)

// DefaultTelegramAuthHashVerifier implements TelegramAuthHashVerifier
// using the official Telegram verification algorithm.
type DefaultTelegramAuthHashVerifier struct{}

// NewTelegramAuthHashVerifier creates a new Telegram auth hash verifier.
func NewTelegramAuthHashVerifier() *DefaultTelegramAuthHashVerifier {
	return &DefaultTelegramAuthHashVerifier{}
}

// Verify verifies the HMAC-SHA256 signature of Telegram authentication data.
// According to Telegram documentation:
// - query: URL query string without hash parameter
// - hash: the provided HMAC-SHA256 signature (hex-encoded)
// - botToken: the Telegram bot token
// - Creates a data-check-string: all fields sorted alphabetically in format "key=<value>\n"
// - Computes secret_key: SHA256(bot_token)
// - Verifies: hex(HMAC-SHA256(data_check_string, secret_key)) == provided hash
func (v *DefaultTelegramAuthHashVerifier) Verify(query string, hash string, botToken string) error {
	log := zerolog.New(zerolog.NewConsoleWriter())

	if hash == "" {
		log.Error().Msg("hash parameter is empty")
		return fmt.Errorf("%w: hash parameter missing", service.ErrInvalidTelegramAuthData)
	}

	// Parse query string using url.Values
	values, err := url.ParseQuery(query)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse query string")
		return fmt.Errorf("%w: invalid query format", service.ErrInvalidTelegramAuthData)
	}

	if len(values) == 0 {
		log.Error().Msg("no fields to verify")
		return fmt.Errorf("%w: no fields to verify", service.ErrInvalidTelegramAuthData)
	}

	// Build data-check-string: sort fields alphabetically and join with "\n"
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var dataCheckStringParts []string
	for _, key := range keys {
		// url.Values stores lists, but we expect single values per key
		// Get the first (and typically only) value
		dataCheckStringParts = append(dataCheckStringParts, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}

	dataCheckString := ""
	for i, part := range dataCheckStringParts {
		if i > 0 {
			dataCheckString += "\n"
		}
		dataCheckString += part
	}

	// Compute secret_key = SHA256(bot_token)
	tokenHash := sha256.Sum256([]byte(botToken))

	// Compute HMAC-SHA256(data_check_string, secret_key)
	h := hmac.New(sha256.New, tokenHash[:])
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	// Compare computed hash with provided hash
	if !hmac.Equal([]byte(computedHash), []byte(hash)) {
		log.Warn().
			Str("provided_hash", hash).
			Str("computed_hash", computedHash).
			Msg("hash verification failed: mismatch")
		return service.ErrInvalidTelegramAuthData
	}

	log.Debug().Msg("hash verification successful")
	return nil
}
