package di

import (
	hydra "github.com/ory/hydra-client-go"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
)

func provideHydra(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (*hydra.APIClient, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		hydraCfg := hydra.NewConfiguration()
		hydraCfg.Servers = hydra.ServerConfigurations{{URL: cfg.Hydra.AdminURL.URL().String()}}

		return hydra.NewAPIClient(hydraCfg), nil
	})
}
