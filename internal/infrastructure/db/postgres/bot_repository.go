package postgres

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/ulbwa/telegram-oidc-provider/internal/domain/entity"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repository"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/db/model"
	"gorm.io/gorm"
)

// GormBotRepository implements port.BotRepositoryPort using GORM.
type GormBotRepository struct {
	gormDB        *gorm.DB
	encryptionKey []byte // AES-256 key (32 bytes)
}

// Compile-time check that GormBotRepository implements port.BotRepositoryPort
var _ repository.BotRepositoryPort = (*GormBotRepository)(nil)

// NewBotRepository creates a new GORM-based bot repository with token encryption.
func NewBotRepository(gormDB *gorm.DB, encryptionKey []byte) (*GormBotRepository, error) {
	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}
	return &GormBotRepository{
		gormDB:        gormDB,
		encryptionKey: encryptionKey,
	}, nil
}

// encryptToken encrypts a bot token using AES-256-GCM.
func (r *GormBotRepository) encryptToken(token string) ([]byte, error) {
	block, err := aes.NewCipher(r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create cipher: %v", repository.ErrEncryptionFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create GCM: %v", repository.ErrEncryptionFailed, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%w: failed to generate nonce: %v", repository.ErrEncryptionFailed, err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
	return ciphertext, nil
}

// decryptToken decrypts a bot token using AES-256-GCM.
func (r *GormBotRepository) decryptToken(encrypted []byte) (string, error) {
	block, err := aes.NewCipher(r.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("%w: failed to create cipher: %v", repository.ErrEncryptionFailed, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("%w: failed to create GCM: %v", repository.ErrEncryptionFailed, err)
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", fmt.Errorf("%w: token data too short", repository.ErrEncryptionFailed)
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("%w: failed to decrypt: %v", repository.ErrEncryptionFailed, err)
	}

	return string(plaintext), nil
}

// toDBModel converts entity.Bot to model.Bot with token encryption.
func (r *GormBotRepository) toDBModel(bot *entity.Bot) (*model.Bot, error) {
	encryptedToken, err := r.encryptToken(bot.Token)
	if err != nil {
		return nil, err
	}

	dbBot := &model.Bot{
		Id:        bot.Id,
		Name:      bot.Name,
		ClientId:  bot.ClientId,
		Username:  bot.Username,
		Token:     encryptedToken,
		CreatedAt: bot.CreatedAt,
	}

	if bot.UpdatedAt != nil {
		dbBot.UpdatedAt = sql.NullTime{Time: *bot.UpdatedAt, Valid: true}
	}

	return dbBot, nil
}

// toEntity converts model.Bot to entity.Bot with token decryption.
func (r *GormBotRepository) toEntity(dbBot *model.Bot) (*entity.Bot, error) {
	decryptedToken, err := r.decryptToken(dbBot.Token)
	if err != nil {
		return nil, err
	}

	bot := &entity.Bot{
		Id:        dbBot.Id,
		Name:      dbBot.Name,
		ClientId:  dbBot.ClientId,
		Username:  dbBot.Username,
		Token:     decryptedToken,
		CreatedAt: dbBot.CreatedAt,
	}

	if dbBot.UpdatedAt.Valid {
		bot.UpdatedAt = &dbBot.UpdatedAt.Time
	}

	return bot, nil
}

// ExistsByID checks whether a bot exists by id.
func (r *GormBotRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	gormDB := GetTx(ctx, r.gormDB)

	var count int64
	if err := gormDB.WithContext(ctx).Model(&model.Bot{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	return count > 0, nil
}

// GetByID retrieves a bot by its ID and populates the provided bot pointer.
func (r *GormBotRepository) GetByID(ctx context.Context, id int64, bot *entity.Bot) error {
	gormDB := GetTx(ctx, r.gormDB)

	var dbBot model.Bot
	if err := gormDB.WithContext(ctx).Where("id = ?", id).First(&dbBot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %v", repository.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	result, err := r.toEntity(&dbBot)
	if err != nil {
		return err
	}

	*bot = *result
	return nil
}

// GetByClientID retrieves a bot by its client ID and populates the provided bot pointer.
func (r *GormBotRepository) GetByClientID(ctx context.Context, clientID string, bot *entity.Bot) error {
	gormDB := GetTx(ctx, r.gormDB)

	var dbBot model.Bot
	if err := gormDB.WithContext(ctx).Where("client_id = ?", clientID).First(&dbBot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %v", repository.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	result, err := r.toEntity(&dbBot)
	if err != nil {
		return err
	}

	*bot = *result
	return nil
}

// Create stores a new bot and updates the provided bot pointer with inserted data.
func (r *GormBotRepository) Create(ctx context.Context, bot *entity.Bot) error {
	gormDB := GetTx(ctx, r.gormDB)

	dbBot, err := r.toDBModel(bot)
	if err != nil {
		return err
	}

	if err := gormDB.WithContext(ctx).Create(dbBot).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%w: client_id already exists", repository.ErrDuplicate)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	// Reload from DB to get all fields including defaults
	result, err := r.toEntity(dbBot)
	if err != nil {
		return err
	}

	*bot = *result
	return nil
}

// Update updates an existing bot and refreshes the provided bot pointer.
func (r *GormBotRepository) Update(ctx context.Context, bot *entity.Bot) error {
	gormDB := GetTx(ctx, r.gormDB)

	dbBot, err := r.toDBModel(bot)
	if err != nil {
		return err
	}

	result := gormDB.WithContext(ctx).Model(&model.Bot{}).Where("id = ?", bot.Id).Updates(dbBot)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%w: client_id already exists", repository.ErrDuplicate)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: %v", repository.ErrNotFound, err)
	}

	// Reload from DB to get updated fields
	var updated model.Bot
	if err := gormDB.WithContext(ctx).Where("id = ?", bot.Id).First(&updated).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %v", repository.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, err)
	}

	reloaded, err := r.toEntity(&updated)
	if err != nil {
		return err
	}

	*bot = *reloaded
	return nil
}

// Delete removes a bot by ID.
func (r *GormBotRepository) Delete(ctx context.Context, id int64) error {
	gormDB := GetTx(ctx, r.gormDB)

	result := gormDB.WithContext(ctx).Delete(&model.Bot{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
			return fmt.Errorf("%w: bot has related records", repository.ErrDatabaseError)
		}
		return fmt.Errorf("%w: %v", repository.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
