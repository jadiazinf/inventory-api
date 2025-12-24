package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/jadiazinf/inventory/internal/adapters/http/dto"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
)

type CustomerHandler struct {
	customerRepo      repositories.CustomerRepository
	customerChildRepo repositories.CustomerChildRepository
}

func NewCustomerHandler(customerRepo repositories.CustomerRepository, customerChildRepo repositories.CustomerChildRepository) *CustomerHandler {
	return &CustomerHandler{
		customerRepo:      customerRepo,
		customerChildRepo: customerChildRepo,
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body dto.CustomerRequest true "Customer data"
// @Success 201 {object} dto.SuccessResponse{data=dto.CustomerResponse}
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req dto.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	customer := req.ToCustomerDomain()

	userID, ok := GetUserID(c)
	if ok {
		customer.CreatedBy = &userID
	}

	if err := h.customerRepo.Create(c.Context(), customer); err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerResponse(customer)
	return dto.SendSuccess(c, fiber.StatusCreated, response, "Customer created successfully")
}

// GetCustomer godoc
// @Summary Get a customer by ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.CustomerResponse}
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	customer, err := h.customerRepo.FindByID(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerResponse(customer)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetCustomerWithChildren godoc
// @Summary Get a customer with their children
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.CustomerResponse}
// @Router /customers/{id}/with-children [get]
func (h *CustomerHandler) GetCustomerWithChildren(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	customer, err := h.customerRepo.GetWithChildren(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerResponse(customer)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetCustomerByTaxID godoc
// @Summary Get a customer by tax ID
// @Tags customers
// @Produce json
// @Param taxId path string true "Tax ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.CustomerResponse}
// @Router /customers/tax-id/{taxId} [get]
func (h *CustomerHandler) GetCustomerByTaxID(c *fiber.Ctx) error {
	taxID := c.Params("taxId")
	if taxID == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "Tax ID is required", nil)
	}

	customer, err := h.customerRepo.FindByTaxID(c.Context(), taxID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerResponse(customer)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// ListCustomers godoc
// @Summary List customers with filters and pagination
// @Tags customers
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param search query string false "Search term"
// @Success 200 {object} dto.SuccessResponse{data=dto.CustomerListResponse}
// @Router /customers [get]
func (h *CustomerHandler) ListCustomers(c *fiber.Ctx) error {
	params := dto.GetPaginationParams(c)
	search := c.Query("search", "")

	filters := repositories.CustomerFilters{
		Search: search,
	}

	customers, total, err := h.customerRepo.List(c.Context(), filters, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerListResponse(customers, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body dto.CustomerRequest true "Customer data"
// @Success 200 {object} dto.SuccessResponse{data=dto.CustomerResponse}
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	customer := req.ToCustomerDomain()
	customer.CustomerID = id

	userID, ok := GetUserID(c)
	if ok {
		customer.UpdatedBy = &userID
	}

	if err := h.customerRepo.Update(c.Context(), customer); err != nil {
		return HandleServiceError(c, err)
	}

	// Reload updated customer
	updated, err := h.customerRepo.FindByID(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerResponse(updated)
	return dto.SendSuccess(c, fiber.StatusOK, response, "Customer updated successfully")
}

// DeleteCustomer godoc
// @Summary Delete a customer (soft delete)
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} dto.SuccessResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.customerRepo.Delete(c.Context(), id); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Customer deleted successfully")
}

// AddChild godoc
// @Summary Add a child to a customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param child body dto.CustomerChildRequest true "Child data"
// @Success 201 {object} dto.SuccessResponse{data=dto.CustomerChildResponse}
// @Router /customers/{id}/children [post]
func (h *CustomerHandler) AddChild(c *fiber.Ctx) error {
	customerID, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.CustomerChildRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	child := req.ToCustomerChildDomain(customerID)

	if err := h.customerChildRepo.Create(c.Context(), child); err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerChildResponse(child)
	return dto.SendSuccess(c, fiber.StatusCreated, response, "Child added successfully")
}

// GetChildren godoc
// @Summary Get all children of a customer
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.CustomerChildResponse}
// @Router /customers/{id}/children [get]
func (h *CustomerHandler) GetChildren(c *fiber.Ctx) error {
	customerID, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	children, err := h.customerChildRepo.FindByCustomer(c.Context(), customerID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	responses := make([]dto.CustomerChildResponse, len(children))
	for i, child := range children {
		responses[i] = dto.ToCustomerChildResponse(&child)
	}

	return dto.SendSuccess(c, fiber.StatusOK, responses, "")
}

// UpdateChild godoc
// @Summary Update a customer child
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param childId path string true "Child ID"
// @Param child body dto.CustomerChildRequest true "Child data"
// @Success 200 {object} dto.SuccessResponse{data=dto.CustomerChildResponse}
// @Router /customers/{id}/children/{childId} [put]
func (h *CustomerHandler) UpdateChild(c *fiber.Ctx) error {
	customerID, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	childID, err := uuid.Parse(c.Params("childId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid child ID", err.Error())
	}

	var req dto.CustomerChildRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	child := req.ToCustomerChildDomain(customerID)
	child.ChildID = childID

	if err := h.customerChildRepo.Update(c.Context(), child); err != nil {
		return HandleServiceError(c, err)
	}

	// Reload updated child
	updated, err := h.customerChildRepo.FindByID(c.Context(), childID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToCustomerChildResponse(updated)
	return dto.SendSuccess(c, fiber.StatusOK, response, "Child updated successfully")
}

// DeleteChild godoc
// @Summary Delete a customer child
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Param childId path string true "Child ID"
// @Success 200 {object} dto.SuccessResponse
// @Router /customers/{id}/children/{childId} [delete]
func (h *CustomerHandler) DeleteChild(c *fiber.Ctx) error {
	childID, err := uuid.Parse(c.Params("childId"))
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid child ID", err.Error())
	}

	if err := h.customerChildRepo.Delete(c.Context(), childID); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Child deleted successfully")
}

// UpdateLoyaltyPoints godoc
// @Summary Update customer loyalty points
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param points body map[string]int true "Points to add/subtract"
// @Success 200 {object} dto.SuccessResponse
// @Router /customers/{id}/loyalty-points [put]
func (h *CustomerHandler) UpdateLoyaltyPoints(c *fiber.Ctx) error {
	customerID, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	var req struct {
		Points int `json:"points"`
	}
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.customerRepo.UpdateLoyaltyPoints(c.Context(), customerID, req.Points); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Loyalty points updated successfully")
}
