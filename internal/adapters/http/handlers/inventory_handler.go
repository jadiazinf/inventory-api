package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/jadiazinf/inventory/internal/adapters/http/dto"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

type InventoryHandler struct {
	inventoryService services.InventoryService
}

func NewInventoryHandler(inventoryService services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// GetInventory godoc
// @Summary Get inventory for a product in a warehouse
// @Tags inventory
// @Produce json
// @Param productId path string true "Product ID"
// @Param warehouseId path string true "Warehouse ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.InventoryResponse}
// @Router /inventory/product/{productId}/warehouse/{warehouseId} [get]
func (h *InventoryHandler) GetInventory(c *fiber.Ctx) error {
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	warehouseID, err := uuid.Parse(c.Params("warehouseId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid warehouse ID", err.Error())
	}

	inventory, err := h.inventoryService.GetInventory(c.Context(), productID, warehouseID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToInventoryResponse(inventory)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetWarehouseInventory godoc
// @Summary Get all inventory for a warehouse
// @Tags inventory
// @Produce json
// @Param warehouseId path string true "Warehouse ID"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.InventoryResponse}
// @Router /inventory/warehouse/{warehouseId} [get]
func (h *InventoryHandler) GetWarehouseInventory(c *fiber.Ctx) error {
	warehouseID, err := uuid.Parse(c.Params("warehouseId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid warehouse ID", err.Error())
	}

	inventories, err := h.inventoryService.GetWarehouseInventory(c.Context(), warehouseID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	responses := make([]dto.InventoryResponse, len(inventories))
	for i, inv := range inventories {
		responses[i] = dto.ToInventoryResponse(&inv)
	}

	return dto.SendSuccess(c, fiber.StatusOK, responses, "")
}

// GetProductInventory godoc
// @Summary Get inventory for a product across all warehouses
// @Tags inventory
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.InventoryResponse}
// @Router /inventory/product/{productId} [get]
func (h *InventoryHandler) GetProductInventory(c *fiber.Ctx) error {
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	inventories, err := h.inventoryService.GetProductInventory(c.Context(), productID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	responses := make([]dto.InventoryResponse, len(inventories))
	for i, inv := range inventories {
		responses[i] = dto.ToInventoryResponse(&inv)
	}

	return dto.SendSuccess(c, fiber.StatusOK, responses, "")
}

// CheckAvailability godoc
// @Summary Check if quantity is available for a product
// @Tags inventory
// @Produce json
// @Param productId query string true "Product ID"
// @Param warehouseId query string true "Warehouse ID"
// @Param quantity query number true "Quantity to check"
// @Success 200 {object} dto.SuccessResponse{data=map[string]bool}
// @Router /inventory/check-availability [get]
func (h *InventoryHandler) CheckAvailability(c *fiber.Ctx) error {
	productID, err := uuid.Parse(c.Query("productId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	warehouseID, err := uuid.Parse(c.Query("warehouseId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid warehouse ID", err.Error())
	}

	quantityStr := c.Query("quantity")
	if quantityStr == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "Quantity is required", nil)
	}

	quantity := 0.0
	if _, err := fmt.Sscanf(quantityStr, "%f", &quantity); err != nil || quantity <= 0 {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid quantity", nil)
	}

	available, err := h.inventoryService.CheckAvailability(c.Context(), productID, warehouseID, quantity)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, map[string]bool{"available": available}, "")
}

// RegisterInboundMovement godoc
// @Summary Register an inbound inventory movement
// @Tags inventory
// @Accept json
// @Produce json
// @Param movement body dto.InboundMovementRequest true "Inbound movement data"
// @Success 201 {object} dto.SuccessResponse
// @Router /inventory/movements/inbound [post]
func (h *InventoryHandler) RegisterInboundMovement(c *fiber.Ctx) error {
	var req dto.InboundMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	userID, ok := GetUserID(c)
	if !ok {
		return dto.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	if err := h.inventoryService.RegisterInboundMovement(
		c.Context(),
		req.ProductID,
		req.WarehouseID,
		userID,
		req.Quantity,
		req.UnitCost,
		req.Currency,
		req.ReferenceType,
		req.ReferenceID,
		req.Notes,
	); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusCreated, nil, "Inbound movement registered successfully")
}

// RegisterOutboundMovement godoc
// @Summary Register an outbound inventory movement
// @Tags inventory
// @Accept json
// @Produce json
// @Param movement body dto.OutboundMovementRequest true "Outbound movement data"
// @Success 201 {object} dto.SuccessResponse
// @Router /inventory/movements/outbound [post]
func (h *InventoryHandler) RegisterOutboundMovement(c *fiber.Ctx) error {
	var req dto.OutboundMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	userID, ok := GetUserID(c)
	if !ok {
		return dto.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	if err := h.inventoryService.RegisterOutboundMovement(
		c.Context(),
		req.ProductID,
		req.WarehouseID,
		userID,
		req.Quantity,
		req.ReferenceType,
		req.ReferenceID,
		req.Notes,
	); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusCreated, nil, "Outbound movement registered successfully")
}

// RegisterAdjustment godoc
// @Summary Register an inventory adjustment
// @Tags inventory
// @Accept json
// @Produce json
// @Param adjustment body dto.AdjustmentRequest true "Adjustment data"
// @Success 201 {object} dto.SuccessResponse
// @Router /inventory/movements/adjustment [post]
func (h *InventoryHandler) RegisterAdjustment(c *fiber.Ctx) error {
	var req dto.AdjustmentRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	userID, ok := GetUserID(c)
	if !ok {
		return dto.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	if err := h.inventoryService.RegisterAdjustment(
		c.Context(),
		req.ProductID,
		req.WarehouseID,
		userID,
		req.Quantity,
		req.Notes,
	); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusCreated, nil, "Adjustment registered successfully")
}

// GetProductMovements godoc
// @Summary Get inventory movements for a product
// @Tags inventory
// @Produce json
// @Param productId path string true "Product ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.SuccessResponse{data=dto.InventoryMovementListResponse}
// @Router /inventory/movements/product/{productId} [get]
func (h *InventoryHandler) GetProductMovements(c *fiber.Ctx) error {
	productID, err := uuid.Parse(c.Params("productId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	params := dto.GetPaginationParams(c)

	movements, total, err := h.inventoryService.GetMovements(c.Context(), productID, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToInventoryMovementListResponse(movements, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetWarehouseMovements godoc
// @Summary Get inventory movements for a warehouse
// @Tags inventory
// @Produce json
// @Param warehouseId path string true "Warehouse ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.SuccessResponse{data=dto.InventoryMovementListResponse}
// @Router /inventory/movements/warehouse/{warehouseId} [get]
func (h *InventoryHandler) GetWarehouseMovements(c *fiber.Ctx) error {
	warehouseID, err := uuid.Parse(c.Params("warehouseId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid warehouse ID", err.Error())
	}

	params := dto.GetPaginationParams(c)

	movements, total, err := h.inventoryService.GetWarehouseMovements(c.Context(), warehouseID, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToInventoryMovementListResponse(movements, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}
