package api

import (
	"context"
	"net/netip"
	"strings"

	"github.com/ulbwa/telegram-oidc-provider/api/generated"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	xlanguage "golang.org/x/text/language"
)

var defaultClientIP = netip.MustParseAddr("127.0.0.1")

func normalizeBCP47Language(value *string) string {
	if value == nil {
		return ""
	}

	raw := strings.TrimSpace(*value)
	if raw == "" {
		return ""
	}

	tags, _, err := xlanguage.ParseAcceptLanguage(raw)
	if err == nil && len(tags) > 0 {
		return tags[0].String()
	}

	tag, err := xlanguage.Parse(raw)
	if err != nil {
		return ""
	}

	return tag.String()
}

func normalizeBCP47LanguagePtr(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := normalizeBCP47Language(value)
	if normalized == "" {
		return nil
	}

	return &normalized
}

// Login user by telegram widget auth data
// (GET /widget/callback)
func (s *server) GetWidgetCallback(ctx context.Context, request generated.GetWidgetCallbackRequestObject) (generated.GetWidgetCallbackResponseObject, error) {
	input := usecase.LoginByWidgetInput{
		LoginChallenge: request.Params.LoginChallenge,
		AuthData:       request.Params.TelegramWidgetAuthData,
		UserAgent:      request.Params.UserAgent,
		Language:       normalizeBCP47LanguagePtr(request.Params.AcceptLanguage),
		ClientIP:       defaultClientIP,
	}

	output, err := s.loginByWidget.Execute(ctx, &input)
	if err != nil {
		return nil, err
	}

	var resp generated.GetWidgetCallback203Response
	resp.Headers.Location = output.RedirectUri
	return resp, nil
}
