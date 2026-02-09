package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

// botRepository implements the BotRepository interface using GORM
type botRepository struct {
	db *gorm.DB
}

// NewBotRepository creates a new BotRepository instance
func NewBotRepository(db *gorm.DB) repositories.BotRepository {
	return &botRepository{
		db: db,
	}
}

// Create creates a new bot in the database
func (r *botRepository) Create(ctx context.Context, bot *domain.Bot) error {
	if bot == nil {
		return fmt.Errorf("%w: bot cannot be nil", repositories.ErrInvalidArgument)
	}

	model := r.domainToModel(bot)
	tx := GetTx(ctx, r.db)

	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return mapError(err, "failed to create bot")
	}

	// Update the domain entity with the created data
	r.modelToDomain(&model, bot)

	return nil
}

// Read retrieves a bot by its ID
func (r *botRepository) Read(ctx context.Context, id int64, bot *domain.Bot) error {
	if bot == nil {
		return fmt.Errorf("%w: bot cannot be nil", repositories.ErrInvalidArgument)
	}

	var model BotModel
	tx := GetTx(ctx, r.db)

	if err := tx.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return mapError(err, fmt.Sprintf("failed to read bot with id %d", id))
	}

	r.modelToDomain(&model, bot)
	return nil
}

// ReadByClientId retrieves a bot by its client ID
func (r *botRepository) ReadByClientId(ctx context.Context, clientId string, bot *domain.Bot) error {
	if bot == nil {
		return fmt.Errorf("%w: bot cannot be nil", repositories.ErrInvalidArgument)
	}

	var model BotModel
	tx := GetTx(ctx, r.db)

	if err := tx.WithContext(ctx).Where("client_id = ?", clientId).First(&model).Error; err != nil {
		return mapError(err, fmt.Sprintf("failed to read bot by client_id %s", clientId))
	}

	r.modelToDomain(&model, bot)
	return nil
}

// Update updates an existing bot in the database
func (r *botRepository) Update(ctx context.Context, bot *domain.Bot) error {
	if bot == nil {
		return fmt.Errorf("%w: bot cannot be nil", repositories.ErrInvalidArgument)
	}

	model := r.domainToModel(bot)
	tx := GetTx(ctx, r.db)

	result := tx.WithContext(ctx).Model(&BotModel{}).Where("id = ?", model.Id).Updates(map[string]interface{}{
		"name":       model.Name,
		"client_id":  model.ClientId,
		"username":   model.Username,
		"token":      model.Token,
		"updated_at": model.UpdatedAt,
	})

	if result.Error != nil {
		return mapError(result.Error, fmt.Sprintf("failed to update bot with id %d", model.Id))
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: bot with id %d not found for update", repositories.ErrNotFound, model.Id)
	}

	return nil
}

// Delete deletes a bot from the database
func (r *botRepository) Delete(ctx context.Context, id int64) error {
	tx := GetTx(ctx, r.db)

	result := tx.WithContext(ctx).Where("id = ?", id).Delete(&BotModel{})
	if result.Error != nil {
		return mapError(result.Error, fmt.Sprintf("failed to delete bot with id %d", id))
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: bot with id %d not found for deletion", repositories.ErrNotFound, id)
	}

	return nil
}

// domainToModel converts a domain Bot entity to a BotModel
func (r *botRepository) domainToModel(bot *domain.Bot) BotModel {
	model := BotModel{
		Id:        bot.Id,
		Name:      bot.Name,
		ClientId:  bot.ClientId,
		Username:  bot.Username,
		Token:     bot.Token,
		CreatedAt: bot.CreatedAt,
	}

	if bot.UpdatedAt != nil {
		model.UpdatedAt = sql.NullTime{
			Time:  *bot.UpdatedAt,
			Valid: true,
		}
	}

	return model
}

// modelToDomain converts a BotModel to a domain Bot entity
func (r *botRepository) modelToDomain(model *BotModel, bot *domain.Bot) {
	bot.Id = model.Id
	bot.Name = model.Name
	bot.ClientId = model.ClientId
	bot.Username = model.Username
	bot.Token = model.Token
	bot.CreatedAt = model.CreatedAt

	if model.UpdatedAt.Valid {
		bot.UpdatedAt = &model.UpdatedAt.Time
	} else {
		bot.UpdatedAt = nil
	}
}
