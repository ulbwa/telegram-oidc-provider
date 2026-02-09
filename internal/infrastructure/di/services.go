package di

import (
	"github.com/samber/do/v2"
	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	domain_services "github.com/ulbwa/telegram-oidc-provider/internal/domain/services"
	infra_services "github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/services"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

func provideServices(injector do.Injector) {
	// Domain services
	do.Provide(injector, func(i do.Injector) (domain_services.IdGenerator, error) {
		return infra_services.NewUUIDv7Generator()
	})

	// Application services
	do.Provide(injector, func(i do.Injector) (app_services.Transactor, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}
		return postgres.NewTransactor(db), nil
	})
}
