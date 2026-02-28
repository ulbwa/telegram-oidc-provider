package web

import (
	"html/template"
	"io"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	"github.com/ulbwa/telegram-oidc-provider/internal/interface/http/web/templates"
)

type server struct {
	errorUri *url.URL

	resolveLoginChallengeUsecase *usecase.ResolveLoginChallenge
}

type renderer struct {
	tmpl *template.Template
}

func (r *renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return r.tmpl.ExecuteTemplate(w, name, data)
}

func NewServer(
	errorUri *url.URL,
	resolveLoginChallengeUsecase *usecase.ResolveLoginChallenge,
) *server {
	return &server{
		errorUri:                     errorUri,
		resolveLoginChallengeUsecase: resolveLoginChallengeUsecase,
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
