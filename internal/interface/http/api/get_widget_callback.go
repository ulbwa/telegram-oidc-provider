package api

import (
	"context"
	"errors"

	"github.com/ulbwa/telegram-oidc-provider/api/generated"
)

// Login user by telegram widget auth data
// (GET /widget/callback)
func (s *server) GetWidgetCallback(ctx context.Context, request generated.GetWidgetCallbackRequestObject) (generated.GetWidgetCallbackResponseObject, error) {
	return nil, errors.New("not implemented")
}
