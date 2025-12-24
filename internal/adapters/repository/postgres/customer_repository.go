package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *gorm.DB) repositories.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	// Check for duplicate TaxID
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Customer{}).
		Where("tax_id = ?", customer.TaxID).Count(&count).Error; err != nil {
		return errors.WrapError(err, "failed to check tax_id uniqueness")
	}
	if count > 0 {
		return errors.AlreadyExists("Customer", "tax_id", customer.TaxID)
	}

	if err := r.db.WithContext(ctx).Create(customer).Error; err != nil {
		return errors.WrapError(err, "failed to create customer")
	}
	return nil
}

func (r *customerRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).
		Preload("Location").
		Preload("Children").
		First(&customer, "customer_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Customer", id.String())
		}
		return nil, errors.WrapError(err, "failed to find customer")
	}
	return &customer, nil
}

func (r *customerRepository) FindByTaxID(ctx context.Context, taxID string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("tax_id = ?", taxID).
		First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Customer")
		}
		return nil, errors.WrapError(err, "failed to find customer by tax_id")
	}
	return &customer, nil
}

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("email = ?", email).
		First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Customer")
		}
		return nil, errors.WrapError(err, "failed to find customer by email")
	}
	return &customer, nil
}

func (r *customerRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).
		Preload("Location").
		Where("firebase_uid = ?", firebaseUID).
		First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Customer")
		}
		return nil, errors.WrapError(err, "failed to find customer by firebase_uid")
	}
	return &customer, nil
}

func (r *customerRepository) List(ctx context.Context, filters repositories.CustomerFilters, limit, offset int) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var total int64

	query := r.buildFilterQuery(r.db.WithContext(ctx).Model(&domain.Customer{}), filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count customers")
	}

	err := query.
		Preload("Location").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&customers).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to list customers")
	}

	return customers, total, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	if err := r.db.WithContext(ctx).Save(customer).Error; err != nil {
		return errors.WrapError(err, "failed to update customer")
	}
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Customer{}, "customer_id = ?", id).Error; err != nil {
		return errors.WrapError(err, "failed to delete customer")
	}
	return nil
}

func (r *customerRepository) GetWithChildren(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).
		Preload("Location").
		Preload("Children").
		First(&customer, "customer_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Customer", id.String())
		}
		return nil, errors.WrapError(err, "failed to find customer with children")
	}
	return &customer, nil
}

func (r *customerRepository) UpdateLoyaltyPoints(ctx context.Context, customerID uuid.UUID, points int) error {
	err := r.db.WithContext(ctx).
		Model(&domain.Customer{}).
		Where("customer_id = ?", customerID).
		Update("loyalty_points", gorm.Expr("loyalty_points + ?", points)).Error

	if err != nil {
		return errors.WrapError(err, "failed to update loyalty points")
	}
	return nil
}

func (r *customerRepository) buildFilterQuery(query *gorm.DB, filters repositories.CustomerFilters) *gorm.DB {
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Type != nil {
		query = query.Where("customer_type = ?", *filters.Type)
	}

	if filters.LocationID != nil {
		query = query.Where("location_id = ?", *filters.LocationID)
	}

	if filters.Search != "" {
		query = query.Where(
			"business_name ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR tax_id ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%",
		)
	}

	return query
}
