package di

import (
	"time"

	"github.com/samber/do/v2"
	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

func provideUsecases(injector do.Injector) {
	// CreateClient use case
	do.Provide(injector, func(i do.Injector) (*usecases.CreateClient, error) {
		hydraClient, err := do.Invoke[app_services.HydraClientManager](i)
		if err != nil {
			return nil, err
		}
		transactor, err := do.Invoke[app_services.Transactor](i)
		if err != nil {
			return nil, err
		}
		botRepo, err := do.Invoke[repositories.BotRepository](i)
		if err != nil {
			return nil, err
		}

		// TODO: Implement TelegramBotVerifier
		return usecases.NewCreateClient(
			nil, // botVerifier - TODO
			hydraClient,
			transactor,
			botRepo,
		), nil
	})

	// LoginByTelegramWidget use case
	do.Provide(injector, func(i do.Injector) (*usecases.LoginByTelegramWidget, error) {
		userRepo, err := do.Invoke[repositories.UserRepository](i)
		if err != nil {
			return nil, err
		}
		loginRepo, err := do.Invoke[repositories.UserBotLoginRepository](i)
		if err != nil {
			return nil, err
		}
		botRepo, err := do.Invoke[repositories.BotRepository](i)
		if err != nil {
			return nil, err
		}
		hydra, err := do.Invoke[app_services.HydraLoginManager](i)
		if err != nil {
			return nil, err
		}
		transactor, err := do.Invoke[app_services.Transactor](i)
		if err != nil {
			return nil, err
		}

		// TODO: Implement Telegram services
		return usecases.NewLoginByTelegramWidget(
			time.Hour*24,
			userRepo,
			loginRepo,
			botRepo,
			hydra,
			transactor,
			nil, // parser - TODO
			nil, // verifier - TODO
			nil, // replayGuard - TODO
			nil, // userFactory - TODO
		), nil
	})

	// LoginByTelegramMiniApp use case
	do.Provide(injector, func(i do.Injector) (*usecases.LoginByTelegramMiniApp, error) {
		userRepo, err := do.Invoke[repositories.UserRepository](i)
		if err != nil {
			return nil, err
		}
		loginRepo, err := do.Invoke[repositories.UserBotLoginRepository](i)
		if err != nil {
			return nil, err
		}
		botRepo, err := do.Invoke[repositories.BotRepository](i)
		if err != nil {
			return nil, err
		}
		hydra, err := do.Invoke[app_services.HydraLoginManager](i)
		if err != nil {
			return nil, err
		}
		transactor, err := do.Invoke[app_services.Transactor](i)
		if err != nil {
			return nil, err
		}

		// TODO: Implement Telegram services
		return usecases.NewLoginByTelegramMiniApp(
			time.Hour*24,
			userRepo,
			loginRepo,
			botRepo,
			hydra,
			transactor,
			nil, // parser - TODO
			nil, // verifier - TODO
			nil, // replayGuard - TODO
			nil, // userFactory - TODO
		), nil
	})
}
