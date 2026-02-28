package di

import (
	"net/url"

	hydra "github.com/ory/hydra-client-go"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
)

func provideUsecases(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (*usecase.SyncBot, error) {
		transactor, err := do.Invoke[service.Transactor](i)
		if err != nil {
			return nil, err
		}

		botRepo, err := do.Invoke[repository.BotRepositoryPort](i)
		if err != nil {
			return nil, err
		}

		botVerifier, err := do.Invoke[service.TelegramTokenVerifier](i)
		if err != nil {
			return nil, err
		}

		return usecase.NewSyncBot(transactor, botRepo, botVerifier)
	})

	do.Provide(injector, func(i do.Injector) (*usecase.ResolveLoginChallenge, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		hydraClient, err := do.Invoke[*hydra.APIClient](i)
		if err != nil {
			return nil, err
		}

		botRepo, err := do.Invoke[repository.BotRepositoryPort](i)
		if err != nil {
			return nil, err
		}

		botUserRepo, err := do.Invoke[repository.BotUserRepositoryPort](i)
		if err != nil {
			return nil, err
		}

		tokenVerifier, err := do.Invoke[service.TelegramTokenVerifier](i)
		if err != nil {
			return nil, err
		}

		var baseUri *url.URL
		if cfg.HTTPServer.BaseUri != (config.URL{}) {
			baseUri = cfg.HTTPServer.BaseUri.URL()
		} else {
			uri, err := buildBaseURL(cfg.HTTPServer.Address)
			if err != nil {
				return nil, err
			}
			baseUri = uri
		}

		return usecase.NewResolveLoginChallenge(
			baseUri,
			cfg.HTTPServer.TelegramAuthURI.URL(),
			hydraClient,
			botRepo,
			botUserRepo,
			tokenVerifier,
		)
	})
}
