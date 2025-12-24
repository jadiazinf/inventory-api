package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return errors.WrapError(err, "failed to create user")
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Location").
		First(&user, "user_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("User", id.String())
		}
		return nil, errors.WrapError(err, "failed to find user")
	}
	return &user, nil
}

func (r *userRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Role.Permissions").
		Preload("Location").
		Where("firebase_uid = ?", firebaseUID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User")
		}
		return nil, errors.WrapError(err, "failed to find user by firebase uid")
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Location").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User")
		}
		return nil, errors.WrapError(err, "failed to find user by email")
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count users")
	}

	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Location").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to list users")
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return errors.WrapError(err, "failed to update user")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.User{}, "user_id = ?", id).Error; err != nil {
		return errors.WrapError(err, "failed to delete user")
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("user_id = ?", userID).
		Update("last_login", now).Error

	if err != nil {
		return errors.WrapError(err, "failed to update last login")
	}
	return nil
}
