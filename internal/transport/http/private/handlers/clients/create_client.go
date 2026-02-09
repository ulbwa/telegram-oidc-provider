package clients

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
	"github.com/ulbwa/telegram-oidc-provider/internal/transport/http/errors"
)

type CreateClientBody struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	RedirectUri string `json:"redirectUri" validate:"required,url"`
	BotToken    string `json:"botToken" validate:"required"`
}

func (c *ClientsController) CreateClient(ctx *fiber.Ctx) error {
	var body CreateClientBody
	if err := ctx.BodyParser(&body); err != nil {
		return errors.ErrInvalidPayload
	}
	if err := c.validator.Struct(body); err != nil {
		return errors.ErrInvalidPayloadValidationFailed
	}

	var ucInput usecases.CreateClientInput
	ucInput.Name = body.Name
	ucInput.RedirectUri = func() *url.URL {
		if redirectUri, err := url.Parse(body.RedirectUri); err != nil {
			panic(fmt.Errorf("invalid redirect URI: %w", err))
		} else {
			return redirectUri
		}
	}()
	ucInput.BotToken = body.BotToken

	ucOutput, err := c.createClientUC.Execute(ctx.UserContext(), &ucInput)
	if err != nil {
		// TODO: map errors from use case to generic HTTP errors
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(ucOutput.Id)
}
