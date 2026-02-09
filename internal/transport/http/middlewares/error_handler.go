package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
	generic_errors "github.com/ulbwa/telegram-oidc-provider/internal/transport/http/errors"
)

type errResponse struct {
	Code    string `json:"code"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message"`
	TraceID string `json:"traceId"`
}

// ErrorHandler returns a Fiber error handler middleware that processes different error types
// and returns appropriate HTTP responses with error details and trace ID.
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var genErr *generic_errors.GenericError
		if ok := errors.As(err, &genErr); ok {
			return processGenError(c, genErr)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return processFiberError(c, fiberErr)
		}

		return processOtherError(c, err)
	}
}

func processGenError(c *fiber.Ctx, err *generic_errors.GenericError) error {
	zerolog.Ctx(c.UserContext()).Error().Err(err).Msg("generic error")

	return c.
		Status(err.HttpCode).
		JSON(errResponse{
			Code:    err.Code,
			Reason:  err.Reason,
			Message: err.Message,
			TraceID: common.GetTraceID(c.UserContext()),
		})
}

func processFiberError(c *fiber.Ctx, err *fiber.Error) error {
	zerolog.Ctx(c.UserContext()).Error().Err(err).Msg("fiber error")

	return c.Status(err.Code).JSON(&errResponse{
		Code:    "HTTP_ERROR",
		Message: err.Message,
		TraceID: common.GetTraceID(c.UserContext()),
	})
}

func processOtherError(c *fiber.Ctx, err error) error {
	zerolog.Ctx(c.UserContext()).Error().Err(err).Msg("internal error")
	return c.Status(generic_errors.ErrInternal.HttpCode).JSON(&errResponse{
		Code:    generic_errors.ErrInternal.Code,
		Message: "An unexpected error has occurred while processing your request.",
		TraceID: common.GetTraceID(c.UserContext()),
	})
}
