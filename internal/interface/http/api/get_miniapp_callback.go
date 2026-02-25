package api

import (
	"context"
	"errors"

	"github.com/ulbwa/telegram-oidc-provider/api/generated"
)

// Login user by telegram mini app auth data
// (GET /miniapp/callback)
func (s *server) GetMiniappCallback(ctx context.Context, request generated.GetMiniappCallbackRequestObject) (generated.GetMiniappCallbackResponseObject, error) {
	return nil, errors.New("not implemented")
}
