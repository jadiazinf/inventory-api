package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// ProductFilters contains filter criteria for product queries
type ProductFilters struct {
	CategoryID     *uuid.UUID
	IsSchoolSupply *bool
	SchoolLevel    *domain.SchoolLevel
	Status         *domain.ProductStatus
	Search         string
	MinPrice       *float64
	MaxPrice       *float64
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	FindBySKU(ctx context.Context, sku string) (*domain.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	List(ctx context.Context, filters ProductFilters, limit, offset int) ([]domain.Product, int64, error)
	SearchByName(ctx context.Context, query string, limit, offset int) ([]domain.Product, int64, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetLowStock(ctx context.Context, warehouseID *uuid.UUID) ([]domain.Product, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	FindByName(ctx context.Context, name string) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	GetTree(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}
