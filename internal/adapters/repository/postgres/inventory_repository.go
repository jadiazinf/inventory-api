package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB) repositories.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		First(&inventory).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Inventory")
		}
		return nil, errors.WrapError(err, "failed to get inventory")
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]domain.Inventory, error) {
	var inventories []domain.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Preload("Warehouse").
		Where("warehouse_id = ?", warehouseID).
		Find(&inventories).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get warehouse inventory")
	}
	return inventories, nil
}

func (r *inventoryRepository) GetByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Inventory, error) {
	var inventories []domain.Inventory
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Warehouse").
		Where("product_id = ?", productID).
		Find(&inventories).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get product inventory")
	}
	return inventories, nil
}

func (r *inventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	if err := r.db.WithContext(ctx).Save(inventory).Error; err != nil {
		return errors.WrapError(err, "failed to update inventory")
	}
	return nil
}

func (r *inventoryRepository) CreateMovement(ctx context.Context, movement *domain.InventoryMovement) error {
	// Create movement - trigger will automatically update inventory table
	if err := r.db.WithContext(ctx).Create(movement).Error; err != nil {
		return errors.WrapError(err, "failed to create inventory movement")
	}
	return nil
}

func (r *inventoryRepository) GetMovements(ctx context.Context, productID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error) {
	var movements []domain.InventoryMovement
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.InventoryMovement{}).
		Where("product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count movements")
	}

	err := query.
		Preload("Product").
		Preload("Warehouse").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&movements).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to get movements")
	}

	return movements, total, nil
}

func (r *inventoryRepository) GetMovementsByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error) {
	var movements []domain.InventoryMovement
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.InventoryMovement{}).
		Where("warehouse_id = ?", warehouseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count movements")
	}

	err := query.
		Preload("Product").
		Preload("Warehouse").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&movements).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to get movements")
	}

	return movements, total, nil
}

func (r *inventoryRepository) CheckAvailability(ctx context.Context, productID, warehouseID uuid.UUID, quantity float64) (bool, error) {
	var inventory domain.Inventory
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		First(&inventory).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No inventory record means 0 available
			return false, nil
		}
		return false, errors.WrapError(err, "failed to check availability")
	}

	return inventory.AvailableQuantity >= quantity, nil
}
