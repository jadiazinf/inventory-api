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

type inventoryService struct {
	inventoryRepo repositories.InventoryRepository
	productRepo   repositories.ProductRepository
	db            *gorm.DB
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	inventoryRepo repositories.InventoryRepository,
	productRepo repositories.ProductRepository,
	db *gorm.DB,
) services.InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		db:            db,
	}
}

// GetInventory retrieves inventory for a product at a warehouse
func (s *inventoryService) GetInventory(ctx context.Context, productID, warehouseID uuid.UUID) (*domain.Inventory, error) {
	return s.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
}

// GetWarehouseInventory retrieves all inventory for a warehouse
func (s *inventoryService) GetWarehouseInventory(ctx context.Context, warehouseID uuid.UUID) ([]domain.Inventory, error) {
	return s.inventoryRepo.GetByWarehouse(ctx, warehouseID)
}

// GetProductInventory retrieves inventory for a product across all warehouses
func (s *inventoryService) GetProductInventory(ctx context.Context, productID uuid.UUID) ([]domain.Inventory, error) {
	return s.inventoryRepo.GetByProduct(ctx, productID)
}

// CheckAvailability checks if quantity is available
func (s *inventoryService) CheckAvailability(ctx context.Context, productID, warehouseID uuid.UUID, quantity float64) (bool, error) {
	return s.inventoryRepo.CheckAvailability(ctx, productID, warehouseID, quantity)
}

// RegisterInboundMovement registers an inbound inventory movement
func (s *inventoryService) RegisterInboundMovement(
	ctx context.Context,
	productID, warehouseID, userID uuid.UUID,
	quantity, unitCost float64,
	currency domain.CurrencyCode,
	referenceType string,
	referenceID *uuid.UUID,
	notes string,
) error {
	// Validate product exists
	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.NotFoundWithID("Product", productID.String())
	}

	if quantity <= 0 {
		return errors.InvalidInput("Quantity must be positive")
	}

	if unitCost <= 0 {
		return errors.InvalidInput("Unit cost must be positive")
	}

	// Create movement
	movement := &domain.InventoryMovement{
		MovementID:    uuid.New(),
		ProductID:     productID,
		WarehouseID:   warehouseID,
		MovementType:  domain.MovementTypeIn,
		Quantity:      quantity,
		Currency:      currency,
		ReferenceType: &referenceType,
		ReferenceID:   referenceID,
		Notes:         &notes,
		CreatedBy:     &userID,
	}

	return s.inventoryRepo.CreateMovement(ctx, movement)
}

// RegisterOutboundMovement registers an outbound inventory movement
func (s *inventoryService) RegisterOutboundMovement(
	ctx context.Context,
	productID, warehouseID, userID uuid.UUID,
	quantity float64,
	referenceType string,
	referenceID *uuid.UUID,
	notes string,
) error {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.NotFoundWithID("Product", productID.String())
	}

	if quantity <= 0 {
		return errors.InvalidInput("Quantity must be positive")
	}

	available, err := s.inventoryRepo.CheckAvailability(ctx, productID, warehouseID, quantity)
	if err != nil {
		return err
	}

	if !available {
		inventory, _ := s.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
		availableQty := 0.0
		if inventory != nil {
			availableQty = inventory.AvailableQuantity
		}
		return errors.InsufficientStock(product.Name, availableQty, quantity)
	}


	movement := &domain.InventoryMovement{
		MovementID:    uuid.New(),
		ProductID:     productID,
		WarehouseID:   warehouseID,
		MovementType:  domain.MovementTypeOut,
		Quantity:      quantity,
		UnitCost:      product.CostPrice,
		Currency:      domain.CurrencyVES,
		ReferenceType: &referenceType,
		ReferenceID:   referenceID,
		Notes:         &notes,
		CreatedBy:     &userID,
	}

	return s.inventoryRepo.CreateMovement(ctx, movement)
}

