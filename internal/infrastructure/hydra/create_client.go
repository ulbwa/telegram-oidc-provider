package hydra

import (
	"context"
	"errors"
	"net/url"

	hydra_client "github.com/ory/hydra-client-go"
	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
)

func (c *HydraClient) CreateClient(
	ctx context.Context,
	name string,
	redirectUri *url.URL,
) (*app_services.CreateClientResponse, error) {
	body := hydra_client.NewOAuth2Client()
	body.SetClientName(name)
	body.SetRedirectUris([]string{redirectUri.String()})
	body.SetGrantTypes([]string{"authorization_code", "refresh_token"})
	body.SetResponseTypes([]string{"code"})
	body.SetScope("openid offline profile ")
	body.SetTokenEndpointAuthMethod("client_secret_basic")

	hydraResp, _, err := c.client.AdminApi.
		CreateOAuth2Client(ctx).
		OAuth2Client(*body).
		Execute()
	if err != nil {
		// TODO: handle error
		return nil, err
	}

	if hydraResp.ClientId == nil {
		return nil, errors.New("hydra response missing client ID")
	}
	if hydraResp.ClientSecret == nil {
		return nil, errors.New("hydra response missing client secret")
	}
	if hydraResp.ClientName == nil {
		return nil, errors.New("hydra response missing client name")
	}
	if len(hydraResp.RedirectUris) == 0 {
		return nil, errors.New("hydra response missing redirect URIs")
	}

	parsedRedirectUri, err := url.Parse(hydraResp.RedirectUris[0])
	if err != nil {
		return nil, err
	}

	return &app_services.CreateClientResponse{
		Id:          *hydraResp.ClientId,
		Secret:      *hydraResp.ClientSecret,
		Name:        *hydraResp.ClientName,
		RedirectUri: parsedRedirectUri,
	}, nil
}
