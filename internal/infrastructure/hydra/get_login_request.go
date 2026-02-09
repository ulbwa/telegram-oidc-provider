package hydra

import (
	"context"
	"errors"
	"net/url"

	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
)

func (c *HydraClient) GetLoginRequest(
	ctx context.Context,
	challenge string,
) (*app_services.LoginRequestResponse, error) {
	hydraResp, _, err := c.client.AdminApi.
		GetLoginRequest(ctx).
		LoginChallenge(challenge).
		Execute()
	if err != nil {
		// TODO: handle error
		return nil, err
	}

	requestUrl, err := url.Parse(hydraResp.RequestUrl)
	if err != nil {
		return nil, err
	}
	if hydraResp.Client.ClientId == nil {
		return nil, errors.New("hydra response missing client ID")
	}
	if hydraResp.SessionId == nil {
		return nil, errors.New("hydra response missing session ID")
	}
	var subject *string
	if hydraResp.Subject != "" {
		subject = &hydraResp.Subject
	}

	var resp app_services.LoginRequestResponse
	resp.RequestUrl = requestUrl
	resp.ClientId = *hydraResp.Client.ClientId
	resp.Subject = subject
	resp.SessionId = *hydraResp.SessionId
	resp.Skip = hydraResp.Skip
	return &resp, nil
}
