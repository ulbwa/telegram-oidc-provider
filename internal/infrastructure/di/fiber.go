package di

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
	domain_services "github.com/ulbwa/telegram-oidc-provider/internal/domain/services"
	"github.com/ulbwa/telegram-oidc-provider/internal/transport/http/middlewares"
	private_clients "github.com/ulbwa/telegram-oidc-provider/internal/transport/http/private/handlers/clients"
	private_login "github.com/ulbwa/telegram-oidc-provider/internal/transport/http/private/handlers/login"
)

type fiberApp struct {
	*fiber.App
}

func (app *fiberApp) Shutdown(ctx context.Context) error {
	return app.ShutdownWithContext(ctx)
}

func provideFiberApp(injector do.Injector) {
	// Validator
	do.Provide(injector, func(i do.Injector) (*validator.Validate, error) {
		return validator.New(), nil
	})

	// Controllers
	do.Provide(injector, func(i do.Injector) (*private_login.LoginController, error) {
		validator, err := do.Invoke[*validator.Validate](i)
		if err != nil {
			return nil, err
		}
		loginByTelegramWidget, err := do.Invoke[*usecases.LoginByTelegramWidget](i)
		if err != nil {
			return nil, err
		}
		loginByTelegramMiniApp, err := do.Invoke[*usecases.LoginByTelegramMiniApp](i)
		if err != nil {
			return nil, err
		}

		return private_login.NewLoginController(validator, loginByTelegramWidget, loginByTelegramMiniApp)
	})

	do.Provide(injector, func(i do.Injector) (*private_clients.ClientsController, error) {
		validator, err := do.Invoke[*validator.Validate](i)
		if err != nil {
			return nil, err
		}
		createClient, err := do.Invoke[*usecases.CreateClient](i)
		if err != nil {
			return nil, err
		}

		return private_clients.NewClientsController(validator, createClient)
	})

	// App
	do.Provide(injector, func(i do.Injector) (*fiberApp, error) {
		cfg, err := do.Invoke[*common.Config](i)
		if err != nil {
			return nil, err
		}

		disableStartup := !cfg.Logger.Console.Enabled || !cfg.Logger.Console.Pretty

		app := fiber.New(fiber.Config{
			ErrorHandler:          middlewares.ErrorHandler(),
			DisableStartupMessage: disableStartup,
		})

		// Logger Middlewares
		logger, err := do.Invoke[zerolog.Logger](i)
		if err != nil {
			return nil, err
		}
		app.Use(middlewares.LoggerMiddleware(&logger))
		app.Use(fiberzerolog.New(fiberzerolog.Config{
			GetLogger: func(c *fiber.Ctx) zerolog.Logger {
				return *zerolog.Ctx(c.UserContext())
			},
		}))

		// TraceID Middleware
		idGen, err := do.Invoke[domain_services.IdGenerator](i)
		if err != nil {
			return nil, err
		}
		app.Use(middlewares.TraceIDMiddleware(idGen))

		// Compression Middleware
		// TODO: configure from config
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed, // 1
		}))

		// Private controllers
		privateGroup := app.Group("/private")

		// Private Clients Controller
		clientController, err := do.Invoke[*private_clients.ClientsController](i)
		if err != nil {
			return nil, err
		}
		clientController.SetupRoutes(privateGroup)

		return &fiberApp{App: app}, nil
	})

	do.Provide(injector, func(i do.Injector) (*fiber.App, error) {
		if app, err := do.Invoke[*fiberApp](i); err != nil {
			return nil, err
		} else {
			return app.App, nil
		}
	})
}
