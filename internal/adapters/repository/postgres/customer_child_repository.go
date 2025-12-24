package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type customerChildRepository struct {
	db *gorm.DB
}

// NewCustomerChildRepository creates a new customer child repository
func NewCustomerChildRepository(db *gorm.DB) repositories.CustomerChildRepository {
	return &customerChildRepository{db: db}
}

func (r *customerChildRepository) Create(ctx context.Context, child *domain.CustomerChild) error {
	if err := r.db.WithContext(ctx).Create(child).Error; err != nil {
		return errors.WrapError(err, "failed to create customer child")
	}
	return nil
}

func (r *customerChildRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.CustomerChild, error) {
	var child domain.CustomerChild
	err := r.db.WithContext(ctx).
		First(&child, "child_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("CustomerChild", id.String())
		}
		return nil, errors.WrapError(err, "failed to find customer child")
	}
	return &child, nil
}

func (r *customerChildRepository) FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]domain.CustomerChild, error) {
	var children []domain.CustomerChild
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&children).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to find children by customer")
	}
	return children, nil
}

func (r *customerChildRepository) Update(ctx context.Context, child *domain.CustomerChild) error {
	if err := r.db.WithContext(ctx).Save(child).Error; err != nil {
		return errors.WrapError(err, "failed to update customer child")
	}
	return nil
}

func (r *customerChildRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.CustomerChild{}, "child_id = ?", id).Error; err != nil {
		return errors.WrapError(err, "failed to delete customer child")
	}
	return nil
}
