package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// SchoolSupplyListRepository defines the interface for school supply list data access
type SchoolSupplyListRepository interface {
	Create(ctx context.Context, list *domain.SchoolSupplyList) error
	CreateWithItems(ctx context.Context, list *domain.SchoolSupplyList, items []domain.SchoolSupplyListItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.SchoolSupplyList, error)
	FindBySchoolLevel(ctx context.Context, level domain.SchoolLevel, schoolYear string) ([]domain.SchoolSupplyList, error)
	GetActive(ctx context.Context) ([]domain.SchoolSupplyList, error)
	List(ctx context.Context, limit, offset int) ([]domain.SchoolSupplyList, int64, error)
	GetItems(ctx context.Context, listID uuid.UUID) ([]domain.SchoolSupplyListItem, error)
	Update(ctx context.Context, list *domain.SchoolSupplyList) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
}

// ReservationFilters contains filter criteria for reservation queries
type ReservationFilters struct {
	CustomerID *uuid.UUID
	StoreID    *uuid.UUID
	Status     *domain.ReservationStatus
	DateFrom   *time.Time
	DateTo     *time.Time
}

// ReservationRepository defines the interface for reservation data access
type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	CreateWithItems(ctx context.Context, reservation *domain.Reservation, items []domain.ReservationItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	FindByNumber(ctx context.Context, reservationNumber string) (*domain.Reservation, error)
	List(ctx context.Context, filters ReservationFilters, limit, offset int) ([]domain.Reservation, int64, error)
	GetItems(ctx context.Context, reservationID uuid.UUID) ([]domain.ReservationItem, error)
	Update(ctx context.Context, reservation *domain.Reservation) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReservationStatus) error
	MarkAsFulfilled(ctx context.Context, id uuid.UUID, fulfilledBy uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID) error
	GetExpired(ctx context.Context) ([]domain.Reservation, error)
	GetExpiringFor(ctx context.Context, within time.Duration) ([]domain.Reservation, error)
}

// PreOrderRepository defines the interface for pre-order data access
type PreOrderRepository interface {
	Create(ctx context.Context, preOrder *domain.PreOrder) error
	CreateWithItems(ctx context.Context, preOrder *domain.PreOrder, items []domain.PreOrderItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.PreOrder, error)
	FindByNumber(ctx context.Context, preOrderNumber string) (*domain.PreOrder, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]domain.PreOrder, error)
	List(ctx context.Context, status *domain.PreOrderStatus, limit, offset int) ([]domain.PreOrder, int64, error)
	GetItems(ctx context.Context, preOrderID uuid.UUID) ([]domain.PreOrderItem, error)
	Update(ctx context.Context, preOrder *domain.PreOrder) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PreOrderStatus) error
	MarkAsReady(ctx context.Context, id uuid.UUID) error
}

// CustomerNotificationRepository defines the interface for notification data access
type CustomerNotificationRepository interface {
	Create(ctx context.Context, notification *domain.CustomerNotification) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.CustomerNotification, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]domain.CustomerNotification, int64, error)
	GetPending(ctx context.Context) ([]domain.CustomerNotification, error)
	GetScheduled(ctx context.Context, before time.Time) ([]domain.CustomerNotification, error)
	MarkAsSent(ctx context.Context, id uuid.UUID) error
	MarkAsFailed(ctx context.Context, id uuid.UUID) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}
