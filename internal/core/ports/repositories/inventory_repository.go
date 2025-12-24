package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// InventoryRepository defines the interface for inventory data access
type InventoryRepository interface {
	GetByProductAndWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*domain.Inventory, error)
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]domain.Inventory, error)
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Inventory, error)
	Update(ctx context.Context, inventory *domain.Inventory) error
	CreateMovement(ctx context.Context, movement *domain.InventoryMovement) error
	GetMovements(ctx context.Context, productID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error)
	GetMovementsByWarehouse(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error)
	CheckAvailability(ctx context.Context, productID, warehouseID uuid.UUID, quantity float64) (bool, error)
}

// WarehouseRepository defines the interface for warehouse data access
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *domain.Warehouse) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error)
	FindByCode(ctx context.Context, code string) (*domain.Warehouse, error)
	FindByStore(ctx context.Context, storeID uuid.UUID) ([]domain.Warehouse, error)
	List(ctx context.Context) ([]domain.Warehouse, error)
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	Delete(ctx context.Context, id uuid.UUID) error
}
