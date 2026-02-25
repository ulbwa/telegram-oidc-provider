package usecase

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	hydraclient "github.com/ory/hydra-client-go"
// 	"github.com/rs/zerolog"
// 	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
// )

// type LinkBotToClient struct {
// 	botRepo repository.BotRepositoryPort
// 	hydra   *hydraclient.APIClient
// }

// type (
// 	LinkBotToClientInput struct {
// 		BotId    int64
// 		ClientID string
// 	}
// 	LinkBotToClientOutput struct {
// 	}
// )

// func (uc *LinkBotToClient) checkClientExists(ctx context.Context, clientId string) error {
// 	client, resp, err := uc.hydra.AdminApi.GetOAuth2Client(ctx, clientId).Execute()
// 	if err != nil {
// 		if errors.Is(err, context.DeadlineExceeded) {
// 			return NewGatewayTimeoutErr("hydra")
// 		}
// 		if resp == nil {
// 			return NewBadGatewayErr("hydra")
// 		}
// 		if resp.StatusCode == 404 {
// 			return NewObjectNotFoundErr("client", clientId)
// 		}
// 		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
// 			return NewBadGatewayErr("hydra")
// 		}
// 		return fmt.Errorf("%w: failed to check client existence", ErrUnexpected)
// 	}
// 	zerolog.
// 		Ctx(ctx).
// 		Debug().
// 		Str("client_id", *client.ClientId).
// 		Msg("client found in hydra")
// 	return nil
// }

// func (uc *LinkBotToClient) Execute(ctx context.Context, input *LinkBotToClientInput) (*LinkBotToClientOutput, error) {
// 	if input == nil {
// 		return nil, errors.New("input is nil")
// 	}

// 	if err := uc.checkClientExists(ctx, input.ClientID); err != nil {
// 		return nil, err
// 	}

// }
