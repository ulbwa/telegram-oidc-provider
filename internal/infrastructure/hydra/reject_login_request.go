package hydra

import (
	"context"
	"net/url"

	hydra_client "github.com/ory/hydra-client-go"
	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
)

func (c *HydraClient) RejectLoginRequest(
	ctx context.Context,
	challenge string,
	error app_services.OAuth2ErrorCode,
) (*app_services.RejectLoginRequestResponse, error) {
	body := hydra_client.NewRejectRequest()
	body.SetError(string(error))

	hydraResp, _, err := c.client.AdminApi.
		RejectLoginRequest(ctx).
		LoginChallenge(challenge).
		RejectRequest(*body).
		Execute()
	if err != nil {
		// TODO: handle error
		return nil, err
	}

	redirectTo, err := url.Parse(hydraResp.RedirectTo)
	if err != nil {
		return nil, err
	}

	return &app_services.RejectLoginRequestResponse{
		RedirectTo: redirectTo,
	}, nil
}
