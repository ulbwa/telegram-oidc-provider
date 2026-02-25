package service

// TelegramAuthHashVerifier verifies HMAC-SHA256 signatures of Telegram authentication data.
type TelegramAuthHashVerifier interface {
	Verify(query string, hash string, botToken string) error
}
