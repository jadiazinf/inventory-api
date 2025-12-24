package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jadiazinf/inventory/internal/adapters/http/dto"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

type ReservationHandler struct {
	reservationService services.ReservationService
}

func NewReservationHandler(reservationService services.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}

// CreateReservation godoc
// @Summary Create a new reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body dto.CreateReservationRequest true "Reservation data"
// @Success 201 {object} dto.SuccessResponse{data=dto.ReservationResponse}
// @Router /reservations [post]
func (h *ReservationHandler) CreateReservation(c *fiber.Ctx) error {
	var req dto.CreateReservationRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	userID, ok := GetUserID(c)
	if !ok {
		return dto.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	serviceReq := req.ToServiceRequest(userID)
	reservation, err := h.reservationService.CreateReservation(c.Context(), serviceReq)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToReservationResponse(reservation)
	return dto.SendSuccess(c, fiber.StatusCreated, response, "Reservation created successfully")
}

// GetReservation godoc
// @Summary Get a reservation by ID
// @Tags reservations
// @Produce json
// @Param id path string true "Reservation ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.ReservationResponse}
// @Router /reservations/{id} [get]
func (h *ReservationHandler) GetReservation(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	reservation, err := h.reservationService.GetReservation(c.Context(), id)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToReservationResponse(reservation)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// GetReservationByNumber godoc
// @Summary Get a reservation by reservation number
// @Tags reservations
// @Produce json
// @Param number path string true "Reservation number"
// @Success 200 {object} dto.SuccessResponse{data=dto.ReservationResponse}
// @Router /reservations/number/{number} [get]
func (h *ReservationHandler) GetReservationByNumber(c *fiber.Ctx) error {
	number := c.Params("number")
	if number == "" {
		return dto.SendError(c, fiber.StatusBadRequest, "Reservation number is required", nil)
	}

	reservation, err := h.reservationService.GetReservationByNumber(c.Context(), number)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToReservationResponse(reservation)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// ListReservations godoc
// @Summary List reservations with filters and pagination
// @Tags reservations
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param status query string false "Status filter"
// @Success 200 {object} dto.SuccessResponse{data=dto.ReservationListResponse}
// @Router /reservations [get]
func (h *ReservationHandler) ListReservations(c *fiber.Ctx) error {
	params := dto.GetPaginationParams(c)
	filters := repositories.ReservationFilters{}

	// Apply status filter if provided
	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.ReservationStatus(statusStr)
		filters.Status = &status
	}

	reservations, total, err := h.reservationService.ListReservations(c.Context(), filters, params.Limit, params.Offset)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToReservationListResponse(reservations, total, params.Limit, params.Offset)
	return dto.SendSuccess(c, fiber.StatusOK, response, "")
}

// ConfirmReservation godoc
// @Summary Confirm a reservation
// @Tags reservations
// @Produce json
// @Param id path string true "Reservation ID"
// @Success 200 {object} dto.SuccessResponse
// @Router /reservations/{id}/confirm [post]
func (h *ReservationHandler) ConfirmReservation(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.reservationService.ConfirmReservation(c.Context(), id); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Reservation confirmed successfully")
}

// FulfillReservation godoc
// @Summary Fulfill a reservation (convert to sale)
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID"
// @Param fulfillment body dto.FulfillReservationRequest true "Fulfillment data"
// @Success 200 {object} dto.SuccessResponse{data=dto.SaleResponse}
// @Router /reservations/{id}/fulfill [post]
func (h *ReservationHandler) FulfillReservation(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.FulfillReservationRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	userID, ok := GetUserID(c)
	if !ok {
		return dto.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	serviceReq := services.FulfillReservationRequest{
		ReservationID:    id,
		PaymentMethod:    req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		ExchangeRate:     req.ExchangeRate,
		UserID:           userID,
	}

	sale, err := h.reservationService.FulfillReservation(c.Context(), serviceReq)
	if err != nil {
		return HandleServiceError(c, err)
	}

	response := dto.ToSaleResponse(sale)
	return dto.SendSuccess(c, fiber.StatusOK, response, "Reservation fulfilled successfully")
}

// CancelReservation godoc
// @Summary Cancel a reservation
// @Tags reservations
// @Produce json
// @Param id path string true "Reservation ID"
// @Param reason query string false "Cancellation reason"
// @Success 200 {object} dto.SuccessResponse
// @Router /reservations/{id}/cancel [post]
func (h *ReservationHandler) CancelReservation(c *fiber.Ctx) error {
	id, err := ParseUUID(c, "id")
	if err != nil {
		return err
	}

	userID, ok := GetUserID(c)
	if !ok {
		return dto.SendError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	reason := c.Query("reason", "Cancelled by user")
	if err := h.reservationService.CancelReservation(c.Context(), id, userID, reason); err != nil {
		return HandleServiceError(c, err)
	}

	return dto.SendSuccess(c, fiber.StatusOK, nil, "Reservation cancelled successfully")
}
