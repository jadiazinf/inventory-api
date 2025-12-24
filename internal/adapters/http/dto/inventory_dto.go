package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// InventoryMovementRequest represents a request to create an inventory movement
type InventoryMovementRequest struct {
	ProductID     uuid.UUID              `json:"product_id" validate:"required"`
	WarehouseID   uuid.UUID              `json:"warehouse_id" validate:"required"`
	MovementType  domain.MovementType    `json:"movement_type" validate:"required"`
	Quantity      float64                `json:"quantity" validate:"required,gt=0"`
	UnitCost      *float64               `json:"unit_cost,omitempty"`
	ReferenceType *string                `json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID             `json:"reference_id,omitempty"`
	Notes         string                 `json:"notes,omitempty"`
}

// InventoryResponse represents inventory in API responses
type InventoryResponse struct {
	InventoryID       uuid.UUID `json:"inventory_id"`
	ProductID         uuid.UUID `json:"product_id"`
	WarehouseID       uuid.UUID `json:"warehouse_id"`
	AvailableQuantity float64   `json:"available_quantity"`
	ReservedQuantity  float64   `json:"reserved_quantity"`
}

// InventoryMovementResponse represents an inventory movement in API responses
type InventoryMovementResponse struct {
	MovementID    uuid.UUID           `json:"movement_id"`
	ProductID     uuid.UUID           `json:"product_id"`
	WarehouseID   uuid.UUID           `json:"warehouse_id"`
	MovementType  domain.MovementType `json:"movement_type"`
	Quantity      float64             `json:"quantity"`
	UnitCost      *float64            `json:"unit_cost,omitempty"`
	ReferenceType *string             `json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID          `json:"reference_id,omitempty"`
	Notes         *string             `json:"notes,omitempty"`
	CreatedBy     *uuid.UUID          `json:"created_by,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}

// InventoryListResponse represents paginated inventory list
type InventoryListResponse struct {
	Inventory []InventoryResponse `json:"inventory"`
	Total     int64               `json:"total"`
	Limit     int                 `json:"limit"`
	Offset    int                 `json:"offset"`
}

// MovementListResponse represents paginated movement list
type MovementListResponse struct {
	Movements []InventoryMovementResponse `json:"movements"`
	Total     int64                       `json:"total"`
	Limit     int                         `json:"limit"`
	Offset    int                         `json:"offset"`
}

// ToInventoryResponse converts domain.Inventory to response
func ToInventoryResponse(i *domain.Inventory) InventoryResponse {
	return InventoryResponse{
		InventoryID:       i.InventoryID,
		ProductID:         i.ProductID,
		WarehouseID:       i.WarehouseID,
		AvailableQuantity: i.AvailableQuantity,
		ReservedQuantity:  i.ReservedQuantity,
	}
}

// ToInventoryMovementResponse converts domain.InventoryMovement to response
func ToInventoryMovementResponse(m *domain.InventoryMovement) InventoryMovementResponse {
	return InventoryMovementResponse{
		MovementID:    m.MovementID,
		ProductID:     m.ProductID,
		WarehouseID:   m.WarehouseID,
		MovementType:  m.MovementType,
		Quantity:      m.Quantity,
		UnitCost:      m.UnitCost,
		ReferenceType: m.ReferenceType,
		ReferenceID:   m.ReferenceID,
		Notes:         m.Notes,
		CreatedBy:     m.CreatedBy,
		CreatedAt:     m.CreatedAt,
	}
}

// ToInventoryListResponse converts inventory slice to list response
func ToInventoryListResponse(inventory []domain.Inventory, total int64, limit, offset int) InventoryListResponse {
	responses := make([]InventoryResponse, len(inventory))
	for i, inv := range inventory {
		responses[i] = ToInventoryResponse(&inv)
	}
	return InventoryListResponse{
		Inventory: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}
}

// ToMovementListResponse converts movement slice to list response
func ToMovementListResponse(movements []domain.InventoryMovement, total int64, limit, offset int) MovementListResponse {
	responses := make([]InventoryMovementResponse, len(movements))
	for i, m := range movements {
		responses[i] = ToInventoryMovementResponse(&m)
	}
	return MovementListResponse{
		Movements: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}
}

// InventoryMovementListResponse is an alias for MovementListResponse
type InventoryMovementListResponse = MovementListResponse

// ToInventoryMovementListResponse is an alias for ToMovementListResponse
func ToInventoryMovementListResponse(movements []domain.InventoryMovement, total int64, limit, offset int) InventoryMovementListResponse {
	return ToMovementListResponse(movements, total, limit, offset)
}

// InboundMovementRequest represents a request to register an inbound movement
type InboundMovementRequest struct {
	ProductID     uuid.UUID           `json:"product_id" validate:"required"`
	WarehouseID   uuid.UUID           `json:"warehouse_id" validate:"required"`
	Quantity      float64             `json:"quantity" validate:"required,gt=0"`
	UnitCost      float64             `json:"unit_cost" validate:"required,gte=0"`
	Currency      domain.CurrencyCode `json:"currency" validate:"required"`
	ReferenceType string              `json:"reference_type" validate:"required"`
	ReferenceID   *uuid.UUID          `json:"reference_id,omitempty"`
	Notes         string              `json:"notes"`
}

// OutboundMovementRequest represents a request to register an outbound movement
type OutboundMovementRequest struct {
	ProductID     uuid.UUID  `json:"product_id" validate:"required"`
	WarehouseID   uuid.UUID  `json:"warehouse_id" validate:"required"`
	Quantity      float64    `json:"quantity" validate:"required,gt=0"`
	ReferenceType string     `json:"reference_type" validate:"required"`
	ReferenceID   *uuid.UUID `json:"reference_id,omitempty"`
	Notes         string     `json:"notes"`
}

// AdjustmentRequest represents a request to register an inventory adjustment
type AdjustmentRequest struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	WarehouseID uuid.UUID `json:"warehouse_id" validate:"required"`
	Quantity    float64   `json:"quantity" validate:"required"`
	Notes       string    `json:"notes" validate:"required"`
}
