package usecases

import (
	"context"
	"net/url"

	"github.com/rs/zerolog"

	app_services "github.com/ulbwa/telegram-oidc-provider/internal/application/services"
	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

type CreateClient struct {
	botVerifier app_services.TelegramBotVerifier
	hydraClient app_services.HydraClientManager
	transactor  app_services.Transactor

	botRepo repositories.BotRepository
}

func NewCreateClient(
	botVerifier app_services.TelegramBotVerifier,
	hydraClient app_services.HydraClientManager,
	transactor app_services.Transactor,

	botRepo repositories.BotRepository,
) *CreateClient {
	return &CreateClient{
		botVerifier: botVerifier,
		hydraClient: hydraClient,
		transactor:  transactor,
		botRepo:     botRepo,
	}
}

type CreateClientInput struct {
	Name        string
	RedirectUri *url.URL
	BotToken    string
}
type CreateClientOutput struct {
	Id string
}

func (uc *CreateClient) deleteClientOnError(ctx context.Context, clientId string) {
	ctx = context.WithoutCancel(ctx)
	if err := uc.hydraClient.DeleteClient(ctx, clientId); err != nil {
		zerolog.
			Ctx(ctx).
			Err(err).
			Str("client_id", clientId).
			Msg("failed to delete client after error")
	}
}

func (uc *CreateClient) Execute(ctx context.Context, input *CreateClientInput) (output *CreateClientOutput, err error) {
	botData, err := uc.botVerifier.Verify(ctx, input.BotToken)
	if err != nil {
		return nil, err
	}

	client, err := uc.hydraClient.CreateClient(ctx, input.Name, input.RedirectUri)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			uc.deleteClientOnError(ctx, client.Id)
			panic(r)
		}
		if err != nil {
			uc.deleteClientOnError(ctx, client.Id)
		}
	}()

	var bot *domain.Bot
	if err := uc.transactor.RunInTransaction(ctx, func(ctx context.Context) error {
		var err error
		bot, err = domain.NewBot(botData.Id, botData.Name, client.Id, botData.Username, input.BotToken)
		if err != nil {
			return err
		}

		return uc.botRepo.Create(ctx, bot)
	}); err != nil {
		return nil, err
	}

	output = &CreateClientOutput{
		Id: client.Id,
	}

	return nil, nil
}
