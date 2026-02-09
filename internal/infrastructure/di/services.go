package di

import (
	hydra_client "github.com/ory/hydra-client-go"
	"github.com/samber/do/v2"
	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
	domain_services "github.com/ulbwa/telegram-oidc-provider/internal/domain/services"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/hydra"
	infra_services "github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/services"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

func provideServices(injector do.Injector) {
	// Domain services
	do.Provide(injector, func(i do.Injector) (domain_services.IdGenerator, error) {
		return infra_services.NewUUIDv7Generator()
	})

	// Application services - Database
	do.Provide(injector, func(i do.Injector) (app_services.Transactor, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}
		return postgres.NewTransactor(db), nil
	})

	// Application services - Hydra
	do.Provide(injector, func(i do.Injector) (*hydra_client.APIClient, error) {
		cfg, err := do.Invoke[*common.Config](i)
		if err != nil {
			return nil, err
		}

		configuration := hydra_client.NewConfiguration()
		configuration.Servers = []hydra_client.ServerConfiguration{
			{
				URL: cfg.Hydra.AdminURL,
			},
		}
		return hydra_client.NewAPIClient(configuration), nil
	})

	do.Provide(injector, func(i do.Injector) (app_services.HydraLoginManager, error) {
		client, err := do.Invoke[*hydra_client.APIClient](i)
		if err != nil {
			return nil, err
		}
		return hydra.NewHydraClient(client)
	})

	do.Provide(injector, func(i do.Injector) (app_services.HydraClientManager, error) {
		client, err := do.Invoke[*hydra_client.APIClient](i)
		if err != nil {
			return nil, err
		}
		return hydra.NewHydraClient(client)
	})
}
