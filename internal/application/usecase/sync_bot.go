//go:generate go-enum --values --names --nocase
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ulbwa/telegram-oidc-provider/internal/application/service"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/pkg/utils"
)

type SyncBot struct {
	transactor    service.Transactor
	botRepo       repository.BotRepositoryPort
	tokenVerifier service.TelegramTokenVerifier
}

func NewSyncBot(
	transactor service.Transactor,
	botRepo repository.BotRepositoryPort,
	tokenVerifier service.TelegramTokenVerifier,
) *SyncBot {
	return &SyncBot{
		transactor:    transactor,
		botRepo:       botRepo,
		tokenVerifier: tokenVerifier,
	}
}

type (
	// SyncBotStatus enum for bot synchronization status
	// ENUM(
	// Created
	// Updated
	// NotUpdated
	// )
	SyncBotStatus string

	SyncBotInput struct {
		BotToken string
	}
	SyncBotOutput struct {
		Id           int64
		Status       SyncBotStatus
		LastSyncedAt time.Time
	}
)

func (uc *SyncBot) verifyBotToken(ctx context.Context, botToken string) (*service.TelegramBotInfo, error) {
	if botInfo, err := uc.tokenVerifier.Verify(ctx, botToken, service.NewVerifyOptions(service.WithSkipCacheRead())); err != nil {
		if errors.Is(err, service.ErrTelegramBotTokenMalformed) {
			return nil, fmt.Errorf(
				"%w: %w",
				ErrInvalidInput,
				NewObjectInvalidErr("bot", "token", utils.Ptr("malformed")))
		}
		if errors.Is(err, service.ErrTelegramBotTokenInvalid) {
			return nil, fmt.Errorf(
				"%w: %w",
				ErrInvalidInput,
				NewObjectInvalidErr("bot", "token", nil))
		}
		return nil, ErrUnexpected
	} else {
		return botInfo, nil
	}
}

func (uc *SyncBot) createBot(ctx context.Context, botInfo *service.TelegramBotInfo, botToken string) (*entity.Bot, error) {
	if botInfo == nil {
		return nil, errors.New("bot info is nil")
	}

	bot, err := entity.NewBot(botInfo.Id, botInfo.Name, botInfo.Username, botToken)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create bot entity", ErrUnexpected)
	}

	if err := uc.botRepo.Create(ctx, bot); err != nil {
		return nil, fmt.Errorf("%w: failed to create bot", ErrUnexpected)
	}

	return bot, nil
}

func (uc *SyncBot) updateBot(ctx context.Context, botInfo *service.TelegramBotInfo, botToken string) (*entity.Bot, bool, error) {
	if botInfo == nil {
		return nil, false, errors.New("bot info is nil")
	}

	var bot entity.Bot
	if err := uc.botRepo.GetByID(ctx, botInfo.Id, &bot); err != nil {
		return nil, false, fmt.Errorf("%w: failed to get bot by ID", ErrUnexpected)
	}
	beforeTouch := bot.ModifiedAt()

	if err := bot.SetName(botInfo.Name); err != nil {
		return nil, false, fmt.Errorf("%w: %v", NewObjectInvalidErr("bot", "name", nil), err)
	}
	if err := bot.SetUsername(botInfo.Username); err != nil {
		return nil, false, fmt.Errorf("%w: %v", NewObjectInvalidErr("bot", "username", nil), err)
	}
	if err := bot.SetToken(botToken); err != nil {
		return nil, false, fmt.Errorf("%w: %v", NewObjectInvalidErr("bot", "token", nil), err)
	}

	afterTouch := bot.ModifiedAt()
	if afterTouch.After(beforeTouch) {
		if err := uc.botRepo.Update(ctx, &bot); err != nil {
			return nil, false, fmt.Errorf("%w: failed to update bot", ErrUnexpected)
		}
		return &bot, true, nil
	} else {
		return &bot, false, nil
	}
}

func (uc *SyncBot) upsertBot(ctx context.Context, botInfo *service.TelegramBotInfo, botToken string) (*entity.Bot, SyncBotStatus, error) {
	exists, err := uc.botRepo.ExistsByID(ctx, botInfo.Id)
	if err != nil {
		return nil, "", fmt.Errorf("%w: failed to check bot existence", ErrUnexpected)
	}
	if exists {
		bot, updated, err := uc.updateBot(ctx, botInfo, botToken)
		if err != nil {
			return nil, "", err
		}
		if updated {
			return bot, SyncBotStatusUpdated, nil
		}
		return bot, SyncBotStatusNotUpdated, nil
	} else {
		bot, err := uc.createBot(ctx, botInfo, botToken)
		if err != nil {
			return nil, "", err
		}
		return bot, SyncBotStatusCreated, nil
	}
}

func (uc *SyncBot) Execute(ctx context.Context, input *SyncBotInput) (*SyncBotOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	botInfo, err := uc.verifyBotToken(ctx, input.BotToken)
	if err != nil {
		return nil, err
	}

	var output SyncBotOutput
	if err := uc.transactor.RunInTransaction(ctx, func(ctx context.Context) error {
		bot, status, err := uc.upsertBot(ctx, botInfo, input.BotToken)
		if err != nil {
			return err
		}
		output.Id = bot.Id
		output.Status = status
		output.LastSyncedAt = bot.ModifiedAt()
		return nil
	}); err != nil {
		return nil, err
	}

	return &output, nil
}
