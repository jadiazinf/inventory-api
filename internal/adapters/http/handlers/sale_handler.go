package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/jadiazinf/inventory/internal/adapters/http/dto"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

type SaleHandler struct {
	saleService services.SaleService
	arService   services.AccountsReceivableService
}

func NewSaleHandler(saleService services.SaleService, arService services.AccountsReceivableService) *SaleHandler {
	return &SaleHandler{
		saleService: saleService,
		arService:   arService,
	}
}

// CreateSale godoc
// @Summary Create a new sale
// @Tags sales
// @Accept json
// @Produce json
// @Param sale body dto.CreateSaleRequest true "Sale data"
// @Success 201 {object} dto.SuccessResponse{data=dto.SaleResponse}
// @Router /sales [post]
func (h *SaleHandler) CreateSale(c *fiber.Ctx) error {
	var req dto.CreateSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	serviceReq := req.ToServiceRequest()
	sale, err := h.saleService.CreateSale(c.Context(), serviceReq)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToSaleResponse(sale)
	return dto.SendSuccess(c, fiber.StatusCreated, response, "Sale created successfully")
}

// CreateCreditSale godoc
// @Summary Create a credit sale with accounts receivable
// @Tags sales
// @Accept json
// @Produce json
// @Param sale body dto.CreateCreditSaleRequest true "Credit sale data"
// @Success 201 {object} dto.SuccessResponse{data=dto.CreditSaleResponse}
// @Router /sales/credit [post]
func (h *SaleHandler) CreateCreditSale(c *fiber.Ctx) error {
	var req dto.CreateCreditSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	serviceReq := req.CreateSaleRequest.ToServiceRequest()
	sale, ar, err := h.saleService.CreateCreditSale(c.Context(), serviceReq, req.CreditDays)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.CreditSaleResponse{
		Sale:              dto.ToSaleResponse(sale),
		AccountsReceivable: nil,
	}
	if ar != nil {
		arResp := dto.ToAccountsReceivableResponse(ar)
		response.AccountsReceivable = &arResp
	}

	return dto.SendSuccess(c, fiber.StatusCreated, response, "Credit sale created successfully")
}

// GetSale godoc
// @Summary Get a sale by ID
// @Tags sales
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.SaleResponse}
// @Router /sales/{id} [get]
func (h *SaleHandler) GetSale(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	sale, err := h.saleService.GetSale(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToSaleResponse(sale)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetSaleByInvoiceNumber godoc
// @Summary Get a sale by invoice number
// @Tags sales
// @Produce json
// @Param invoice path string true "Invoice number"
// @Success 200 {object} dto.SuccessResponse{data=dto.SaleResponse}
// @Router /sales/invoice/{invoice} [get]
func (h *SaleHandler) GetSaleByInvoiceNumber(c *fiber.Ctx) error {
	invoice := c.Params("invoice")
	if invoice == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "Invoice number is required", nil)
	}

	sale, err := h.saleService.GetSaleByInvoiceNumber(c.Context(), invoice)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToSaleResponse(sale)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// ListSales godoc
// @Summary List sales with filters and pagination
// @Tags sales
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.SuccessResponse{data=dto.SaleListResponse}
// @Router /sales [get]
func (h *SaleHandler) ListSales(c *fiber.Ctx) error {
	params := dto.GetPaginationParams(c)
	filters := repositories.SaleFilters{}

	sales, total, err := h.saleService.ListSales(c.Context(), filters, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToSaleListResponse(sales, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// CancelSale godoc
// @Summary Cancel a sale
// @Tags sales
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} dto.SuccessResponse
// @Router /sales/{id}/cancel [post]
func (h *SaleHandler) CancelSale(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	reason := c.Query("reason", "Cancelled by user")
	if err := h.saleService.CancelSale(c.Context(), id, reason); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Sale cancelled successfully")
}

// GetDailySales godoc
// @Summary Get sales for a specific date
// @Tags sales
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.SaleResponse}
// @Router /sales/daily [get]
func (h *SaleHandler) GetDailySales(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "Date is required", nil)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", err.Error())
	}

	sales, err := h.saleService.GetDailySales(c.Context(), nil, date)
	if err != nil {
		return HandleServiceError(c, err)
	}

	responses := make([]dto.SaleResponse, len(sales))
	for i, sale := range sales {
		responses[i] = dto.ToSaleResponse(&sale)
	}

	return dto.SendSuccess(c, fiber.StatusOK, responses, "")
}

// RegisterPayment godoc
// @Summary Register a payment on accounts receivable
// @Tags accounts-receivable
// @Accept json
// @Produce json
// @Param id path string true "Receivable ID"
// @Param payment body dto.PaymentRequest true "Payment data"
// @Success 200 {object} dto.SuccessResponse
// @Router /accounts-receivable/{id}/payments [post]
func (h *SaleHandler) RegisterPayment(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	userID, _ := GetUserID(c)
	if err := h.arService.RegisterPayment(c.Context(), id, req.Amount, req.Currency, req.PaymentMethod, req.Reference, req.Notes, userID); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Payment registered successfully")
}

// GetAccountsReceivable godoc
// @Summary Get accounts receivable by ID
// @Tags accounts-receivable
// @Produce json
// @Param id path string true "Receivable ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.AccountsReceivableResponse}
// @Router /accounts-receivable/{id} [get]
func (h *SaleHandler) GetAccountsReceivable(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	ar, err := h.arService.GetAccountsReceivable(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToAccountsReceivableResponse(ar)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}
