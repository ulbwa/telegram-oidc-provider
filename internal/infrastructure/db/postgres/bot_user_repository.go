package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/db/model"
	"gorm.io/gorm"
)

// GormBotUserRepository implements port.BotUserRepositoryPort using GORM.
type GormBotUserRepository struct {
	gormDB *gorm.DB
}

// Compile-time check that GormBotUserRepository implements port.BotUserRepositoryPort
var _ repository.BotUserRepositoryPort = (*GormBotUserRepository)(nil)

// NewBotUserRepository creates a new GORM-based bot user repository.
func NewBotUserRepository(gormDB *gorm.DB) *GormBotUserRepository {
	return &GormBotUserRepository{gormDB: gormDB}
}

// toDBModel converts entity.BotUser to model.BotUser.
func (r *GormBotUserRepository) toDBModel(botUser *entity.BotUser) *model.BotUser {
	dbBotUser := &model.BotUser{
		BotId:       botUser.BotId,
		UserId:      botUser.UserId,
		FirstName:   botUser.User.FirstName,
		IP:          botUser.IP,
		LastLoginAt: botUser.LastLoginAt,
		CreatedAt:   botUser.CreatedAt,
	}

	// Handle nullable fields from User
	if botUser.User.LastName != nil {
		dbBotUser.LastName = sql.NullString{String: *botUser.User.LastName, Valid: true}
	}

	if botUser.User.Username != nil {
		dbBotUser.Username = sql.NullString{String: *botUser.User.Username, Valid: true}
	}

	if botUser.User.PhotoUrl != nil {
		dbBotUser.PhotoUrl = sql.NullString{String: botUser.User.PhotoUrl.String(), Valid: true}
	}

	if botUser.User.IsPremium != nil {
		dbBotUser.IsPremium = sql.NullBool{Bool: *botUser.User.IsPremium, Valid: true}
	}

	if botUser.UserAgent != nil {
		dbBotUser.UserAgent = sql.NullString{String: *botUser.UserAgent, Valid: true}
	}

	if botUser.Language != nil {
		dbBotUser.Language = sql.NullString{String: *botUser.Language, Valid: true}
	}

	if botUser.UpdatedAt != nil {
		dbBotUser.UpdatedAt = sql.NullTime{Time: *botUser.UpdatedAt, Valid: true}
	}

	return dbBotUser
}

// toEntity converts model.BotUser to entity.BotUser.
func (r *GormBotUserRepository) toEntity(dbBotUser *model.BotUser) (*entity.BotUser, error) {
	user := entity.User{
		FirstName: dbBotUser.FirstName,
	}

	if dbBotUser.LastName.Valid {
		user.LastName = &dbBotUser.LastName.String
	}

	if dbBotUser.Username.Valid {
		user.Username = &dbBotUser.Username.String
	}

	if dbBotUser.PhotoUrl.Valid {
		photoUrl, err := url.Parse(dbBotUser.PhotoUrl.String)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid photo URL in database: %v", repository.ErrCorruptedData, err)
		}
		user.PhotoUrl = photoUrl
	}

	if dbBotUser.IsPremium.Valid {
		user.IsPremium = &dbBotUser.IsPremium.Bool
	}

	botUser := &entity.BotUser{
		BotId:       dbBotUser.BotId,
		UserId:      dbBotUser.UserId,
		User:        user,
		IP:          dbBotUser.IP,
		LastLoginAt: dbBotUser.LastLoginAt,
		CreatedAt:   dbBotUser.CreatedAt,
	}

	if dbBotUser.UserAgent.Valid {
		botUser.UserAgent = &dbBotUser.UserAgent.String
	}

	if dbBotUser.Language.Valid {
		botUser.Language = &dbBotUser.Language.String
	}

	if dbBotUser.UpdatedAt.Valid {
		botUser.UpdatedAt = &dbBotUser.UpdatedAt.Time
	}

	return botUser, nil
}

// GetByBotAndUser retrieves a bot user by bot ID and user ID and populates the provided botUser pointer.
func (r *GormBotUserRepository) GetByBotAndUser(ctx context.Context, botID, userID int64, botUser *entity.BotUser) error {
	gormDB := GetTx(ctx, r.gormDB)

	var dbBotUser model.BotUser
	if err := gormDB.WithContext(ctx).Where("bot_id = ? AND user_id = ?", botID, userID).First(&dbBotUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	result, err := r.toEntity(&dbBotUser)
	if err != nil {
		return err
	}

	*botUser = *result
	return nil
}

// GetByBot retrieves all users for a specific bot.
func (r *GormBotUserRepository) GetByBot(ctx context.Context, botID int64) ([]*entity.BotUser, error) {
	gormDB := GetTx(ctx, r.gormDB)

	var dbBotUsers []model.BotUser
	if err := gormDB.WithContext(ctx).Where("bot_id = ?", botID).Find(&dbBotUsers).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	botUsers := make([]*entity.BotUser, 0, len(dbBotUsers))
	for i := range dbBotUsers {
		botUser, err := r.toEntity(&dbBotUsers[i])
		if err != nil {
			return nil, err
		}
		botUsers = append(botUsers, botUser)
	}

	return botUsers, nil
}

// Create stores a new bot user and updates the provided botUser pointer with inserted data.
func (r *GormBotUserRepository) Create(ctx context.Context, botUser *entity.BotUser) error {
	gormDB := GetTx(ctx, r.gormDB)

	dbBotUser := r.toDBModel(botUser)

	if err := gormDB.WithContext(ctx).Create(dbBotUser).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%w: bot user already exists", repository.ErrDuplicate)
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fmt.Errorf("%w: bot does not exist", repository.ErrNotFound)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	// Reload from DB to get all fields including defaults
	result, err := r.toEntity(dbBotUser)
	if err != nil {
		return err
	}

	*botUser = *result
	return nil
}

// Update updates an existing bot user and refreshes the provided botUser pointer.
func (r *GormBotUserRepository) Update(ctx context.Context, botUser *entity.BotUser) error {
	gormDB := GetTx(ctx, r.gormDB)

	dbBotUser := r.toDBModel(botUser)

	result := gormDB.WithContext(ctx).
		Model(&model.BotUser{}).
		Where("bot_id = ? AND user_id = ?", botUser.BotId, botUser.UserId).
		Updates(dbBotUser)

	if result.Error != nil {
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	// Reload from DB to get updated fields
	var updated model.BotUser
	if err := gormDB.WithContext(ctx).
		Where("bot_id = ? AND user_id = ?", botUser.BotId, botUser.UserId).
		First(&updated).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	reloaded, err := r.toEntity(&updated)
	if err != nil {
		return err
	}

	*botUser = *reloaded
	return nil
}

// Delete removes a bot user.
func (r *GormBotUserRepository) Delete(ctx context.Context, botID, userID int64) error {
	gormDB := GetTx(ctx, r.gormDB)

	result := gormDB.WithContext(ctx).
		Where("bot_id = ? AND user_id = ?", botID, userID).
		Delete(&model.BotUser{})

	if result.Error != nil {
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
