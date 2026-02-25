package di

import (
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/db/postgres"
	"gorm.io/gorm"
)

func provideRepositories(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (service.Transactor, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return postgres.NewGormTransactor(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repository.BotUserRepositoryPort, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return postgres.NewBotUserRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repository.BotRepositoryPort, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		return postgres.NewBotRepository(db, []byte(cfg.Security.BotToken.EncryptionKey))
	})
}
