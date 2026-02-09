package usecases

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strconv"
	"time"

	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

type LoginBySession struct {
	userRepo  repositories.UserRepository
	loginRepo repositories.UserBotLoginRepository
	botRepo   repositories.BotRepository

	hydra      app_services.HydraLoginManager
	transactor app_services.Transactor
}

type LoginBySessionInput struct {
	UserIp         net.IP
	UserAgent      *string
	AcceptLanguage *string
	LoginChallenge string
}
type LoginBySessionOutput struct {
	RedirectTo *url.URL
}

func (uc *LoginBySession) upsertUserBotLogin(ctx context.Context, user *domain.User, bot *domain.Bot, input *LoginBySessionInput) error {
	var userLogin domain.UserBotLogin
	if err := uc.loginRepo.Read(ctx, user.Id, bot.Id, &userLogin); err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			return err
		}
		userLogin = domain.UserBotLogin{
			UserId:      user.Id,
			BotId:       bot.Id,
			IP:          input.UserIp,
			UserAgent:   input.UserAgent,
			Language:    input.AcceptLanguage,
			LastLoginAt: time.Now(),
		}
		if err := uc.loginRepo.Create(ctx, &userLogin); err != nil {
			return err
		}
	} else {
		if err := userLogin.SetIP(input.UserIp); err != nil {
			return err
		}
		if input.UserAgent != nil {
			userLogin.SetUserAgent(input.UserAgent)
		}
		if input.AcceptLanguage != nil {
			userLogin.SetLanguage(input.AcceptLanguage)
		}
		userLogin.UpdateLastLogin()
		if err := uc.loginRepo.Update(ctx, &userLogin); err != nil {
			return err
		}
	}

	return nil
}

func (uc *LoginBySession) accept(
	ctx context.Context,
	input *LoginBySessionInput,
	user *domain.User,
) (output *LoginBySessionOutput, err error) {
	acceptLoginRequestResp, err := uc.hydra.AcceptLoginRequest(ctx, input.LoginChallenge, user, true)
	if err != nil {
		return nil, err
	}
	output = &LoginBySessionOutput{
		RedirectTo: acceptLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginBySession) reject(
	ctx context.Context,
	input *LoginBySessionInput,
	code app_services.OAuth2ErrorCode,
) (output *LoginBySessionOutput, err error) {
	rejectLoginRequestResp, err := uc.hydra.RejectLoginRequest(ctx, input.LoginChallenge, code)
	if err != nil {
		return nil, err
	}
	output = &LoginBySessionOutput{
		RedirectTo: rejectLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginBySession) refresh(
	ctx context.Context,
	input *LoginBySessionInput,
	loginRequest *app_services.LoginRequestResponse,
) (output *LoginBySessionOutput, err error) {
	if err := uc.hydra.RevokeLoginSession(ctx, loginRequest.SessionId); err != nil {
		return nil, err
	}

	rejectLoginRequestResp, err := uc.hydra.RejectLoginRequest(ctx, input.LoginChallenge, app_services.OIDCErrorLoginRequired)
	if err != nil {
		return nil, err
	}
	output = &LoginBySessionOutput{
		RedirectTo: rejectLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginBySession) Execute(ctx context.Context, input *LoginBySessionInput) (output *LoginBySessionOutput, err error) {
	loginRequest, err := uc.hydra.GetLoginRequest(ctx, input.LoginChallenge)
	if err != nil {
		return nil, err
	}
	if !loginRequest.IsAuthenticated() {
		// User is not authenticated, force re-login, without revoking login session
		return uc.reject(ctx, input, app_services.OIDCErrorLoginRequired)
	}
	userId, err := strconv.ParseInt(*loginRequest.Subject, 10, 64)
	if err != nil {
		return uc.refresh(ctx, input, loginRequest)
	}

	var refresh bool
	var reject bool
	var rejectReason app_services.OAuth2ErrorCode
	var user *domain.User

	if err := uc.transactor.RunInTransaction(ctx, func(ctx context.Context) error {
		var bot domain.Bot
		if err := uc.botRepo.ReadByClientId(ctx, loginRequest.ClientId, &bot); err != nil {
			if errors.Is(err, repositories.ErrNotFound) {
				reject = true
				rejectReason = app_services.OAuth2ErrorUnauthorizedClient
				return err
			}
			return err
		}

		if err := uc.userRepo.Read(ctx, userId, user); err != nil {
			if errors.Is(err, repositories.ErrNotFound) {
				refresh = true
				return err
			}
			return err
		}

		if err := uc.upsertUserBotLogin(ctx, user, &bot, input); err != nil {
			reject = true
			rejectReason = app_services.OAuth2ErrorServerError
			return err
		}

		return nil
	}); err != nil {
		if !refresh && !reject {
			return nil, err
		}
	}

	if refresh {
		return uc.refresh(ctx, input, loginRequest)
	} else if reject {
		return uc.reject(ctx, input, rejectReason)
	} else {
		return uc.accept(ctx, input, user)
	}
}
