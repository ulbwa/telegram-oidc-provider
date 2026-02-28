//go:generate go-enum --values --names --nocase
package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	hydra "github.com/ory/hydra-client-go"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/pkg/utils"
)

type ResolveLoginChallenge struct {
	baseUri         *url.URL
	telegramAuthUri *url.URL

	transactor    service.Transactor
	hydra         *hydra.APIClient
	botRepo       repository.BotRepositoryPort
	botUserRepo   repository.BotUserRepositoryPort
	tokenVerifier service.TelegramTokenVerifier
}

func NewResolveLoginChallenge(
	baseUri *url.URL,
	telegramAuthUri *url.URL,
	transactor service.Transactor,
	hydraClient *hydra.APIClient,
	botRepo repository.BotRepositoryPort,
	botUserRepo repository.BotUserRepositoryPort,
	tokenVerifier service.TelegramTokenVerifier,
) (*ResolveLoginChallenge, error) {
	if baseUri == nil {
		return nil, errors.New("base URI is nil")
	}
	if telegramAuthUri == nil {
		return nil, errors.New("telegram auth URI is nil")
	}
	if transactor == nil {
		return nil, errors.New("transactor is nil")
	}
	if hydraClient == nil {
		return nil, errors.New("hydra client is nil")
	}
	if botRepo == nil {
		return nil, errors.New("bot repository is nil")
	}
	if botUserRepo == nil {
		return nil, errors.New("bot user repository is nil")
	}
	if tokenVerifier == nil {
		return nil, errors.New("token verifier is nil")
	}

	return &ResolveLoginChallenge{
		baseUri:         baseUri,
		telegramAuthUri: telegramAuthUri,
		transactor:      transactor,
		hydra:           hydraClient,
		botRepo:         botRepo,
		botUserRepo:     botUserRepo,
		tokenVerifier:   tokenVerifier,
	}, nil
}

type (
	// ResolveLoginChallengeAction enum for login challenge handling action
	// ENUM(
	//     Redirect
	//     Render
	// )
	ResolveLoginChallengeAction string

	ResolveLoginChallengeInput struct {
		LoginChallenge string
	}
	ResolveLoginChallengeOutput struct {
		Action             ResolveLoginChallengeAction
		RedirectUri        *string
		WidgetUri          *string
		MiniAppCallbackUri *string
	}
)

func (uc *ResolveLoginChallenge) verifyChallenge(challenge string) error {
	if challenge == "" {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("login", "challenge", nil))
	}
	return nil
}

func (uc *ResolveLoginChallenge) getBot(ctx context.Context, clientId string) (*entity.Bot, error) {
	var bot entity.Bot
	if err := uc.botRepo.GetByClientID(ctx, clientId, &bot); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectNotFoundErr("client", clientId))
		}
		return nil, ErrUnexpected
	}
	return &bot, nil
}

func (uc *ResolveLoginChallenge) verifyBotToken(ctx context.Context, botToken string) error {
	if _, err := uc.tokenVerifier.Verify(ctx, botToken, service.NewVerifyOptions(service.WithSkipCacheRead())); err != nil {
		if errors.Is(err, service.ErrTelegramBotTokenMalformed) {
			return fmt.Errorf(
				"%w: %w",
				ErrInvalidInput,
				NewObjectInvalidErr("bot", "token", utils.Ptr("malformed")))
		}
		if errors.Is(err, service.ErrTelegramBotTokenInvalid) {
			return fmt.Errorf(
				"%w: %w",
				ErrInvalidInput,
				NewObjectInvalidErr("bot", "token", nil))
		}
		return ErrUnexpected
	}
	return nil
}

func (uc *ResolveLoginChallenge) getLoginRequest(ctx context.Context, loginChallenge string) (*hydra.LoginRequest, error) {
	loginRequest, resp, err := uc.hydra.AdminApi.
		GetLoginRequest(ctx).
		LoginChallenge(loginChallenge).
		Execute()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, NewGatewayTimeoutErr("hydra")
		}
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
				return nil, fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("login", "challenge", nil))
			}
			if resp.StatusCode >= http.StatusInternalServerError {
				return nil, NewBadGatewayErr("hydra")
			}
			return nil, ErrUnexpected
		}
		return nil, NewBadGatewayErr("hydra")
	}

	return loginRequest, nil
}

