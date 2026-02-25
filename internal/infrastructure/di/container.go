package di

import (
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
)

func NewContainer(cfg *config.Config) do.Injector {
	injector := do.New()

	// Config
	do.ProvideValue(injector, cfg)

	provideZerolog(injector)
	provideGorm(injector)
	provideRedis(injector)
	provideHydra(injector)
	provideServices(injector)
	provideRepositories(injector)
	provideUsecases(injector)
	provideEchoApp(injector)

	return injector
}
