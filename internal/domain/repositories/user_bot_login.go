package repositories

import (
	"context"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type UserBotLoginRepository interface {
	Create(
		ctx context.Context,
		login *entities.UserBotLogin,
	) error

	Read(
		ctx context.Context,
		userId,
		botId int64,
		login *entities.UserBotLogin,
	) error

	Update(
		ctx context.Context,
		login *entities.UserBotLogin,
	) error

	Delete(
		ctx context.Context,
		userId,
		botId int64,
	) error

	ReadByBotId(
		ctx context.Context,
		botId int64,
		page *string,
		logins *[]entities.UserBotLogin,
	) error
}
