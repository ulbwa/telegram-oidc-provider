package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"

	domain "github.com/ulbwa/telegram-oidc-provider/internal/domain/entities"
	"github.com/ulbwa/telegram-oidc-provider/internal/domain/repositories"
)

// userRepository implements the UserRepository interface using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if user == nil {
		return fmt.Errorf("%w: user cannot be nil", repositories.ErrInvalidArgument)
	}

	model := r.domainToModel(user)
	
	tx := GetTx(ctx, r.db)
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return mapError(err, "failed to create user")
	}

	// Update the domain entity with the created data
	if err := r.modelToDomain(&model, user); err != nil {
		return fmt.Errorf("%w: failed to convert model to domain: %v", repositories.ErrOperationFailed, err)
	}
	
	return nil
}

// Read retrieves a user by their ID
func (r *userRepository) Read(ctx context.Context, id int64, user *domain.User) error {
	if user == nil {
		return fmt.Errorf("%w: user cannot be nil", repositories.ErrInvalidArgument)
	}

	var model UserModel
	tx := GetTx(ctx, r.db)
	
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return mapError(err, fmt.Sprintf("failed to read user with id %d", id))
	}

	if err := r.modelToDomain(&model, user); err != nil {
		return fmt.Errorf("%w: failed to convert model to domain: %v", repositories.ErrOperationFailed, err)
	}
	return nil
}

// Update updates an existing user in the database
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if user == nil {
		return fmt.Errorf("%w: user cannot be nil", repositories.ErrInvalidArgument)
	}

	model := r.domainToModel(user)
	tx := GetTx(ctx, r.db)

	result := tx.WithContext(ctx).Model(&UserModel{}).Where("id = ?", model.Id).Updates(map[string]interface{}{
		"first_name": model.FirstName,
		"last_name":  model.LastName,
		"username":   model.Username,
		"photo_url":  model.PhotoUrl,
		"is_premium": model.IsPremium,
		"updated_at": model.UpdatedAt,
	})

	if result.Error != nil {
		return mapError(result.Error, fmt.Sprintf("failed to update user with id %d", model.Id))
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: user with id %d not found for update", repositories.ErrNotFound, model.Id)
	}

	return nil
}

// Delete deletes a user from the database
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	tx := GetTx(ctx, r.db)
	
	result := tx.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return mapError(result.Error, fmt.Sprintf("failed to delete user with id %d", id))
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: user with id %d not found for deletion", repositories.ErrNotFound, id)
	}
	
	return nil
}

// domainToModel converts a domain User entity to a UserModel
func (r *userRepository) domainToModel(user *domain.User) UserModel {
	model := UserModel{
		Id:        user.Id,
		FirstName: user.FirstName,
		CreatedAt: user.CreatedAt,
	}

	if user.LastName != nil {
		model.LastName = sql.NullString{String: *user.LastName, Valid: true}
	}

	if user.Username != nil {
		model.Username = sql.NullString{String: *user.Username, Valid: true}
	}

	if user.PhotoUrl != nil {
		model.FromPhotoURL(user.PhotoUrl)
	}

	if user.IsPremium != nil {
		model.IsPremium = sql.NullBool{Bool: *user.IsPremium, Valid: true}
	}

	if user.UpdatedAt != nil {
		model.UpdatedAt = sql.NullTime{Time: *user.UpdatedAt, Valid: true}
	}

	return model
}

// modelToDomain converts a UserModel to a domain User entity
func (r *userRepository) modelToDomain(model *UserModel, user *domain.User) error {
	user.Id = model.Id
	user.FirstName = model.FirstName
	user.CreatedAt = model.CreatedAt

	if model.LastName.Valid {
		user.LastName = &model.LastName.String
	} else {
		user.LastName = nil
	}

	if model.Username.Valid {
		user.Username = &model.Username.String
	} else {
		user.Username = nil
	}

	if model.PhotoUrl.Valid {
		photoUrl, err := model.ToPhotoURL()
		if err != nil {
			return fmt.Errorf("failed to parse photo URL: %w", err)
		}
		user.PhotoUrl = photoUrl
	} else {
		user.PhotoUrl = nil
	}

	if model.IsPremium.Valid {
		user.IsPremium = &model.IsPremium.Bool
	} else {
		user.IsPremium = nil
	}

	if model.UpdatedAt.Valid {
		user.UpdatedAt = &model.UpdatedAt.Time
	} else {
		user.UpdatedAt = nil
	}

	return nil
}
