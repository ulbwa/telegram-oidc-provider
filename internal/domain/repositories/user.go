package repositories

import (
	"context"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user *entities.User,
	) error

	Read(
		ctx context.Context,
		id int64,
		user *entities.User,
	) error

	Update(
		ctx context.Context,
		user *entities.User,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}