// RegisterAdjustment registers an inventory adjustment
func (s *inventoryService) RegisterAdjustment(
	ctx context.Context,
	productID, warehouseID, userID uuid.UUID,
	quantity float64,
	notes string,
) error {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.NotFoundWithID("Product", productID.String())
	}

	inventory, err := s.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	movementType := domain.MovementTypeAdjustment
	adjustQty := quantity
	if quantity < 0 {
		adjustQty = -quantity
		if adjustQty > inventory.AvailableQuantity {
			return errors.InsufficientStock(product.Name, inventory.AvailableQuantity, adjustQty)
		}
	}


	movement := &domain.InventoryMovement{
		MovementID:    uuid.New(),
		ProductID:     productID,
		WarehouseID:   warehouseID,
		MovementType:  movementType,
		Quantity:      adjustQty,
		UnitCost:      product.CostPrice,
		Currency:      domain.CurrencyVES,
		ReferenceType: stringPtr("ADJUSTMENT"),
		Notes:         &notes,
		CreatedBy:     &userID,
	}

	return s.inventoryRepo.CreateMovement(ctx, movement)
}

// RegisterReservation registers an inventory reservation
func (s *inventoryService) RegisterReservation(
	ctx context.Context,
	productID, warehouseID, userID uuid.UUID,
	quantity float64,
	referenceID uuid.UUID,
) error {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.NotFoundWithID("Product", productID.String())
	}

	if quantity <= 0 {
		return errors.InvalidInput("Quantity must be positive")
	}

	available, err := s.inventoryRepo.CheckAvailability(ctx, productID, warehouseID, quantity)
	if err != nil {
		return err
	}

	if !available {
		inventory, _ := s.inventoryRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
		availableQty := 0.0
		if inventory != nil {
			availableQty = inventory.AvailableQuantity
		}
		return errors.InsufficientStock(product.Name, availableQty, quantity)
	}


	movement := &domain.InventoryMovement{
		MovementID:    uuid.New(),
		ProductID:     productID,
		WarehouseID:   warehouseID,
		MovementType:  domain.MovementTypeReservation,
		Quantity:      quantity,
		UnitCost:      product.CostPrice,
		Currency:      domain.CurrencyVES,
		ReferenceType: stringPtr("RESERVATION"),
		ReferenceID:   &referenceID,
		CreatedBy:     &userID,
	}

	return s.inventoryRepo.CreateMovement(ctx, movement)
}

// ReleaseReservation releases a reservation
func (s *inventoryService) ReleaseReservation(
	ctx context.Context,
	productID, warehouseID, userID uuid.UUID,
	quantity float64,
	referenceID uuid.UUID,
) error {
	if quantity <= 0 {
		return errors.InvalidInput("Quantity must be positive")
	}

	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.NotFoundWithID("Product", productID.String())
	}


	movement := &domain.InventoryMovement{
		MovementID:    uuid.New(),
		ProductID:     productID,
		WarehouseID:   warehouseID,
		MovementType:  domain.MovementTypeReservationRelease,
		Quantity:      quantity,
		UnitCost:      product.CostPrice,
		Currency:      domain.CurrencyVES,
		ReferenceType: stringPtr("RESERVATION_RELEASE"),
		ReferenceID:   &referenceID,
		CreatedBy:     &userID,
	}

	return s.inventoryRepo.CreateMovement(ctx, movement)
}

// GetMovements retrieves movement history for a product
func (s *inventoryService) GetMovements(ctx context.Context, productID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error) {
	return s.inventoryRepo.GetMovements(ctx, productID, limit, offset)
}

// GetWarehouseMovements retrieves movement history for a warehouse
func (s *inventoryService) GetWarehouseMovements(ctx context.Context, warehouseID uuid.UUID, limit, offset int) ([]domain.InventoryMovement, int64, error) {
	return s.inventoryRepo.GetMovementsByWarehouse(ctx, warehouseID, limit, offset)
}
