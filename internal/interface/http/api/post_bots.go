package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ulbwa/telegram-oidc-provider/api/generated"
	"github.com/ulbwa/telegram-oidc-provider/internal/application/usecase"
	"github.com/ulbwa/telegram-oidc-provider/pkg/utils"
)

// Sync Telegram bot by token
// (POST /bots)
func (s *server) PostBots(ctx context.Context, request generated.PostBotsRequestObject) (generated.PostBotsResponseObject, error) {
	input := usecase.SyncBotInput{
		BotToken: request.Body.Token,
	}
	output, err := s.syncBot.Execute(ctx, &input)
	if err != nil {
		code, resp, err := handleError(err)
		if err != nil {
			return nil, err
		}
		switch code {
		case http.StatusBadRequest:
			return generated.PostBots400JSONResponse(*resp), nil
		case http.StatusInternalServerError:
			return generated.PostBots500JSONResponse(*resp), nil
		default:
			return nil, errors.New("unexpected error code from error handler")
		}
	}

	location := s.baseUri.JoinPath("bots", strconv.FormatInt(output.Id, 10))

	if output.Status == usecase.SyncBotStatusCreated {
		var httpResp generated.PostBots201JSONResponse
		httpResp.Body.Id = output.Id
		httpResp.Headers.Location = location.String()
		return httpResp, nil
	} else {
		var httpResp generated.PostBots200JSONResponse
		httpResp.Body.Id = output.Id
		httpResp.Body.LastSyncedAt = &output.LastSyncedAt
		if output.Status == usecase.SyncBotStatusUpdated {
			httpResp.Body.Updated = utils.Ptr(true)
		}
		httpResp.Headers.Location = location.String()
		return httpResp, nil
	}
}
