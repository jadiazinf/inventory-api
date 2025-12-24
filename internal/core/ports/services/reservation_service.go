package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
)

// ReservationItem represents an item in a reservation request
type ReservationItem struct {
	ProductID uuid.UUID
	Quantity  float64
}

// CreateReservationRequest represents a request to create a reservation
type CreateReservationRequest struct {
	CustomerID     uuid.UUID
	ChildID        *uuid.UUID
	ListID         *uuid.UUID
	StoreID        uuid.UUID
	Items          []ReservationItem
	DepositAmount  float64
	Currency       domain.CurrencyCode
	ExpirationDays int // Days until reservation expires
	Notes          *string
	UserID         uuid.UUID
}

// FulfillReservationRequest represents a request to fulfill a reservation
type FulfillReservationRequest struct {
	ReservationID    uuid.UUID
	PaymentMethod    domain.PaymentMethod
	PaymentReference *string
	ExchangeRate     *float64
	UserID           uuid.UUID
}

// ReservationService defines the interface for reservation business logic
type ReservationService interface {
	CreateReservation(ctx context.Context, req CreateReservationRequest) (*domain.Reservation, error)
	GetReservation(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	GetReservationByNumber(ctx context.Context, reservationNumber string) (*domain.Reservation, error)
	ListReservations(ctx context.Context, filters repositories.ReservationFilters, limit, offset int) ([]domain.Reservation, int64, error)

	// Workflow operations
	ConfirmReservation(ctx context.Context, id uuid.UUID) error
	FulfillReservation(ctx context.Context, req FulfillReservationRequest) (*domain.Sale, error)
	CancelReservation(ctx context.Context, id, userID uuid.UUID, reason string) error

	// Maintenance operations
	ExpireReservations(ctx context.Context) (int, error)
	SendReminders(ctx context.Context, hoursBeforeExpiration int) (int, error)
}

// SchoolSupplyListService defines the interface for school supply list business logic
type SchoolSupplyListService interface {
	CreateList(ctx context.Context, list *domain.SchoolSupplyList, items []domain.SchoolSupplyListItem) error
	GetList(ctx context.Context, id uuid.UUID) (*domain.SchoolSupplyList, error)
	GetActiveListsByLevel(ctx context.Context, level domain.SchoolLevel, schoolYear string) ([]domain.SchoolSupplyList, error)
	PublishList(ctx context.Context, id uuid.UUID) error
	ArchiveList(ctx context.Context, id uuid.UUID) error
	CheckAvailability(ctx context.Context, listID, warehouseID uuid.UUID) (map[uuid.UUID]float64, error) // Returns productID -> available quantity
}

// PreOrderService defines the interface for pre-order business logic
type PreOrderService interface {
	CreatePreOrder(ctx context.Context, customerID, storeID uuid.UUID, items []ReservationItem, depositAmount float64, currency domain.CurrencyCode, userID uuid.UUID) (*domain.PreOrder, error)
	GetPreOrder(ctx context.Context, id uuid.UUID) (*domain.PreOrder, error)
	ListPreOrders(ctx context.Context, status *domain.PreOrderStatus, limit, offset int) ([]domain.PreOrder, int64, error)
	ConfirmPreOrder(ctx context.Context, id, userID uuid.UUID) error
	MarkAsReady(ctx context.Context, id uuid.UUID) error
	SendReadyNotification(ctx context.Context, id uuid.UUID) error
	DeliverPreOrder(ctx context.Context, id uuid.UUID, paymentMethod domain.PaymentMethod, paymentReference *string, userID uuid.UUID) (*domain.Sale, error)
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	SendReservationConfirmation(ctx context.Context, reservationID uuid.UUID) error
	SendReservationReminder(ctx context.Context, reservationID uuid.UUID) error
	SendPreOrderReady(ctx context.Context, preOrderID uuid.UUID) error
	SendCustomNotification(ctx context.Context, customerID uuid.UUID, notificationType domain.NotificationType, subject, message string) error
	ProcessPendingNotifications(ctx context.Context) (int, error)
}
