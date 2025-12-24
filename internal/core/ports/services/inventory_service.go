package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// InventoryService defines the interface for inventory business logic
type InventoryService interface {
	GetInventory(ctx context.Context, productID, warehouseID uuid.UUID) (*domain.Inventory, error)
	GetWarehouseInventory(ctx context.Context, warehouseID uuid.UUID) ([]domain.Inventory, error)
	GetProductInventory(ctx context.Context, productID uuid.UUID) ([]domain.Inventory, error)
	CheckAvailability(ctx context.Context, productID, warehouseID uuid.UUID, quantity float64) (bool, error)

	// Movement operations
	RegisterInboundMovement(ctx context.Context, productID, warehouseID, userID uuid.UUID, quantity, unitCost float64, currency domain.CurrencyCode, referenceType string, referenceID *uuid.UUID, notes string) error
	RegisterOutboundMovement(ctx context.Context, productID, warehouseID, userID uuid.UUID, quantity float64, referenceType string, referenceID *uuid.UUID, notes string) error
	RegisterAdjustment(ctx context.Context, productID, warehouseID, userID uuid.UUID, quantity float64, notes string) error
	RegisterReservation(ctx context.Context, productID, warehouseID, userID uuid.UUID, quantity float64, referenceID uuid.UUID) error
	ReleaseReservation(ctx context.Context, productID, warehouseID, userID uuid.UUID, quantity float64, referenceID uuid.UUID) error

	// History and reporting
	GetMovements(ctx context.Context, productID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error)
	GetWarehouseMovements(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error)
}
