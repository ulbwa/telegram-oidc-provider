package errors

import (
	"net/http"
)

var (
	CodeBadRequest   = "BAD_REQUEST"
	CodeUnauthorized = "UNAUTHORIZED"
)

var (
	ErrBadRequest = &GenericError{
		Code:     CodeBadRequest,
		HttpCode: http.StatusBadRequest,
		Message:  "The request is malformed or contains invalid parameters.",
	}

	ErrInvalidQuery = ErrBadRequest.
			WithReason("INVALID_QUERY").
			WithMsg("The query parameters are invalid.")

	ErrInvalidPayload = ErrBadRequest.
				WithReason("INVALID_PAYLOAD").
				WithMsg("The request payload is invalid.")

	ErrInvalidPayloadValidationFailed = ErrInvalidPayload.
						WithReason("INVALID_PAYLOAD").
						WithMsg("The request payload validation has failed.")

	ErrUnableToParseJSON = ErrInvalidPayload.
				WithMsg("Unable to parse the request JSON payload.")

	ErrValidationFailed = ErrInvalidPayload.
				WithMsg("The request payload validation has failed.")

	ErrInternal = &GenericError{
		Code:     "INTERNAL",
		HttpCode: http.StatusInternalServerError,
		Message:  "An internal server error occurred.",
	}
)
