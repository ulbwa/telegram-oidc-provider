package web

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
)

func (s *server) Login(c echo.Context) error {
	output, err := s.resolveLoginChallengeUsecase.Execute(c.Request().Context(), &usecase.ResolveLoginChallengeInput{
		LoginChallenge: c.QueryParam("login_challenge"),
	})
	if err != nil {
		zerolog.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to execute login usecase")

		if errors.Is(err, usecase.ErrInvalidInput) {
			var objectNotFoundErr *usecase.ObjectNotFoundErr
			if errors.As(err, &objectNotFoundErr) && objectNotFoundErr.Object == "client" {
				return s.fallbackToErrorPage(c, ErrCodeInvalidClient)
			}

			var objectInvalidErr *usecase.ObjectInvalidErr
			if errors.As(err, &objectInvalidErr) && objectInvalidErr.Object == "bot" && objectInvalidErr.Field == "token" {
				return s.fallbackToErrorPage(c, ErrCodeInvalidBotCredentials)
			}

			return s.fallbackToErrorPage(c, ErrCodeInvalidRequest)
		}

		return s.fallbackToErrorPage(c, ErrCodeInternalError)
	}

	if output == nil {
		return s.fallbackToErrorPage(c, ErrCodeInternalError)
	}

	if output.Action == usecase.ResolveLoginChallengeActionRedirect {
		if output.RedirectUri == nil {
			return s.fallbackToErrorPage(c, ErrCodeInternalError)
		}
		return c.Redirect(http.StatusFound, *output.RedirectUri)
	}

	if output.WidgetUri == nil || output.MiniAppCallbackUri == nil {
		return s.fallbackToErrorPage(c, ErrCodeInternalError)
	}

	return c.Render(http.StatusOK, "login", map[string]any{
		"WidgetUri":          *output.WidgetUri,
		"MiniAppCallbackUri": *output.MiniAppCallbackUri,
	})
}
