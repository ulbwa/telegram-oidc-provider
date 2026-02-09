package login

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
	"github.com/ulbwa/telegram-oidc-provider/internal/transport/http/errors"
)

func (c *LoginController) LoginByTelegramMiniApp(ctx *fiber.Ctx) error {
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

	var ucInput usecases.LoginByTelegramMiniAppInput
	ucInput.LoginChallenge = loginChallenge
	ucInput.InitData = queryWithoutChallenge.Encode()

	ucOutput, err := c.loginByTelegramMiniAppUC.Execute(ctx.UserContext(), &ucInput)
	if err != nil {
		// TODO: map errors
		return err
	}

	return ctx.Redirect(ucOutput.RedirectTo.String())
}