func (uc *ResolveLoginChallenge) buildRenderOutput(loginChallenge string, bot *entity.Bot) *ResolveLoginChallengeOutput {
	origin := *uc.baseUri
	origin = *origin.JoinPath("/login")

	widgetCallbackUri := *uc.baseUri
	widgetCallbackUri = *origin.JoinPath("/widget/callback")
	widgetCallbackUriQuery := widgetCallbackUri.Query()
	widgetCallbackUriQuery.Set("login_challenge", loginChallenge)
	widgetCallbackUri.RawQuery = widgetCallbackUriQuery.Encode()

	widgetUri := *uc.telegramAuthUri
	widgetUriQuery := widgetUri.Query()
	widgetUriQuery.Set("bot_id", fmt.Sprintf("%d", bot.Id))
	widgetUriQuery.Set("origin", origin.String())
	widgetUriQuery.Set("request_access", "write")
	widgetUriQuery.Set("return_to", widgetCallbackUri.String())
	widgetUri.RawQuery = widgetUriQuery.Encode()

	miniappCallbackUri := *uc.baseUri
	miniappCallbackUri = *origin.JoinPath("/miniapp/callback")
	miniappCallbackUriQuery := miniappCallbackUri.Query()
	miniappCallbackUriQuery.Set("login_challenge", loginChallenge)
	miniappCallbackUri.RawQuery = miniappCallbackUriQuery.Encode()

	return &ResolveLoginChallengeOutput{
		Action:             ResolveLoginChallengeActionRender,
		WidgetUri:          utils.Ptr(widgetUri.String()),
		MiniAppCallbackUri: utils.Ptr(miniappCallbackUri.String()),
	}
}

func (uc *ResolveLoginChallenge) buildRedirectOutput(redirectUri string) *ResolveLoginChallengeOutput {
	return &ResolveLoginChallengeOutput{
		Action:      ResolveLoginChallengeActionRedirect,
		RedirectUri: utils.Ptr(redirectUri),
	}
}

func (uc *ResolveLoginChallenge) updateLastLogin(ctx context.Context, botId, userId int64) error {
	var user entity.BotUser
	if err := uc.botUserRepo.GetByBotAndUser(ctx, botId, userId, &user); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewObjectNotFoundErr("user", userId)
		}
		return ErrUnexpected
	}

	user.UpdateLastLogin()
	if err := uc.botUserRepo.Update(ctx, &user); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewObjectNotFoundErr("user", userId)
		}
		return ErrUnexpected
	}

	return nil
}

func (uc *ResolveLoginChallenge) parseSubjectUserId(subject string) (int64, error) {
	if subject == "" {
		return 0, NewObjectInvalidErr("login", "subject", utils.Ptr("empty"))
	}
	userId, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return 0, NewObjectInvalidErr("login", "subject", nil)
	}
	return userId, nil
}

func (uc *ResolveLoginChallenge) acceptLoginRequest(ctx context.Context, loginChallenge string, userId int64) (*hydra.CompletedRequest, error) {
	acceptReq := hydra.NewAcceptLoginRequest(strconv.FormatInt(userId, 10))

	completed, resp, err := uc.hydra.AdminApi.
		AcceptLoginRequest(ctx).
		LoginChallenge(loginChallenge).
		AcceptLoginRequest(*acceptReq).
		Execute()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, NewGatewayTimeoutErr("hydra")
		}
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {
				return nil, fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("login", "challenge", nil))
			}
			if resp.StatusCode >= http.StatusInternalServerError {
				return nil, NewBadGatewayErr("hydra")
			}
			return nil, ErrUnexpected
		}
		return nil, NewBadGatewayErr("hydra")
	}

	return completed, nil
}

func (uc *ResolveLoginChallenge) mapRejectError(err error) (string, int64, string) {
	if err == nil {
		return "server_error", http.StatusInternalServerError, "unexpected authentication error"
	}

	var gatewayTimeoutErr *GatewayTimeoutErr
	if errors.As(err, &gatewayTimeoutErr) {
		return "temporarily_unavailable", http.StatusServiceUnavailable, "authentication service is temporarily unavailable"
	}

	var badGatewayErr *BadGatewayErr
	if errors.As(err, &badGatewayErr) {
		return "temporarily_unavailable", http.StatusServiceUnavailable, "authentication service is temporarily unavailable"
	}

	var objectInvalidErr *ObjectInvalidErr
	if errors.As(err, &objectInvalidErr) {
		if objectInvalidErr.Object == "bot" && objectInvalidErr.Field == "token" {
			return "unauthorized_client", http.StatusBadRequest, "client is linked to invalid bot credentials"
		}
		if objectInvalidErr.Object == "login" && objectInvalidErr.Field == "challenge" {
			return "invalid_request", http.StatusBadRequest, "invalid login challenge"
		}
		return "invalid_request", http.StatusBadRequest, "invalid authentication request"
	}

	var objectNotFoundErr *ObjectNotFoundErr
	if errors.As(err, &objectNotFoundErr) {
		if objectNotFoundErr.Object == "client" {
			return "unauthorized_client", http.StatusBadRequest, "oauth2 client is not linked to bot configuration"
		}
		return "access_denied", http.StatusForbidden, "authentication cannot be completed"
	}

	if errors.Is(err, ErrInvalidInput) {
		return "invalid_request", http.StatusBadRequest, "invalid authentication request"
	}

	if errors.Is(err, ErrUnexpected) {
		return "server_error", http.StatusInternalServerError, "internal authentication error"
	}

	return "server_error", http.StatusInternalServerError, "internal authentication error"
}

