package di

import (
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
)

func NewContainer(cfg *common.Config) do.Injector {
	injector := do.New()

	// Config
	do.ProvideValue(injector, cfg)

	provideZerolog(injector)
	provideServices(injector)
	provideRepositories(injector)
	provideUsecases(injector)
	provideFiberApp(injector)

	return injector
}
