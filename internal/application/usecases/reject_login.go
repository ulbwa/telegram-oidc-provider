package usecases

import (
	"context"
	"net/url"

	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
)

type RejectLogin struct {
	hydra app_services.HydraLoginManager
}

type RejectLoginInput struct {
	LoginChallenge string
}

type RejectLoginOutput struct {
	RedirectTo *url.URL
}

func (uc *RejectLogin) Execute(ctx context.Context, input *RejectLoginInput) (output *RejectLoginOutput, err error) {
	rejectLoginRequestResp, err := uc.hydra.RejectLoginRequest(
		ctx,
		input.LoginChallenge,
		app_services.OAuth2ErrorAccessDenied,
	)
	if err != nil {
		return nil, err
	}

	output = &RejectLoginOutput{
		RedirectTo: rejectLoginRequestResp.RedirectTo,
	}

	return output, nil

}
