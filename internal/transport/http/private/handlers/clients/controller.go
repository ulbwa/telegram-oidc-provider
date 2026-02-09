package clients

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecases"
)

type ClientsController struct {
	validator *validator.Validate

	createClientUC *usecases.CreateClient
}

func NewClientsController(
	validator *validator.Validate,
	createClient *usecases.CreateClient,
) (*ClientsController, error) {
	return &ClientsController{
		validator: validator,

		createClientUC: createClient,
	}, nil
}

func (c *ClientsController) SetupRoutes(router fiber.Router) {
	group := router.Group("/clients")
	group.Post("", c.CreateClient)
}
