package login

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
)

type LoginController struct {
	validator *validator.Validate

	loginByTelegramWidgetUC  *usecases.LoginByTelegramWidget
	loginByTelegramMiniAppUC *usecases.LoginByTelegramMiniApp
}

func NewLoginController(
	validator *validator.Validate,
	loginByTelegramWidget *usecases.LoginByTelegramWidget,
	loginByTelegramMiniApp *usecases.LoginByTelegramMiniApp,
) (*LoginController, error) {
	return &LoginController{
		validator: validator,

		loginByTelegramWidgetUC:  loginByTelegramWidget,
		loginByTelegramMiniAppUC: loginByTelegramMiniApp,
	}, nil
}

func (c *LoginController) SetupRoutes(router fiber.Router) {
	group := router.Group("/login")
	group.Get("/widget", c.LoginByTelegramWidget)
	group.Get("/miniapp", c.LoginByTelegramMiniApp)
}
