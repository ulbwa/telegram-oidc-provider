package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// LoggerMiddleware injects the logger into the request context for use in request handlers.
func LoggerMiddleware(logger *zerolog.Logger) fiber.Handler {
	if logger == nil {
		panic("logger is nil")
	}

	return func(c *fiber.Ctx) error {
		c.SetUserContext(logger.WithContext(c.UserContext()))
		return c.Next()
	}
}
