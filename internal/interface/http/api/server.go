package api

import (
	"errors"
	"net/url"

	"github.com/ulbwa/telegram-oidc-provider/api/generated"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
)

type server struct {
	baseUri *url.URL
	syncBot *usecase.SyncBot
}

var _ generated.StrictServerInterface = (*server)(nil)

func NewServer(
	baseUri *url.URL,
	syncBot *usecase.SyncBot,
) (generated.StrictServerInterface, error) {
	if baseUri == nil {
		return nil, errors.New("baseUri cannot be nil")
	}
	if baseUri.Scheme == "" || baseUri.Host == "" {
		return nil, errors.New("baseUri must have scheme and host")
	}

	return &server{
		baseUri: baseUri,
		syncBot: syncBot,
	}, nil
}
