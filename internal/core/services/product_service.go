package services

import (
	"context"


	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

type productService struct {
	productRepo   repositories.ProductRepository
	inventoryRepo repositories.InventoryRepository
	db            *gorm.DB
}

// NewProductService creates a new product service
func NewProductService(
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
	db *gorm.DB,
) services.ProductService {
	return &productService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		db:            db,
	}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	// Validate SKU uniqueness
	existing, err := s.productRepo.FindBySKU(ctx, product.SKU)
	if err == nil && existing != nil {
		return errors.AlreadyExists("Product", "SKU", product.SKU)
	}

	// Validate barcode uniqueness if provided
	if product.Barcode != nil && *product.Barcode != "" {
		existing, err := s.productRepo.FindByBarcode(ctx, *product.Barcode)
		if err == nil && existing != nil {
			return errors.AlreadyExists("Product", "barcode", *product.Barcode)
		}
	}

	// Validate basic fields
	if product.Name == "" {
		return errors.InvalidInput("Product name is required")
	}

	if product.SellingPrice <= 0 {
		return errors.InvalidInput("Sale price must be positive")
	}

	// Set default status if not provided
	if product.Status == "" {
		product.Status = domain.ProductStatusActive
	}

	// Generate UUID if not provided
	if product.ProductID == uuid.Nil {
		product.ProductID = uuid.New()
	}

	return s.productRepo.Create(ctx, product)
}

// GetProduct retrieves a product by ID
func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

// GetProductBySKU retrieves a product by SKU
func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	return s.productRepo.FindBySKU(ctx, sku)
}

// GetProductByBarcode retrieves a product by barcode
func (s *productService) GetProductByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	return s.productRepo.FindByBarcode(ctx, barcode)
}

// ListProducts lists products with filters
func (s *productService) ListProducts(ctx context.Context, filters repositories.ProductFilters, limit, offset int) ([]domain.Product, int64, error) {
	return s.productRepo.List(ctx, filters, limit, offset)
}

// SearchProducts searches products by query
func (s *productService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]domain.Product, int64, error) {
	filters := repositories.ProductFilters{
		Search: query,
	}
	return s.productRepo.List(ctx, filters, limit, offset)
}

// UpdateProduct updates a product
func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	// Validate product exists
	_, err := s.productRepo.FindByID(ctx, product.ProductID)
	if err != nil {
		return errors.NotFoundWithID("Product", product.ProductID.String())
	}

	// Validate SKU uniqueness if changed
	if product.SKU != "" {
		existing, err := s.productRepo.FindBySKU(ctx, product.SKU)
		if err == nil && existing != nil && existing.ProductID != product.ProductID {
			return errors.AlreadyExists("Product", "SKU", product.SKU)
		}
	}

	// Validate barcode uniqueness if changed
	if product.Barcode != nil && *product.Barcode != "" {
		existing, err := s.productRepo.FindByBarcode(ctx, *product.Barcode)
		if err == nil && existing != nil && existing.ProductID != product.ProductID {
			return errors.AlreadyExists("Product", "barcode", *product.Barcode)
		}
	}

	return s.productRepo.Update(ctx, product)
}

// UpdatePrice updates product price with audit trail
func (s *productService) UpdatePrice(ctx context.Context, productID uuid.UUID, newPrice float64, currency domain.CurrencyCode, reason string) error {
	// Validate product exists
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.NotFoundWithID("Product", productID.String())
	}

	if newPrice <= 0 {
		return errors.InvalidInput("Price must be positive")
	}


	// Update price in transaction
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update product price
		product.SellingPrice = newPrice
		if err := tx.Save(product).Error; err != nil {
			return errors.WrapError(err, "failed to update product price")
		}


		return nil
	})
}

// DeleteProduct soft deletes a product
func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	// Validate product exists
	_, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return errors.NotFoundWithID("Product", id.String())
	}

	// TODO: Check if product has any inventory or sales history
	// For now, just soft delete
	return s.productRepo.Delete(ctx, id)
}


// GetLowStockProducts retrieves products with low stock (stub)
func (s *productService) GetLowStockProducts(ctx context.Context, warehouseID *uuid.UUID) ([]domain.Product, error) {
	// TODO: Implement when inventory methods are available
	return []domain.Product{}, nil
}
