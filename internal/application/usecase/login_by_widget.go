package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strconv"
	"time"

	hydra "github.com/ory/hydra-client-go"
	"github.com/rs/zerolog"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/pkg/utils"
)

type LoginByWidget struct {
	transactor        service.Transactor
	hydra             *hydra.APIClient
	widgetDataParser  service.TelegramWidgetDataParser
	authHashVerifier  service.TelegramAuthHashVerifier
	tokenVerifier     service.TelegramTokenVerifier
	botRepo           repository.BotRepositoryPort
	botUserRepo       repository.BotUserRepositoryPort
	authDataFreshness time.Duration
}

func NewLoginByWidget(
	transactor service.Transactor,
	hydraClient *hydra.APIClient,
	widgetDataParser service.TelegramWidgetDataParser,
	authHashVerifier service.TelegramAuthHashVerifier,
	tokenVerifier service.TelegramTokenVerifier,
	botRepo repository.BotRepositoryPort,
	botUserRepo repository.BotUserRepositoryPort,
	authDataFreshness time.Duration,
) (*LoginByWidget, error) {
	if transactor == nil {
		return nil, errors.New("transactor is nil")
	}
	if hydraClient == nil {
		return nil, errors.New("hydra client is nil")
	}
	if widgetDataParser == nil {
		return nil, errors.New("widget data parser is nil")
	}
	if authHashVerifier == nil {
		return nil, errors.New("auth hash verifier is nil")
	}
	if tokenVerifier == nil {
		return nil, errors.New("token verifier is nil")
	}
	if botRepo == nil {
		return nil, errors.New("bot repository is nil")
	}
	if botUserRepo == nil {
		return nil, errors.New("bot user repository is nil")
	}
	if authDataFreshness <= 0 {
		return nil, errors.New("auth data freshness must be positive")
	}

	return &LoginByWidget{
		transactor:        transactor,
		hydra:             hydraClient,
		widgetDataParser:  widgetDataParser,
		authHashVerifier:  authHashVerifier,
		tokenVerifier:     tokenVerifier,
		botRepo:           botRepo,
		botUserRepo:       botUserRepo,
		authDataFreshness: authDataFreshness,
	}, nil
}

type (
	LoginByWidgetInput struct {
		LoginChallenge string
		AuthData       map[string]any
		UserAgent      *string
		Language       *string
		ClientIP       netip.Addr
	}
	LoginByWidgetOutput struct {
		RedirectUri string
	}
)

func (uc *LoginByWidget) verifyChallenge(loginChallenge string) error {
	if loginChallenge == "" {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("login", "challenge", nil))
	}

	return nil
}

func (uc *LoginByWidget) verifyIP(clientIP netip.Addr) error {
	if !clientIP.IsValid() || clientIP.IsUnspecified() {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("request", "client_ip", nil))
	}

	return nil
}

func (uc *LoginByWidget) getLoginRequest(ctx context.Context, loginChallenge string) (*hydra.LoginRequest, error) {
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

func (uc *LoginByWidget) getBot(ctx context.Context, clientId string) (*entity.Bot, error) {
	var bot entity.Bot
	if err := uc.botRepo.GetByClientID(ctx, clientId, &bot); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectNotFoundErr("client", clientId))
		}
		return nil, ErrUnexpected
	}
	return &bot, nil
}

func (uc *LoginByWidget) verifyBotToken(ctx context.Context, botToken string) error {
	if _, err := uc.tokenVerifier.Verify(ctx, botToken, nil); err != nil {
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

func (uc *LoginByWidget) verifyAuthData(authData *service.TelegramAuthData, botToken string) error {
	if authData == nil || authData.User == nil {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("telegram_auth_data", "user", nil))
	}
	if authData.IsExpired(uc.authDataFreshness) {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("telegram_auth_data", "auth_date", utils.Ptr("expired")))
	}
	if err := uc.authHashVerifier.Verify(authData.Raw, authData.Hash, botToken); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("telegram_auth_data", "hash", nil))
	}

	return nil
}

func (uc *LoginByWidget) parseAndVerifyAuthData(authDataParams map[string]any, botToken string) (*service.TelegramAuthData, error) {
	if len(authDataParams) == 0 {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("telegram_auth_data", "payload", nil))
	}

	authData, err := uc.widgetDataParser.Parse(authDataParams)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("telegram_auth_data", "payload", nil))
	}

	if err := uc.verifyAuthData(authData, botToken); err != nil {
		return nil, err
	}

	return authData, nil
}

