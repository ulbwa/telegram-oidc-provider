package hydra

import (
	"context"
	"fmt"
	"net/url"

	hydra_client "github.com/ory/hydra-client-go"
	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

func (c *HydraClient) AcceptLoginRequest(
	ctx context.Context,
	challenge string,
	user *domain.User,
	remember bool,
) (*app_services.AcceptedLoginRequestResponse, error) {
	body := hydra_client.NewAcceptLoginRequest(fmt.Sprintf("%d", user.Id))
	body.SetRemember(remember)
	body.SetRememberFor(3600)

	hydraResp, _, err := c.client.AdminApi.
		AcceptLoginRequest(ctx).
		LoginChallenge(challenge).
		AcceptLoginRequest(*body).
		Execute()
	if err != nil {
		// TODO: handle error
		return nil, err
	}

	redirectTo, err := url.Parse(hydraResp.RedirectTo)
	if err != nil {
		return nil, err
	}

	return &app_services.AcceptedLoginRequestResponse{
		RedirectTo: redirectTo,
	}, nil
}
