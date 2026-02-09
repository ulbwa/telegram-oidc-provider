package di

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/storage/postgres"
)

func provideRepositories(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (repositories.BotRepository, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}
		return postgres.NewBotRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repositories.UserRepository, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}
		return postgres.NewUserRepository(db), nil
	})

	do.Provide(injector, func(i do.Injector) (repositories.UserBotLoginRepository, error) {
		db, err := do.Invoke[*gorm.DB](i)
		if err != nil {
			return nil, err
		}
		return postgres.NewUserBotLoginRepository(db), nil
	})
}
