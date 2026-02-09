package hydra

import (
	"context"
)

func (c *HydraClient) DeleteClient(
	ctx context.Context,
	id string,
) error {
	_, err := c.client.AdminApi.
		DeleteOAuth2Client(ctx, id).
		Execute()
	// TODO: handle error
	return err
}
