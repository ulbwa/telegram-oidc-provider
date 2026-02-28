package di

import (
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
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
}
