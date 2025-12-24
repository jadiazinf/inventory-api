package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// SaleFilters contains filter criteria for sale queries
type SaleFilters struct {
	CustomerID    *uuid.UUID
	StoreID       *uuid.UUID
	SalespersonID *uuid.UUID
	SaleType      *domain.SaleType
	Status        *domain.SaleStatus
	DateFrom      *time.Time
	DateTo        *time.Time
}

// SaleRepository defines the interface for sale data access
type SaleRepository interface {
	Create(ctx context.Context, sale *domain.Sale) error
	CreateWithDetails(ctx context.Context, sale *domain.Sale, details []domain.SaleDetail) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
	FindByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Sale, error)
	List(ctx context.Context, filters SaleFilters, limit, offset int) ([]domain.Sale, int64, error)
	GetDetails(ctx context.Context, saleID uuid.UUID) ([]domain.SaleDetail, error)
	Update(ctx context.Context, sale *domain.Sale) error
	Cancel(ctx context.Context, id uuid.UUID) error
	GetDailySales(ctx context.Context, storeID *uuid.UUID, date time.Time) ([]domain.Sale, error)
	GetSalesByPeriod(ctx context.Context, from, to time.Time) ([]domain.Sale, error)
}

// AccountsReceivableRepository defines the interface for accounts receivable data access
type AccountsReceivableRepository interface {
	Create(ctx context.Context, receivable *domain.AccountsReceivable) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.AccountsReceivable, error)
	FindBySale(ctx context.Context, saleID uuid.UUID) (*domain.AccountsReceivable, error)
	FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]domain.AccountsReceivable, error)
	GetOverdue(ctx context.Context) ([]domain.AccountsReceivable, error)
	List(ctx context.Context, status *domain.AccountStatus, limit, offset int) ([]domain.AccountsReceivable, int64, error)
	Update(ctx context.Context, receivable *domain.AccountsReceivable) error
	AddPayment(ctx context.Context, payment *domain.CustomerPayment) error
	GetPayments(ctx context.Context, receivableID uuid.UUID) ([]domain.CustomerPayment, error)
}
