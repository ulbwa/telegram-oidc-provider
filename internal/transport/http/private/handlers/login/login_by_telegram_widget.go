package login

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
	"github.com/ulbwa/telegram-oidc-provider/internal/transport/http/errors"
)

type LoginByTelegramWidgetBody struct {
	LoginChallenge string `json:"loginChallenge" validate:"required"`
	WidgetData     string `json:"widgetData"     validate:"required"`
}

func (c *LoginController) LoginByTelegramWidget(ctx *fiber.Ctx) error {
	query, err := url.ParseQuery(string(ctx.Request().URI().QueryString()))
	if err != nil {
		return errors.ErrInvalidQuery
	}
	loginChallenge := query.Get("login_challenge")
	if loginChallenge == "" {
		return errors.ErrInvalidQuery
	}
	queryWithoutChallenge := query
	queryWithoutChallenge.Del("login_challenge")

	var ucInput usecases.LoginByTelegramWidgetInput
	ucInput.LoginChallenge = loginChallenge
	ucInput.WidgetData = queryWithoutChallenge.Encode()

	ucOutput, err := c.loginByTelegramWidgetUC.Execute(ctx.UserContext(), &ucInput)
	if err != nil {
		// TODO: map errors
		return err
	}

	return ctx.Redirect(ucOutput.RedirectTo.String())
}
