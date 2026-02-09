package services

import (
	"context"
	"net/url"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type OAuth2ErrorCode string

const (
	OAuth2ErrorInvalidRequest          OAuth2ErrorCode = "invalid_request"
	OAuth2ErrorUnauthorizedClient      OAuth2ErrorCode = "unauthorized_client"
	OAuth2ErrorAccessDenied            OAuth2ErrorCode = "access_denied"
	OAuth2ErrorUnsupportedResponseType OAuth2ErrorCode = "unsupported_response_type"
	OAuth2ErrorInvalidScope            OAuth2ErrorCode = "invalid_scope"
	OAuth2ErrorServerError             OAuth2ErrorCode = "server_error"
	OAuth2ErrorUnavailable             OAuth2ErrorCode = "temporarily_unavailable"

	OIDCErrorInsideRequired           OAuth2ErrorCode = "interaction_required"
	OIDCErrorLoginRequired            OAuth2ErrorCode = "login_required"
	OIDCErrorConsentRequired          OAuth2ErrorCode = "consent_required"
	OIDCErrorAccountSelectionRequired OAuth2ErrorCode = "account_selection_required"
	OIDCErrorInvalidRequestURI        OAuth2ErrorCode = "invalid_request_uri"
	OIDCErrorInvalidRequestObject     OAuth2ErrorCode = "invalid_request_object"
	OIDCErrorRequestNotSupported      OAuth2ErrorCode = "request_not_supported"
	OIDCErrorRequestURINotSupported   OAuth2ErrorCode = "request_uri_not_supported"
	OIDCErrorRegistrationNotSupported OAuth2ErrorCode = "registration_not_supported"
)

type LoginRequestResponse struct {
	RequestUrl *url.URL
	ClientId   string
	Subject    *string
	SessionId  string
	Skip       bool
}

func (l *LoginRequestResponse) IsAuthenticated() bool {
	return l.Subject != nil
}

func (l *LoginRequestResponse) SkipRequired() bool {
	return l.Skip && l.IsAuthenticated()
}

type AcceptedLoginRequestResponse struct {
	RedirectTo *url.URL
}

type RejectLoginRequestResponse struct {
	RedirectTo *url.URL
}

type CreateClientResponse struct {
	Id          string
	Secret      string
	Name        string
	RedirectUri *url.URL
}

type HydraLoginManager interface {
	GetLoginRequest(
		ctx context.Context,
		challenge string,
	) (*LoginRequestResponse, error)

	AcceptLoginRequest(
		ctx context.Context,
		challenge string,
		user *domain.User,
		remember bool,
	) (*AcceptedLoginRequestResponse, error)

	RejectLoginRequest(
		ctx context.Context,
		challenge string,
		error OAuth2ErrorCode,
	) (*RejectLoginRequestResponse, error)

	RevokeLoginSession(
		ctx context.Context,
		sessionId string,
	) error
}

type HydraClientManager interface {
	CreateClient(
		ctx context.Context,
		name string,
		redirectUri *url.URL,
	) (*CreateClientResponse, error)

	DeleteClient(
		ctx context.Context,
		id string,
	) error
}
