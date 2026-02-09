package usecases

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

type LoginByTelegramWidget struct {
	ttl time.Duration

	userRepo  repositories.UserRepository
	loginRepo repositories.UserBotLoginRepository
	botRepo   repositories.BotRepository

	hydra       app_services.HydraLoginManager
	transactor  app_services.Transactor
	parser      app_services.TelegramWidgetDataParser
	verifier    app_services.TelegramHashVerifier
	replayGuard app_services.TelegramReplayGuard
	userFactory app_services.TelegramUserFactory
}

type LoginByTelegramWidgetInput struct {
	UserIp         net.IP
	UserAgent      *string
	AcceptLanguage *string
	LoginChallenge string
	WidgetData     string
}
type LoginByTelegramWidgetOutput struct {
	RedirectTo *url.URL
}

func (uc *LoginByTelegramWidget) upsertUser(ctx context.Context, userData *app_services.TelegramUserData) (*domain.User, error) {
	var user domain.User
	if err := uc.userRepo.Read(ctx, userData.Id, &user); err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			return nil, err
		}
		newUser, err := uc.userFactory.CreateUser(userData)
		if err != nil {
			return nil, err
		}
		user = *newUser
		if err := uc.userRepo.Create(ctx, &user); err != nil {
			return nil, err
		}
	} else {
		if err := user.SetFirstName(userData.FirstName); err != nil {
			return nil, err
		}
		if err := user.SetLastName(userData.LastName); err != nil {
			return nil, err
		}
		if err := user.SetUsername(userData.Username); err != nil {
			return nil, err
		}
		if err := user.SetPhotoUrl(userData.PhotoUrl); err != nil {
			return nil, err
		}
		if err := uc.userRepo.Update(ctx, &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (uc *LoginByTelegramWidget) accept(ctx context.Context, input *LoginByTelegramWidgetInput, user *domain.User) (output *LoginByTelegramWidgetOutput, err error) {
	acceptLoginRequestResp, err := uc.hydra.AcceptLoginRequest(ctx, input.LoginChallenge, user, true)
	if err != nil {
		return nil, err
	}
	output = &LoginByTelegramWidgetOutput{
		RedirectTo: acceptLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginByTelegramWidget) reject(
	ctx context.Context,
	input *LoginByTelegramWidgetInput,
	code app_services.OAuth2ErrorCode,
) (output *LoginByTelegramWidgetOutput, err error) {
	rejectLoginRequestResp, err := uc.hydra.RejectLoginRequest(ctx, input.LoginChallenge, code)
	if err != nil {
		return nil, err
	}
	output = &LoginByTelegramWidgetOutput{
		RedirectTo: rejectLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginByTelegramWidget) upsertUserBotLogin(ctx context.Context, user *domain.User, bot *domain.Bot, input *LoginByTelegramWidgetInput, authData *app_services.TelegramAuthData) error {
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
			Language:    authData.User.LanguageCode,
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
		if authData.User.LanguageCode != nil {
			userLogin.SetLanguage(authData.User.LanguageCode)
		} else if input.AcceptLanguage != nil {
			userLogin.SetLanguage(input.AcceptLanguage)
		}
		userLogin.UpdateLastLogin()
		if err := uc.loginRepo.Update(ctx, &userLogin); err != nil {
			return err
		}
	}

	return nil
}

func (uc *LoginByTelegramWidget) Execute(ctx context.Context, input *LoginByTelegramWidgetInput) (output *LoginByTelegramWidgetOutput, err error) {
	loginRequest, err := uc.hydra.GetLoginRequest(ctx, input.LoginChallenge)
	if err != nil {
		return nil, err
	}
	if loginRequest.SkipRequired() {
		return uc.reject(ctx, input, app_services.OAuth2ErrorServerError)
	}

	authData, err := uc.parser.Parse(input.WidgetData)
	if err != nil || authData.IsExpired(uc.ttl) {
		return uc.reject(ctx, input, app_services.OAuth2ErrorInvalidRequest)
	}

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
			reject = true
			rejectReason = app_services.OAuth2ErrorServerError
			return err
		}

		if err := uc.verifier.Verify(authData.Raw, bot.Token); err != nil {
			reject = true
			rejectReason = app_services.OAuth2ErrorInvalidRequest
			return err
		}

		if err := uc.replayGuard.CheckAndMarkUsed(ctx, authData.Hash, uc.ttl); err != nil {
			reject = true
			rejectReason = app_services.OAuth2ErrorInvalidRequest
			return err
		}

		user, err = uc.upsertUser(ctx, authData.User)
		if err != nil {
			reject = true
			rejectReason = app_services.OAuth2ErrorServerError
			return err
		}

		if err := uc.upsertUserBotLogin(ctx, user, &bot, input, authData); err != nil {
			reject = true
			rejectReason = app_services.OAuth2ErrorServerError
			return err
		}

		return nil
	}); err != nil {
		if !reject {
			return nil, err
		}
	}

	if reject {
		return uc.reject(ctx, input, rejectReason)
	} else {
		return uc.accept(ctx, input, user)
	}
}
