package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

// SaleItemRequest represents an item in a sale
type SaleItemRequest struct {
	ProductID      uuid.UUID `json:"product_id" validate:"required"`
	Quantity       float64   `json:"quantity" validate:"required,gt=0"`
	UnitPrice      *float64  `json:"unit_price,omitempty"`
	DiscountAmount *float64  `json:"discount_amount,omitempty"`
}

// CreateSaleRequest represents a request to create a sale
type CreateSaleRequest struct {
	CustomerID       *uuid.UUID           `json:"customer_id,omitempty"`
	StoreID          uuid.UUID            `json:"store_id" validate:"required"`
	WarehouseID      uuid.UUID            `json:"warehouse_id" validate:"required"`
	SaleType         domain.SaleType      `json:"sale_type" validate:"required"`
	Currency         domain.CurrencyCode  `json:"currency" validate:"required"`
	ExchangeRate     *float64             `json:"exchange_rate,omitempty"`
	DiscountAmount   *float64             `json:"discount_amount,omitempty"`
	PaymentMethod    domain.PaymentMethod `json:"payment_method" validate:"required"`
	PaymentReference *string              `json:"payment_reference,omitempty"`
	Notes            *string              `json:"notes,omitempty"`
	SalespersonID    uuid.UUID            `json:"salesperson_id" validate:"required"`
	Items            []SaleItemRequest    `json:"items" validate:"required,min=1"`
}

// CreateCreditSaleRequest represents a request to create a credit sale
type CreateCreditSaleRequest struct {
	CreateSaleRequest
	CreditDays int `json:"credit_days" validate:"required,gt=0"`
}

// SaleDetailResponse represents a sale detail in API responses
type SaleDetailResponse struct {
	DetailID       uuid.UUID `json:"sale_detail_id"`
	SaleID         uuid.UUID `json:"sale_id"`
	ProductID      uuid.UUID `json:"product_id"`
	Quantity       float64   `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	DiscountAmount float64   `json:"discount_amount"`
	Subtotal       float64   `json:"subtotal"`
	TaxAmount      float64   `json:"tax_amount"`
	Total          float64   `json:"total"`
}

// SaleResponse represents a sale in API responses
type SaleResponse struct {
	SaleID           uuid.UUID               `json:"sale_id"`
	InvoiceNumber    string                  `json:"invoice_number"`
	CustomerID       *uuid.UUID              `json:"customer_id,omitempty"`
	StoreID          *uuid.UUID              `json:"store_id,omitempty"`
	WarehouseID      *uuid.UUID              `json:"warehouse_id,omitempty"`
	SaleType         domain.SaleType         `json:"sale_type"`
	Status           domain.SaleStatus       `json:"status"`
	Currency         domain.CurrencyCode     `json:"currency"`
	ExchangeRate     *float64                `json:"exchange_rate,omitempty"`
	Subtotal         float64                 `json:"subtotal"`
	DiscountAmount   float64                 `json:"discount_amount"`
	TaxAmount        float64                 `json:"tax_amount"`
	TotalAmount      float64                 `json:"total_amount"`
	PaymentMethod    *domain.PaymentMethod   `json:"payment_method,omitempty"`
	PaymentReference *string                 `json:"payment_reference,omitempty"`
	Notes            *string                 `json:"notes,omitempty"`
	SalespersonID    *uuid.UUID              `json:"salesperson_id,omitempty"`
	Details          []SaleDetailResponse    `json:"details,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
}

// SaleListResponse represents paginated sale list
type SaleListResponse struct {
	Sales  []SaleResponse `json:"sales"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// CreditSaleResponse represents a credit sale with AR info
type CreditSaleResponse struct {
	Sale              SaleResponse               `json:"sale"`
	AccountsReceivable *AccountsReceivableResponse `json:"accounts_receivable,omitempty"`
}

// AccountsReceivableResponse represents an AR in API responses
type AccountsReceivableResponse struct {
	ReceivableID uuid.UUID           `json:"receivable_id"`
	SaleID       *uuid.UUID          `json:"sale_id,omitempty"`
	CustomerID   uuid.UUID           `json:"customer_id"`
	TotalAmount  float64             `json:"total_amount"`
	PaidAmount   float64             `json:"paid_amount"`
	Balance      float64             `json:"balance"`
	Currency     domain.CurrencyCode `json:"currency"`
	DueDate      time.Time           `json:"due_date"`
	Status       domain.AccountStatus `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
}

