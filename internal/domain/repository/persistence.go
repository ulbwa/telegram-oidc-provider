package repository

import (
	"context"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
)

// BotRepositoryPort defines the interface for bot data access
type BotRepositoryPort interface {
	// GetByID retrieves a bot by its ID and populates the provided bot pointer.
	GetByID(ctx context.Context, id int64, bot *entity.Bot) error

	// GetByClientID retrieves a bot by its client ID (unique) and populates the provided bot pointer.
	GetByClientID(ctx context.Context, clientID string, bot *entity.Bot) error

	// Create stores a new bot and populates the pointer with inserted data.
	Create(ctx context.Context, bot *entity.Bot) error

	// Update updates an existing bot and refreshes the provided pointer.
	Update(ctx context.Context, bot *entity.Bot) error

	// Delete removes a bot by ID.
	Delete(ctx context.Context, id int64) error

	// ExistsByID checks if a bot exists by its ID.
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

// BotUserRepositoryPort defines the interface for bot_user data access
type BotUserRepositoryPort interface {
	// GetByBotAndUser retrieves a bot user by bot ID and user ID and populates the provided botUser pointer.
	GetByBotAndUser(ctx context.Context, botID, userID int64, botUser *entity.BotUser) error

	// GetByBot retrieves all users for a specific bot.
	GetByBot(ctx context.Context, botID int64) ([]*entity.BotUser, error)

	// Create stores a new bot user and populates the pointer with inserted data.
	Create(ctx context.Context, botUser *entity.BotUser) error

	// Update updates an existing bot user and refreshes the provided pointer.
	Update(ctx context.Context, botUser *entity.BotUser) error

	// Delete removes a bot user.
	Delete(ctx context.Context, botID, userID int64) error
}
