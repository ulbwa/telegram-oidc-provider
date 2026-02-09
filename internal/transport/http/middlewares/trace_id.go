package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/services"
)

// TraceIDMiddleware extracts or generates a trace ID from request headers,
// adds it to the response headers, and propagates it through the context and logger.
func TraceIDMiddleware(idGen services.IdGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get or generate trace ID
		requestID := c.Get(fiber.HeaderXRequestID)
		if requestID == "" {
			requestID = idGen.Generate()
		}

		// Set trace ID in response header for client debugging
		c.Set(fiber.HeaderXRequestID, requestID)

		ctx := c.UserContext()

		// Add trace ID to logger
		zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("trace_id", requestID)
		})

		// Add trace ID to context
		ctx = common.WithTraceID(ctx, requestID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
