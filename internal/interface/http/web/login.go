package web

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
)

func (s *server) Login(c echo.Context) error {
	input := usecase.ResolveLoginChallengeInput{
		LoginChallenge: c.QueryParam("login_challenge"),
	}
	output, err := s.resolveLoginChallengeUsecase.Execute(c.Request().Context(), &input)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			return s.fallbackToErrorPage(c, ErrCodeInvalidRequest)
		}

		return s.fallbackToErrorPage(c, ErrCodeInternalError)
	}

	switch output.Action {
	case usecase.ResolveLoginChallengeActionRedirect:
		return c.Redirect(http.StatusFound, *output.RedirectUri)
	case usecase.ResolveLoginChallengeActionRender:
		return c.Render(http.StatusOK, "login", map[string]any{
			"WidgetUri":          *output.WidgetUri,
			"MiniAppCallbackUri": *output.MiniAppCallbackUri,
		})
	default:
		return s.fallbackToErrorPage(c, ErrCodeInternalError)
	}
}
