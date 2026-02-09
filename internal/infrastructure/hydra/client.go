package hydra

import (
	"errors"

	hydra_client "github.com/ory/hydra-client-go"
)

type HydraClient struct {
	client *hydra_client.APIClient
}

func NewHydraClient(client *hydra_client.APIClient) (*HydraClient, error) {
	if client == nil {
		return nil, errors.New("hydra client cannot be nil")
	}
	return &HydraClient{
		client: client,
	}, nil
}
