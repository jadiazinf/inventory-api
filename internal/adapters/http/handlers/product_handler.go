package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/jadiazinf/inventory/internal/adapters/http/dto"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.ProductRequest true "Product data"
// @Success 201 {object} dto.SuccessResponse{data=dto.ProductResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	product := req.ToProductDomain()

	// Get user from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if ok {
		product.CreatedBy = &userID
	}

	if err := h.productService.CreateProduct(c.Context(), product); err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToProductResponse(product)
	return dto.SendSuccess(c, fiber.StatusCreated, response, "Product created successfully")
}

// GetProduct godoc
// @Summary Get a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.ProductResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	product, err := h.productService.GetProduct(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToProductResponse(product)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetProductBySKU godoc
// @Summary Get a product by SKU
// @Tags products
// @Produce json
// @Param sku path string true "Product SKU"
// @Success 200 {object} dto.SuccessResponse{data=dto.ProductResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/sku/{sku} [get]
func (h *ProductHandler) GetProductBySKU(c *fiber.Ctx) error {
	sku := c.Params("sku")
	if sku == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "SKU is required", nil)
	}

	product, err := h.productService.GetProductBySKU(c.Context(), sku)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToProductResponse(product)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// ListProducts godoc
// @Summary List products with pagination
// @Tags products
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param search query string false "Search term"
// @Success 200 {object} dto.SuccessResponse{data=dto.ProductListResponse}
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	params := dto.GetPaginationParams(c)
	search := c.Query("search", "")

	filters := repositories.ProductFilters{
		Search: search,
	}

	products, total, err := h.productService.ListProducts(c.Context(), filters, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToProductListResponse(products, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// SearchProducts godoc
// @Summary Search products
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.SuccessResponse{data=dto.ProductListResponse}
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "Search query is required", nil)
	}

	params := dto.GetPaginationParams(c)

	products, total, err := h.productService.SearchProducts(c.Context(), query, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToProductListResponse(products, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// UpdateProduct godoc
// @Summary Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.ProductRequest true "Product data"
// @Success 200 {object} dto.SuccessResponse{data=dto.ProductResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	product := req.ToProductDomain()
	product.ProductID = id

	// Get user from context
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if ok {
		product.UpdatedBy = &userID
	}

	if err := h.productService.UpdateProduct(c.Context(), product); err != nil {
		return HandleServiceError(c, err)
	}

	// Reload updated product
	updated, err := h.productService.GetProduct(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToProductResponse(updated)
	return dto.SendSuccess(c, fiber.StatusOK, response, "Product updated successfully")
}

// UpdatePrice godoc
// @Summary Update product price
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param price body dto.UpdatePriceRequest true "Price data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id}/price [put]
func (h *ProductHandler) UpdatePrice(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	var req dto.UpdatePriceRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.productService.UpdatePrice(c.Context(), id, req.NewPrice, req.Currency, req.Reason); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Price updated successfully")
}

// DeleteProduct godoc
// @Summary Delete a product (soft delete)
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid product ID", err.Error())
	}

	if err := h.productService.DeleteProduct(c.Context(), id); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Product deleted successfully")
}
