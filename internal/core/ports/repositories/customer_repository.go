package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// CustomerFilters contains filter criteria for customer queries
type CustomerFilters struct {
	Status     *domain.CustomerStatus
	Type       *domain.CustomerType
	LocationID *uuid.UUID
	Search     string
}

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	FindByTaxID(ctx context.Context, taxID string) (*domain.Customer, error)
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.Customer, error)
	List(ctx context.Context, filters CustomerFilters, limit, offset int) ([]domain.Customer, int64, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithChildren(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	UpdateLoyaltyPoints(ctx context.Context, customerID uuid.UUID, points int) error
}

// CustomerChildRepository defines the interface for customer children data access
type CustomerChildRepository interface {
	Create(ctx context.Context, child *domain.CustomerChild) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.CustomerChild, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]domain.CustomerChild, error)
	Update(ctx context.Context, child *domain.CustomerChild) error
	Delete(ctx context.Context, id uuid.UUID) error
}
