package hydra

import (
	"context"
)

func (c *HydraClient) RevokeLoginSession(
	ctx context.Context,
	sessionId string,
) error {
	_, err := c.client.AdminApi.
		RevokeAuthenticationSession(ctx).
		Subject(sessionId).
		Execute()
		// TODO: handle error
	return err
}
