package di

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/api/generated"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
	apihttp "github.com/ulbwa/telegram-oidc-provider/internal/interface/http/api"
)

func provideEchoApp(injector do.Injector) {
	do.Provide(injector, func(i do.Injector) (*echo.Echo, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		syncBot, err := do.Invoke[*usecase.SyncBot](i)
		if err != nil {
			return nil, err
		}

		var baseUri *url.URL
		if cfg.HTTPServer.BaseUri != (config.URL{}) {
			baseUri = cfg.HTTPServer.BaseUri.URL()
		} else {
			uri, err := buildBaseURL(cfg.HTTPServer.Address)
			if err != nil {
				return nil, err
			}
			baseUri = uri
		}

		apiServer, err := apihttp.NewServer(baseUri, syncBot)
		if err != nil {
			return nil, err
		}

		echoApp := echo.New()
		echoApp.HideBanner = shouldHideEchoBanner(cfg)
		echoApp.HidePort = shouldHideEchoBanner(cfg)

		// Add request/response validation middleware for API endpoints
		spec, err := generated.GetSwagger()
		if err != nil {
			return nil, fmt.Errorf("failed to get OpenAPI spec: %w", err)
		}
		echoApp.Use(echo_middleware.OapiRequestValidator(spec))

		generated.RegisterHandlers(echoApp, generated.NewStrictHandler(apiServer, nil))

		return echoApp, nil
	})
}

func shouldHideEchoBanner(cfg *config.Config) bool {
	return cfg.Logger.Console.Enabled && !cfg.Logger.Console.Pretty
}

func buildBaseURL(address string) (*url.URL, error) {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return url.Parse(address)
	}
	if strings.HasPrefix(address, ":") {
		return url.Parse(fmt.Sprintf("http://localhost%s", address))
	}
	return url.Parse(fmt.Sprintf("http://%s", address))
}
