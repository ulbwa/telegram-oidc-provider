package telegram

import (
	"context"
	"errors"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
)

type DefaultTelegramTokenVerifier struct {
	tokenCache service.TelegramTokenVerificationCache
}

var _ service.TelegramTokenVerifier = (*DefaultTelegramTokenVerifier)(nil)

func NewTelegramTokenVerifier(tokenCache service.TelegramTokenVerificationCache) *DefaultTelegramTokenVerifier {
	return &DefaultTelegramTokenVerifier{
		tokenCache: tokenCache,
	}
}

func (s *DefaultTelegramTokenVerifier) getLogger(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx).With().Str("service", "defaultTelegramTokenVerifier").Logger()
	return &logger
}

func (s *DefaultTelegramTokenVerifier) cacheTokenInvalid(token string) {
	if s.tokenCache == nil {
		return
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.tokenCache.CacheTokenInvalid(bgCtx, token); err != nil {
			s.getLogger(bgCtx).Err(err).Msg("failed to cache invalid token")
		}
	}()
}

func (s *DefaultTelegramTokenVerifier) cacheTokenValid(token string) {
	if s.tokenCache == nil {
		return
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.tokenCache.CacheTokenValid(bgCtx, token); err != nil {
			s.getLogger(bgCtx).Err(err).Msg("failed to cache valid token")
		}
	}()
}

func (s *DefaultTelegramTokenVerifier) Verify(ctx context.Context, token string, opts *service.VerifyOptions) (*service.TelegramBotInfo, error) {
	if opts == nil {
		opts = service.NewVerifyOptions()
	}

	log := s.getLogger(ctx).With().Interface("options", opts).Logger()

	if !opts.SkipCacheRead && s.tokenCache != nil {
		isValid, err := s.tokenCache.GetTokenStatus(ctx, token)
		if err == nil {
			if !isValid {
				return nil, service.ErrTelegramBotTokenInvalid
			}
			return &service.TelegramBotInfo{}, nil
		}
		if err != service.ErrTokenNotInCache {
			log.Err(err).Msg("failed to check cache, continuing to Telegram API")
		}
	}

	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{DisableTokenCheck: true})
	if err != nil {
		log.Err(err).Msg("failed to create Telegram bot with provided token")
		s.cacheTokenInvalid(token)

		return nil, service.ErrTelegramBotTokenMalformed
	}

	me, err := bot.GetMe(nil)
	if err != nil {
		log.Err(err).Msg("failed to call GetMe with provided token, token is likely invalid")
		s.cacheTokenInvalid(token)

		return nil, service.ErrTelegramBotTokenInvalid
	}
	if !me.IsBot {
		return nil, errors.New("provided token does not belong to a bot")
	}

	var botInfo service.TelegramBotInfo
	botInfo.Id = me.Id
	botInfo.Name = me.FirstName
	botInfo.Username = me.Username
	s.cacheTokenValid(token)

	return &botInfo, nil
}
