package web

import (
	"html/template"
	"io"
	"net/url"

	"github.com/labstack/echo/v4"
	hydra "github.com/ory/hydra-client-go"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/internal/interface/http/web/templates"
)

type server struct {
	baseUri         *url.URL
	errorUri        *url.URL
	telegramAuthUri *url.URL // url.Parse("https://oauth.telegram.org/auth")

	hydra       *hydra.APIClient
	transactor  service.Transactor
	botRepo     repository.BotRepositoryPort
	botVerifier service.TelegramTokenVerifier
}

type renderer struct {
	tmpl *template.Template
}

func (r *renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return r.tmpl.ExecuteTemplate(w, name, data)
}

func NewServer(
	baseUri *url.URL,
	errorUri *url.URL,
	telegramAuthUri *url.URL,
	hydraClient *hydra.APIClient,
	transactor service.Transactor,
	botRepo repository.BotRepositoryPort,
	botVerifier service.TelegramTokenVerifier,
) *server {
	return &server{
		baseUri:         baseUri,
		errorUri:        errorUri,
		telegramAuthUri: telegramAuthUri,
		hydra:           hydraClient,
		transactor:      transactor,
		botRepo:         botRepo,
		botVerifier:     botVerifier,
	}
}

func NewRenderer() (echo.Renderer, error) {
	tmpl, err := template.New("web").
		New("login").Parse(templates.LoginTemplate())
	if err != nil {
		return nil, err
	}

	if _, err = tmpl.New("error").Parse(templates.ErrorTemplate()); err != nil {
		return nil, err
	}

	if _, err = tmpl.New("consent").Parse(templates.ConsentTemplate()); err != nil {
		return nil, err
	}

	return &renderer{tmpl: tmpl}, nil
}

func (s *server) Register(e *echo.Echo) {
	e.GET("/login", s.Login)
	e.GET("/error", s.Error)
}
