package repositories

import (
	"context"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type BotRepository interface {
	Create(
		ctx context.Context,
		bot *entities.Bot,
	) error

	Read(
		ctx context.Context,
		id int64,
		bot *entities.Bot,
	) error

	ReadByClientId(
		ctx context.Context,
		clientId string,
		bot *entities.Bot,
	) error

	Update(
		ctx context.Context,
		bot *entities.Bot,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}
