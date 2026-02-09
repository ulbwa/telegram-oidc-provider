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

type LoginByTelegramMiniApp struct {
	ttl time.Duration

	userRepo  repositories.UserRepository
	loginRepo repositories.UserBotLoginRepository
	botRepo   repositories.BotRepository

	hydra       app_services.HydraLoginManager
	transactor  app_services.Transactor
	parser      app_services.TelegramInitDataParser
	verifier    app_services.TelegramHashVerifier
	replayGuard app_services.TelegramReplayGuard
	userFactory app_services.TelegramUserFactory
}

type LoginByTelegramMiniAppInput struct {
	UserIp         net.IP
	UserAgent      *string
	AcceptLanguage *string
	LoginChallenge string
	InitData       string
}
type LoginByTelegramMiniAppOutput struct {
	RedirectTo *url.URL
}

func (uc *LoginByTelegramMiniApp) upsertUser(ctx context.Context, userData *app_services.TelegramUserData) (*domain.User, error) {
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
		user.SetFirstName(userData.FirstName)
		user.SetLastName(userData.LastName)
		user.SetUsername(userData.Username)
		user.SetPhotoUrl(userData.PhotoUrl)
		user.SetIsPremium(userData.IsPremium)
		if err := uc.userRepo.Update(ctx, &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (uc *LoginByTelegramMiniApp) accept(ctx context.Context, input *LoginByTelegramMiniAppInput, user *domain.User) (output *LoginByTelegramMiniAppOutput, err error) {
	acceptLoginRequestResp, err := uc.hydra.AcceptLoginRequest(ctx, input.LoginChallenge, user, false)
	if err != nil {
		return nil, err
	}
	output = &LoginByTelegramMiniAppOutput{
		RedirectTo: acceptLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginByTelegramMiniApp) reject(
	ctx context.Context,
	input *LoginByTelegramMiniAppInput,
	code app_services.OAuth2ErrorCode,
) (output *LoginByTelegramMiniAppOutput, err error) {
	rejectLoginRequestResp, err := uc.hydra.RejectLoginRequest(ctx, input.LoginChallenge, code)
	if err != nil {
		return nil, err
	}
	output = &LoginByTelegramMiniAppOutput{
		RedirectTo: rejectLoginRequestResp.RedirectTo,
	}

	return output, nil
}

func (uc *LoginByTelegramMiniApp) upsertUserBotLogin(ctx context.Context, user *domain.User, bot *domain.Bot, input *LoginByTelegramMiniAppInput, authData *app_services.TelegramAuthData) error {
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
		userLogin.SetIP(input.UserIp)
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

func (uc *LoginByTelegramMiniApp) Execute(ctx context.Context, input *LoginByTelegramMiniAppInput) (output *LoginByTelegramMiniAppOutput, err error) {
	loginRequest, err := uc.hydra.GetLoginRequest(ctx, input.LoginChallenge)
	if err != nil {
		return nil, err
	}
	if loginRequest.SkipRequired() {
		return uc.reject(ctx, input, app_services.OAuth2ErrorServerError)
	}

	authData, err := uc.parser.Parse(input.InitData)
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
