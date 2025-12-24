package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

// ReservationItemRequest represents an item in a reservation
type ReservationItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  float64   `json:"quantity" validate:"required,gt=0"`
}

// CreateReservationRequest represents a request to create a reservation
type CreateReservationRequest struct {
	CustomerID     uuid.UUID                `json:"customer_id" validate:"required"`
	ChildID        *uuid.UUID               `json:"child_id,omitempty"`
	ListID         *uuid.UUID               `json:"list_id,omitempty"`
	StoreID        uuid.UUID                `json:"store_id" validate:"required"`
	Items          []ReservationItemRequest `json:"items" validate:"required,min=1"`
	DepositAmount  float64                  `json:"deposit_amount" validate:"gte=0"`
	Currency       domain.CurrencyCode      `json:"currency" validate:"required"`
	ExpirationDays int                      `json:"expiration_days" validate:"required,gt=0"`
	Notes          *string                  `json:"notes,omitempty"`
}

// FulfillReservationRequest represents a request to fulfill a reservation
type FulfillReservationRequest struct {
	PaymentMethod    domain.PaymentMethod `json:"payment_method" validate:"required"`
	PaymentReference *string              `json:"payment_reference,omitempty"`
	ExchangeRate     *float64             `json:"exchange_rate,omitempty"`
}

// ReservationItemResponse represents a reservation item in API responses
// ReservationItemResponse represents a reservation item in API responses
type ReservationItemResponse struct {
	ReservationItemID uuid.UUID `json:"reservation_item_id"`
	ReservationID     uuid.UUID `json:"reservation_id"`
	ProductID         uuid.UUID `json:"product_id"`
	Quantity          float64   `json:"quantity"`
	ReservedQuantity  float64   `json:"reserved_quantity"`
	FulfilledQuantity float64   `json:"fulfilled_quantity"`
	UnitPrice         float64   `json:"unit_price"`
	TotalAmount       float64   `json:"total_amount"`
	IsFulfilled       bool      `json:"is_fulfilled"`
}

// ReservationResponse represents a reservation in API responses
type ReservationResponse struct {
	ReservationID     uuid.UUID                `json:"reservation_id"`
	ReservationNumber string                   `json:"reservation_number"`
	CustomerID        uuid.UUID                `json:"customer_id"`
	ChildID           *uuid.UUID               `json:"child_id,omitempty"`
	ListID            *uuid.UUID               `json:"list_id,omitempty"`
	StoreID           *uuid.UUID               `json:"store_id,omitempty"`
	Status            domain.ReservationStatus `json:"status"`
	ReservationDate   time.Time                `json:"reservation_date"`
	ExpirationDate    time.Time                `json:"expiration_date"`
	PickupDate        *time.Time               `json:"pickup_date,omitempty"`
	TotalAmount       float64                  `json:"total_amount"`
	DepositAmount     float64                  `json:"deposit_amount"`
	Balance           float64                  `json:"balance"`
	Currency          domain.CurrencyCode      `json:"currency"`
	Notes             *string                  `json:"notes,omitempty"`
	ReminderSentAt    *time.Time               `json:"reminder_sent_at,omitempty"`
	FulfilledAt       *time.Time               `json:"fulfilled_at,omitempty"`
	Items             []ReservationItemResponse `json:"items,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
}

// ReservationListResponse represents paginated reservation list
type ReservationListResponse struct {
	Reservations []ReservationResponse `json:"reservations"`
	Total        int64                 `json:"total"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
}

// ToCreateReservationServiceRequest converts DTO to service request
func (r *CreateReservationRequest) ToServiceRequest(userID uuid.UUID) services.CreateReservationRequest {
	items := make([]services.ReservationItem, len(r.Items))
	for i, item := range r.Items {
		items[i] = services.ReservationItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return services.CreateReservationRequest{
		CustomerID:     r.CustomerID,
		ChildID:        r.ChildID,
		ListID:         r.ListID,
		StoreID:        r.StoreID,
		Items:          items,
		DepositAmount:  r.DepositAmount,
		Currency:       r.Currency,
		ExpirationDays: r.ExpirationDays,
		Notes:          r.Notes,
		UserID:         userID,
	}
}

// ToFulfillReservationServiceRequest converts DTO to service request
func (r *FulfillReservationRequest) ToServiceRequest(reservationID, userID uuid.UUID) services.FulfillReservationRequest {
	return services.FulfillReservationRequest{
		ReservationID:    reservationID,
		PaymentMethod:    r.PaymentMethod,
		PaymentReference: r.PaymentReference,
		ExchangeRate:     r.ExchangeRate,
		UserID:           userID,
	}
}

// ToReservationItemResponse converts domain.ReservationItem to response
// ToReservationItemResponse converts domain.ReservationItem to response
func ToReservationItemResponse(i *domain.ReservationItem) ReservationItemResponse {
	return ReservationItemResponse{
		ReservationItemID: i.ReservationItemID,
		ReservationID:     i.ReservationID,
		ProductID:         i.ProductID,
		Quantity:          i.Quantity,
		ReservedQuantity:  i.ReservedQuantity,
		FulfilledQuantity: i.FulfilledQuantity,
		UnitPrice:         i.UnitPrice,
		TotalAmount:       i.TotalAmount,
		IsFulfilled:       i.IsFulfilled,
	}
}

// ToReservationResponse converts domain.Reservation to response
func ToReservationResponse(r *domain.Reservation) ReservationResponse {
	var items []ReservationItemResponse
	if r.Items != nil {
		items = make([]ReservationItemResponse, len(r.Items))
		for i, item := range r.Items {
			items[i] = ToReservationItemResponse(&item)
		}
	}

	return ReservationResponse{
		ReservationID:     r.ReservationID,
		ReservationNumber: r.ReservationNumber,
		CustomerID:        r.CustomerID,
		ChildID:           r.ChildID,
		ListID:            r.ListID,
		StoreID:           r.StoreID,
		Status:            r.Status,
		ReservationDate:   r.ReservationDate,
		ExpirationDate:    r.ExpirationDate,
		PickupDate:        r.PickupDate,
		TotalAmount:       r.TotalAmount,
		DepositAmount:     r.DepositAmount,
		Balance:           r.Balance,
		Currency:          r.Currency,
		Notes:             r.Notes,
		ReminderSentAt:    r.ReminderSentAt,
		FulfilledAt:       r.FulfilledAt,
		Items:             items,
		CreatedAt:         r.CreatedAt,
	}
}

// ToReservationListResponse converts reservation slice to list response
func ToReservationListResponse(reservations []domain.Reservation, total int64, limit, offset int) ReservationListResponse {
	responses := make([]ReservationResponse, len(reservations))
	for i, r := range reservations {
		responses[i] = ToReservationResponse(&r)
	}
	return ReservationListResponse{
		Reservations: responses,
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}
}