// PaymentRequest represents a payment on AR
type PaymentRequest struct {
	Amount        float64              `json:"amount" validate:"required,gt=0"`
	Currency      domain.CurrencyCode  `json:"currency" validate:"required"`
	PaymentMethod domain.PaymentMethod `json:"payment_method" validate:"required"`
	Reference     *string              `json:"reference,omitempty"`
	Notes         *string              `json:"notes,omitempty"`
}

// ToCreateSaleServiceRequest converts DTO to service request
func (r *CreateSaleRequest) ToServiceRequest() services.CreateSaleRequest {
	items := make([]services.SaleItem, len(r.Items))
	for i, item := range r.Items {
		discountAmt := 0.0
		if item.DiscountAmount != nil {
			discountAmt = *item.DiscountAmount
		}
		items[i] = services.SaleItem{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			DiscountAmount: discountAmt,
		}
	}

	discountAmt := 0.0
	if r.DiscountAmount != nil {
		discountAmt = *r.DiscountAmount
	}

	return services.CreateSaleRequest{
		CustomerID:       r.CustomerID,
		StoreID:          r.StoreID,
		WarehouseID:      r.WarehouseID,
		SaleType:         r.SaleType,
		Currency:         r.Currency,
		ExchangeRate:     r.ExchangeRate,
		DiscountAmount:   discountAmt,
		PaymentMethod:    &r.PaymentMethod,
		PaymentReference: r.PaymentReference,
		Notes:            r.Notes,
		SalespersonID:    r.SalespersonID,
		Items:            items,
	}
}

// ToSaleDetailResponse converts domain.SaleDetail to response
func ToSaleDetailResponse(d *domain.SaleDetail) SaleDetailResponse {
	return SaleDetailResponse{
		DetailID:       d.DetailID,
		SaleID:         d.SaleID,
		ProductID:      d.ProductID,
		Quantity:       d.Quantity,
		UnitPrice:      d.UnitPrice,
		DiscountAmount: d.DiscountAmount,
		Subtotal:       d.Subtotal,
		TaxAmount:      d.TaxAmount,
		Total:          d.Total,
	}
}

// ToSaleResponse converts domain.Sale to response
func ToSaleResponse(s *domain.Sale) SaleResponse {
	var details []SaleDetailResponse
	if s.Details != nil {
		details = make([]SaleDetailResponse, len(s.Details))
		for i, d := range s.Details {
			details[i] = ToSaleDetailResponse(&d)
		}
	}

	return SaleResponse{
		SaleID:           s.SaleID,
		InvoiceNumber:    s.InvoiceNumber,
		CustomerID:       s.CustomerID,
		StoreID:          s.StoreID,
		WarehouseID:      s.WarehouseID,
		SaleType:         s.SaleType,
		Status:           s.Status,
		Currency:         s.Currency,
		ExchangeRate:     s.ExchangeRate,
		Subtotal:         s.Subtotal,
		DiscountAmount:   s.DiscountAmount,
		TaxAmount:        s.TaxAmount,
		TotalAmount:      s.TotalAmount,
		PaymentMethod:    s.PaymentMethod,
		PaymentReference: s.PaymentReference,
		Notes:            s.Notes,
		SalespersonID:    s.SalespersonID,
		Details:          details,
		CreatedAt:        s.CreatedAt,
	}
}

// ToAccountsReceivableResponse converts domain.AccountsReceivable to response
func ToAccountsReceivableResponse(ar *domain.AccountsReceivable) AccountsReceivableResponse {
	return AccountsReceivableResponse{
		ReceivableID: ar.ReceivableID,
		SaleID:       ar.SaleID,
		CustomerID:   ar.CustomerID,
		TotalAmount:  ar.TotalAmount,
		PaidAmount:   ar.PaidAmount,
		Balance:      ar.Balance,
		Currency:     ar.Currency,
		DueDate:      ar.DueDate,
		Status:       ar.Status,
		CreatedAt:    ar.CreatedAt,
	}
}

// ToSaleListResponse converts sale slice to list response
func ToSaleListResponse(sales []domain.Sale, total int64, limit, offset int) SaleListResponse {
	responses := make([]SaleResponse, len(sales))
	for i, s := range sales {
		responses[i] = ToSaleResponse(&s)
	}
	return SaleListResponse{
		Sales:  responses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}
