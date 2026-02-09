package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

// userBotLoginRepository implements the UserBotLoginRepository interface using GORM
type userBotLoginRepository struct {
	db *gorm.DB
}

// NewUserBotLoginRepository creates a new UserBotLoginRepository instance
func NewUserBotLoginRepository(db *gorm.DB) repositories.UserBotLoginRepository {
	return &userBotLoginRepository{
		db: db,
	}
}

// Create creates a new user bot login in the database
func (r *userBotLoginRepository) Create(ctx context.Context, login *entities.UserBotLogin) error {
	if login == nil {
		return fmt.Errorf("%w: login cannot be nil", repositories.ErrInvalidArgument)
	}

	model := r.domainToModel(login)
	tx := GetTx(ctx, r.db)

	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return mapError(err, fmt.Sprintf("failed to create user_bot_login for user %d and bot %d", login.UserId, login.BotId))
	}

	// Update the domain entity with the created data
	r.modelToDomain(&model, login)

	return nil
}

// Read retrieves a user bot login by user ID and bot ID
func (r *userBotLoginRepository) Read(ctx context.Context, userId, botId int64, login *entities.UserBotLogin) error {
	if login == nil {
		return fmt.Errorf("%w: login cannot be nil", repositories.ErrInvalidArgument)
	}

	var model UserBotLoginModel
	tx := GetTx(ctx, r.db)

	if err := tx.WithContext(ctx).Where("user_id = ? AND bot_id = ?", userId, botId).First(&model).Error; err != nil {
		return mapError(err, fmt.Sprintf("failed to read user_bot_login for user %d and bot %d", userId, botId))
	}

	r.modelToDomain(&model, login)
	return nil
}

// Update updates an existing user bot login in the database
func (r *userBotLoginRepository) Update(ctx context.Context, login *entities.UserBotLogin) error {
	if login == nil {
		return fmt.Errorf("%w: login cannot be nil", repositories.ErrInvalidArgument)
	}

	model := r.domainToModel(login)
	tx := GetTx(ctx, r.db)

	result := tx.WithContext(ctx).Model(&UserBotLoginModel{}).
		Where("user_id = ? AND bot_id = ?", model.UserId, model.BotId).
		Updates(map[string]interface{}{
			"ip":            model.IP,
			"user_agent":    model.UserAgent,
			"language":      model.Language,
			"last_login_at": model.LastLoginAt,
			"updated_at":    model.UpdatedAt,
		})

	if result.Error != nil {
		return mapError(result.Error, fmt.Sprintf("failed to update user_bot_login for user %d and bot %d", model.UserId, model.BotId))
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: user_bot_login for user %d and bot %d not found for update", repositories.ErrNotFound, model.UserId, model.BotId)
	}

	return nil
}

// Delete deletes a user bot login from the database
func (r *userBotLoginRepository) Delete(ctx context.Context, userId, botId int64) error {
	tx := GetTx(ctx, r.db)

	result := tx.WithContext(ctx).Where("user_id = ? AND bot_id = ?", userId, botId).Delete(&UserBotLoginModel{})
	if result.Error != nil {
		return mapError(result.Error, fmt.Sprintf("failed to delete user_bot_login for user %d and bot %d", userId, botId))
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: user_bot_login for user %d and bot %d not found for deletion", repositories.ErrNotFound, userId, botId)
	}

	return nil
}

// ReadByBotId retrieves all user bot logins for a specific bot with pagination
func (r *userBotLoginRepository) ReadByBotId(ctx context.Context, botId int64, page *string, logins *[]entities.UserBotLogin) error {
	if logins == nil {
		return fmt.Errorf("%w: logins cannot be nil", repositories.ErrInvalidArgument)
	}

	tx := GetTx(ctx, r.db)
	query := tx.WithContext(ctx).Where("bot_id = ?", botId).Order("last_login_at DESC")

	// Apply cursor-based pagination if page token is provided
	if page != nil && *page != "" {
		// Assuming page token is the last_login_at timestamp in RFC3339 format
		query = query.Where("last_login_at < ?", *page)
	}

	// Limit the number of results (adjust as needed)
	const pageSize = 100
	query = query.Limit(pageSize)

	var models []UserBotLoginModel
	if err := query.Find(&models).Error; err != nil {
		return mapError(err, fmt.Sprintf("failed to read user_bot_logins for bot %d", botId))
	}

	// Convert models to domain entities
	*logins = make([]entities.UserBotLogin, 0, len(models))
	for _, model := range models {
		var login entities.UserBotLogin
		r.modelToDomain(&model, &login)
		*logins = append(*logins, login)
	}

	return nil
}

// domainToModel converts a domain UserBotLogin entity to a UserBotLoginModel
func (r *userBotLoginRepository) domainToModel(login *entities.UserBotLogin) UserBotLoginModel {
	model := UserBotLoginModel{
		UserId:      login.UserId,
		BotId:       login.BotId,
		LastLoginAt: login.LastLoginAt,
		CreatedAt:   login.CreatedAt,
	}

	model.FromIP(login.IP)

	if login.UserAgent != nil {
		model.UserAgent = sql.NullString{String: *login.UserAgent, Valid: true}
	}

	if login.Language != nil {
		model.Language = sql.NullString{String: *login.Language, Valid: true}
	}

	if login.UpdatedAt != nil {
		model.UpdatedAt = sql.NullTime{Time: *login.UpdatedAt, Valid: true}
	}

	return model
}

// modelToDomain converts a UserBotLoginModel to a domain UserBotLogin entity
func (r *userBotLoginRepository) modelToDomain(model *UserBotLoginModel, login *entities.UserBotLogin) {
	login.UserId = model.UserId
	login.BotId = model.BotId
	login.IP = model.ToIP()
	login.LastLoginAt = model.LastLoginAt
	login.CreatedAt = model.CreatedAt

	if model.UserAgent.Valid {
		login.UserAgent = &model.UserAgent.String
	} else {
		login.UserAgent = nil
	}

	if model.Language.Valid {
		login.Language = &model.Language.String
	} else {
		login.Language = nil
	}

	if model.UpdatedAt.Valid {
		login.UpdatedAt = &model.UpdatedAt.Time
	} else {
		login.UpdatedAt = nil
	}
}