func (uc *ResolveLoginChallenge) rejectLoginRequest(ctx context.Context, loginChallenge string, reason error) (*ResolveLoginChallengeOutput, error) {
	reasonDebug := "unknown"
	if reason != nil {
		reasonDebug = reason.Error()
	}

	oauth2Error, statusCode, description := uc.mapRejectError(reason)
	rejectReq := hydra.NewRejectRequest()
	rejectReq.SetError(oauth2Error)
	rejectReq.SetStatusCode(int64(statusCode))
	rejectReq.SetErrorDescription(description)
	rejectReq.SetErrorHint("authentication request was rejected")
	rejectReq.SetErrorDebug(reasonDebug)

	completed, resp, err := uc.hydra.AdminApi.
		RejectLoginRequest(ctx).
		LoginChallenge(loginChallenge).
		RejectRequest(*rejectReq).
		Execute()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, NewGatewayTimeoutErr("hydra")
		}
		if resp != nil {
			if resp.StatusCode >= http.StatusInternalServerError {
				return nil, NewBadGatewayErr("hydra")
			}
			return nil, ErrUnexpected
		}
		return nil, NewBadGatewayErr("hydra")
	}

	if completed == nil || completed.RedirectTo == "" {
		return nil, ErrUnexpected
	}

	return uc.buildRedirectOutput(completed.RedirectTo), nil
}

func (uc *ResolveLoginChallenge) rejectAfterChallenge(ctx context.Context, loginChallenge string, reason error) (*ResolveLoginChallengeOutput, error) {
	zerolog.Ctx(ctx).Warn().
		Err(reason).
		Str("login_challenge", loginChallenge).
		Msg("resolve login challenge failed, rejecting login request in hydra")

	output, err := uc.rejectLoginRequest(ctx, loginChallenge, reason)
	if err != nil {
		zerolog.Ctx(ctx).Error().
			Err(err).
			Str("login_challenge", loginChallenge).
			Msg("failed to reject login request in hydra")
		return nil, err
	}

	return output, nil
}

func (uc *ResolveLoginChallenge) Execute(ctx context.Context, input *ResolveLoginChallengeInput) (*ResolveLoginChallengeOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	challenge := input.LoginChallenge
	if err := uc.verifyChallenge(challenge); err != nil {
		return nil, err
	}

	loginRequest, err := uc.getLoginRequest(ctx, challenge)
	if err != nil {
		return uc.rejectAfterChallenge(ctx, challenge, err)
	}
	if loginRequest == nil || loginRequest.Client.ClientId == nil {
		return uc.rejectAfterChallenge(ctx, challenge, ErrUnexpected)
	}
	clientId := *loginRequest.Client.ClientId

	bot, err := uc.getBot(ctx, clientId)
	if err != nil {
		return uc.rejectAfterChallenge(ctx, challenge, err)
	}

	if err := uc.verifyBotToken(ctx, bot.Token); err != nil {
		return uc.rejectAfterChallenge(ctx, challenge, err)
	}

	if loginRequest.Skip {
		skipUserId, err := uc.parseSubjectUserId(loginRequest.Subject)
		if err == nil {
			err = uc.transactor.RunInTransaction(ctx, func(txCtx context.Context) error {
				return uc.updateLastLogin(txCtx, bot.Id, skipUserId)
			})
		}

		if err == nil {
			completed, acceptErr := uc.acceptLoginRequest(ctx, loginRequest.Challenge, skipUserId)
			if acceptErr == nil && completed != nil && completed.RedirectTo != "" {
				return uc.buildRedirectOutput(completed.RedirectTo), nil
			}
			if acceptErr != nil {
				err = acceptErr
			} else {
				err = ErrUnexpected
			}
		}

		zerolog.Ctx(ctx).Warn().
			Err(err).
			Str("login_challenge", challenge).
			Str("client_id", clientId).
			Str("subject", loginRequest.Subject).
			Msg("skip login failed, falling back to interactive login UI")
	}

	return uc.buildRenderOutput(challenge, bot), nil
}