func (uc *LoginByWidget) ensureBotUserExists(
	ctx context.Context,
	botId int64,
	tgUser *service.TelegramUserData,
	clientIP netip.Addr,
	userAgent *string,
	language *string,
) error {
	var botUser entity.BotUser
	if err := uc.botUserRepo.GetByBotAndUser(ctx, botId, tgUser.Id, &botUser); err == nil {
		return nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		zerolog.Ctx(ctx).Error().
			Err(err).
			Int64("bot_id", botId).
			Int64("user_id", tgUser.Id).
			Msg("failed to load bot user by bot and user ids")
		return ErrUnexpected
	}

	user, err := entity.NewUser(tgUser.FirstName, tgUser.LastName, tgUser.Username, tgUser.PhotoUrl, tgUser.IsPremium)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("user", "profile", nil))
	}

	newBotUser, err := entity.NewBotUser(botId, tgUser.Id, user, clientIP, userAgent, language)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectInvalidErr("bot_user", "data", nil))
	}

	if err := uc.botUserRepo.Create(ctx, newBotUser); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil
		}
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%w: %w", ErrInvalidInput, NewObjectNotFoundErr("bot", botId))
		}
		zerolog.Ctx(ctx).Error().
			Err(err).
			Int64("bot_id", botId).
			Int64("user_id", tgUser.Id).
			Msg("failed to create bot user")
		return ErrUnexpected
	}

	return nil
}

func (uc *LoginByWidget) acceptLoginRequest(ctx context.Context, loginChallenge string, userId int64) (*hydra.CompletedRequest, error) {
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

func (uc *LoginByWidget) mapRejectError(err error) (string, int64, string) {
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

func (uc *LoginByWidget) rejectLoginRequest(ctx context.Context, loginChallenge string, reason error) (string, error) {
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
			return "", NewGatewayTimeoutErr("hydra")
		}
		if resp != nil {
			if resp.StatusCode >= http.StatusInternalServerError {
				return "", NewBadGatewayErr("hydra")
			}
			return "", ErrUnexpected
		}
		return "", NewBadGatewayErr("hydra")
	}

	if completed == nil || completed.RedirectTo == "" {
		return "", ErrUnexpected
	}

	return completed.RedirectTo, nil
}

func (uc *LoginByWidget) rejectAndBuildOutput(ctx context.Context, loginChallenge string, reason error) (*LoginByWidgetOutput, error) {
	redirectUri, rejectErr := uc.rejectLoginRequest(ctx, loginChallenge, reason)
	if rejectErr != nil {
		return nil, rejectErr
	}

	return &LoginByWidgetOutput{RedirectUri: redirectUri}, nil
}

func (uc *LoginByWidget) Execute(ctx context.Context, input *LoginByWidgetInput) (*LoginByWidgetOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if err := uc.verifyChallenge(input.LoginChallenge); err != nil {
		return nil, err
	}

	if err := uc.verifyIP(input.ClientIP); err != nil {
		return nil, err
	}

	loginRequest, err := uc.getLoginRequest(ctx, input.LoginChallenge)
	if err != nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, err)
	}
	if loginRequest == nil || loginRequest.Client.ClientId == nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, ErrUnexpected)
	}

	clientId := *loginRequest.Client.ClientId

	bot, err := uc.getBot(ctx, clientId)
	if err != nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, err)
	}

	if err := uc.verifyBotToken(ctx, bot.Token); err != nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, err)
	}

	authData, err := uc.parseAndVerifyAuthData(input.AuthData, bot.Token)
	if err != nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, err)
	}

	if err := uc.transactor.RunInTransaction(ctx, func(txCtx context.Context) error {
		return uc.ensureBotUserExists(
			txCtx,
			bot.Id,
			authData.User,
			input.ClientIP,
			input.UserAgent,
			input.Language,
		)
	}); err != nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, err)
	}

	completed, err := uc.acceptLoginRequest(ctx, input.LoginChallenge, authData.User.Id)
	if err != nil {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, err)
	}

	if completed == nil || completed.RedirectTo == "" {
		return uc.rejectAndBuildOutput(ctx, input.LoginChallenge, ErrUnexpected)
	}

	return &LoginByWidgetOutput{RedirectUri: completed.RedirectTo}, nil
}
