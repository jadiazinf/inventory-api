package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
)

// SaleItem represents an item in a sale request
type SaleItem struct {
	ProductID      uuid.UUID
	Quantity       float64
	UnitPrice      *float64 // Optional, will use product price if not provided
	DiscountAmount float64
}

// CreateSaleRequest represents a request to create a sale
type CreateSaleRequest struct {
	CustomerID       *uuid.UUID
	StoreID          uuid.UUID
	WarehouseID      uuid.UUID
	SaleType         domain.SaleType
	Items            []SaleItem
	DiscountAmount   float64
	Currency         domain.CurrencyCode
	ExchangeRate     *float64
	PaymentMethod    *domain.PaymentMethod
	PaymentReference *string
	Notes            *string
	SalespersonID    uuid.UUID
}

// SaleService defines the interface for sale business logic
type SaleService interface {
	CreateSale(ctx context.Context, req CreateSaleRequest) (*domain.Sale, error)
	GetSale(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
	GetSaleByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Sale, error)
	ListSales(ctx context.Context, filters repositories.SaleFilters, limit, offset int) ([]domain.Sale, int64, error)
	CancelSale(ctx context.Context, id uuid.UUID, reason string) error
	GetDailySales(ctx context.Context, storeID *uuid.UUID, date time.Time) ([]domain.Sale, error)
	GetSalesByPeriod(ctx context.Context, from, to time.Time) ([]domain.Sale, error)

	// Credit sales
	CreateCreditSale(ctx context.Context, req CreateSaleRequest, creditDays int) (*domain.Sale, *domain.AccountsReceivable, error)
}

// AccountsReceivableService defines the interface for accounts receivable business logic
type AccountsReceivableService interface {
	GetAccountsReceivable(ctx context.Context, id uuid.UUID) (*domain.AccountsReceivable, error)
	GetCustomerReceivables(ctx context.Context, customerID uuid.UUID) ([]domain.AccountsReceivable, error)
	GetOverdueReceivables(ctx context.Context) ([]domain.AccountsReceivable, error)
	RegisterPayment(ctx context.Context, receivableID uuid.UUID, amount float64, currency domain.CurrencyCode, paymentMethod domain.PaymentMethod, reference, notes *string, userID uuid.UUID) error
	GetPaymentHistory(ctx context.Context, receivableID uuid.UUID) ([]domain.CustomerPayment, error)
}
