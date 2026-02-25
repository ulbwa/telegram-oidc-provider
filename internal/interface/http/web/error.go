package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorCode string

var (
	ErrCodeInternalError         ErrorCode = "internal_error"
	ErrCodeInvalidRequest        ErrorCode = "invalid_request"
	ErrCodeInvalidClient         ErrorCode = "invalid_client"
	ErrCodeInvalidBotCredentials ErrorCode = "invalid_bot_credentials"
)

func (s *server) fallbackToErrorPage(c echo.Context, errCode ErrorCode) error {
	uri := s.errorUri
	uriQuery := uri.Query()
	uriQuery.Set("error", string(errCode))
	uri.RawQuery = uriQuery.Encode()
	return c.Redirect(http.StatusFound, uri.String())
}

func (s *server) Error(c echo.Context) error {
	errCode := c.QueryParam("error")
	if errCode == "" {
		errCode = string(ErrCodeInternalError)
	}
	return c.Render(http.StatusOK, "error", map[string]any{
		"ErrorCode": errCode,
	})
}
