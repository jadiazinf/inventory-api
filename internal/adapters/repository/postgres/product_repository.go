package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	// Check for duplicate SKU
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Product{}).Where("sku = ?", product.SKU).Count(&count).Error; err != nil {
		return errors.WrapError(err, "failed to check SKU uniqueness")
	}
	if count > 0 {
		return errors.AlreadyExists("Product", "SKU", product.SKU)
	}

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return errors.WrapError(err, "failed to create product")
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Unit").
		Preload("Supplier").
		First(&product, "product_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Product", id.String())
		}
		return nil, errors.WrapError(err, "failed to find product")
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Unit").
		Preload("Supplier").
		Where("sku = ?", sku).
		First(&product).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Product")
		}
		return nil, errors.WrapError(err, "failed to find product by SKU")
	}
	return &product, nil
}

func (r *productRepository) FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Unit").
		Preload("Supplier").
		Where("barcode = ?", barcode).
		First(&product).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Product")
		}
		return nil, errors.WrapError(err, "failed to find product by barcode")
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, filters repositories.ProductFilters, limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.buildFilterQuery(r.db.WithContext(ctx).Model(&domain.Product{}), filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count products")
	}

	err := query.
		Preload("Category").
		Preload("Unit").
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&products).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to list products")
	}

	return products, total, nil
}

func (r *productRepository) SearchByName(ctx context.Context, query string, limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	searchQuery := r.db.WithContext(ctx).Model(&domain.Product{}).
		Where("name ILIKE ? OR sku ILIKE ?", "%"+query+"%", "%"+query+"%")

	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count search results")
	}

	err := searchQuery.
		Preload("Category").
		Preload("Unit").
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&products).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to search products")
	}

	return products, total, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return errors.WrapError(err, "failed to update product")
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Product{}, "product_id = ?", id).Error; err != nil {
		return errors.WrapError(err, "failed to delete product")
	}
	return nil
}

func (r *productRepository) GetLowStock(ctx context.Context, warehouseID *uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product

	query := r.db.WithContext(ctx).
		Joins("INNER JOIN inventory ON inventory.product_id = products.product_id").
		Where("inventory.available_quantity <= products.min_stock")

	if warehouseID != nil {
		query = query.Where("inventory.warehouse_id = ?", *warehouseID)
	}

	err := query.
		Preload("Category").
		Preload("Unit").
		Group("products.product_id").
		Find(&products).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get low stock products")
	}

	return products, nil
}

func (r *productRepository) buildFilterQuery(query *gorm.DB, filters repositories.ProductFilters) *gorm.DB {
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if filters.IsSchoolSupply != nil {
		query = query.Where("is_school_supply = ?", *filters.IsSchoolSupply)
	}

	if filters.SchoolLevel != nil {
		query = query.Where("? = ANY(grade_levels)", *filters.SchoolLevel)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	if filters.MinPrice != nil {
		query = query.Where("selling_price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("selling_price <= ?", *filters.MaxPrice)
	}

	return query
}
