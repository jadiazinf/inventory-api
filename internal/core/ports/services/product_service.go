package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	ListProducts(ctx context.Context, filters repositories.ProductFilters, limit, offset int) ([]domain.Product, int64, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]domain.Product, int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	UpdatePrice(ctx context.Context, productID uuid.UUID, newPrice float64, currency domain.CurrencyCode, reason string) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetLowStockProducts(ctx context.Context, warehouseID *uuid.UUID) ([]domain.Product, error)
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	GetCategoryTree(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}
