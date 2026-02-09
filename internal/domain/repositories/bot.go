package repositories

import (
	"context"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type BotRepository interface {
	Create(
		ctx context.Context,
		bot *domain.Bot,
	) error

	Read(
		ctx context.Context,
		id int64,
		bot *domain.Bot,
	) error

	ReadByClientId(
		ctx context.Context,
		clientId string,
		bot *domain.Bot,
	) error

	Update(
		ctx context.Context,
		bot *domain.Bot,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}
