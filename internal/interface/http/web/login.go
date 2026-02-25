package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
)

func (s *server) Login(c echo.Context) error {
	loginChallenge := c.QueryParam("login_challenge")
	if loginChallenge == "" {
		return s.fallbackToErrorPage(c, "invalid_request")
	}

	loginInfo, _, err := s.hydra.AdminApi.GetLoginRequest(c.Request().Context()).LoginChallenge(loginChallenge).Execute()
	if err != nil {
		zerolog.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to get login request from hydra")
		return s.fallbackToErrorPage(c, ErrCodeInvalidRequest)
	}

	var bot entity.Bot
	if err := s.botRepo.GetByClientID(c.Request().Context(), *loginInfo.Client.ClientId, &bot); err != nil {
		zerolog.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to get bot by client id")
		return s.fallbackToErrorPage(c, ErrCodeInvalidClient)
	}

	if _, err := s.botVerifier.Verify(c.Request().Context(), bot.Token, nil); err != nil {
		zerolog.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to verify bot token")
		if errors.Is(err, service.ErrTelegramBotTokenInvalid) {
			return s.fallbackToErrorPage(c, ErrCodeInvalidBotCredentials)
		} else {
			return s.fallbackToErrorPage(c, ErrCodeInternalError)
		}
	}

	origin := s.baseUri
	origin = origin.JoinPath("/login")

	widgetCallbackUri := s.baseUri
	widgetCallbackUri = origin.JoinPath("/widget/callback")
	widgetCallbackUriQuery := widgetCallbackUri.Query()
	widgetCallbackUriQuery.Set("login_challenge", loginChallenge)
	widgetCallbackUri.RawQuery = widgetCallbackUriQuery.Encode()

	widgetUri := s.telegramAuthUri
	widgetUriQuery := widgetUri.Query()
	widgetUriQuery.Set("bot_id", fmt.Sprintf("%d", bot.Id))
	widgetUriQuery.Set("origin", origin.String())
	widgetUriQuery.Set("request_access", "write")
	widgetUriQuery.Set("return_to", widgetCallbackUri.String())
	widgetUri.RawQuery = widgetUriQuery.Encode()

	miniappCallbackUri := s.baseUri
	miniappCallbackUri = origin.JoinPath("/miniapp/callback")
	miniappCallbackUriQuery := miniappCallbackUri.Query()
	miniappCallbackUriQuery.Set("login_challenge", loginChallenge)
	miniappCallbackUri.RawQuery = miniappCallbackUriQuery.Encode()

	return c.Render(http.StatusOK, "login", map[string]any{
		"WidgetUri":          widgetUri.String(),
		"MiniAppCallbackUri": miniappCallbackUri.String(),
	})
}
