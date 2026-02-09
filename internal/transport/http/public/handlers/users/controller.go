package users

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UsersController struct {
	validator *validator.Validate
}

func (c *UsersController) SetupRoutes(app *fiber.App) {
	group := app.Group("/users")
	group.Get("", c.ListUsers)
	group.Get("/:id", c.GetUser)
}
