package repositories

import (
	"context"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user *domain.User,
	) error

	Read(
		ctx context.Context,
		id int64,
		user *domain.User,
	) error

	Update(
		ctx context.Context,
		user *domain.User,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}
