package repositories

import (
	"context"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type UserBotLoginRepository interface {
	Create(
		ctx context.Context,
		login *domain.UserBotLogin,
	) error

	Read(
		ctx context.Context,
		userId,
		botId int64,
		login *domain.UserBotLogin,
	) error

	Update(
		ctx context.Context,
		login *domain.UserBotLogin,
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
		logins *[]domain.UserBotLogin,
	) error
}
