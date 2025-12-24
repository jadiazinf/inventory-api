package domain

import (
	"time"

	"github.com/google/uuid"
)

// Sale represents a sale transaction
type Sale struct {
	SaleID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"sale_id"`
	InvoiceNumber   string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"invoice_number"`
	CustomerID      *uuid.UUID     `gorm:"type:uuid" json:"customer_id,omitempty"`
	StoreID         *uuid.UUID     `gorm:"type:uuid" json:"store_id,omitempty"`
	SaleDate        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"sale_date"`
	SaleType        SaleType       `gorm:"type:sale_type;default:'CASH'" json:"sale_type"`
	Status          SaleStatus     `gorm:"type:sale_status;default:'COMPLETED'" json:"status"`
	Subtotal        float64        `gorm:"type:decimal(15,2);default:0" json:"subtotal"`
	DiscountAmount  float64        `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TaxAmount       float64        `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	TotalAmount     float64        `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Currency        CurrencyCode   `gorm:"type:currency_code;default:'VES'" json:"currency"`
	ExchangeRate    *float64       `gorm:"type:decimal(15,4)" json:"exchange_rate,omitempty"`
	PaymentMethod   *PaymentMethod `gorm:"type:payment_method" json:"payment_method,omitempty"`
	PaymentReference *string       `gorm:"type:varchar(100)" json:"payment_reference,omitempty"`
	Notes           *string        `gorm:"type:text" json:"notes,omitempty"`
	WarehouseID     *uuid.UUID     `gorm:"type:uuid" json:"warehouse_id,omitempty"`
	SalespersonID   *uuid.UUID     `gorm:"type:uuid" json:"salesperson_id,omitempty"`
	ReservationID   *uuid.UUID     `gorm:"type:uuid" json:"reservation_id,omitempty"`
	PreOrderID      *uuid.UUID     `gorm:"type:uuid" json:"pre_order_id,omitempty"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy       *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relations
	Customer    *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Store       *Store        `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Warehouse   *Warehouse    `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Salesperson *User         `gorm:"foreignKey:SalespersonID" json:"salesperson,omitempty"`
	Details     []SaleDetail  `gorm:"foreignKey:SaleID" json:"details,omitempty"`
}

func (Sale) TableName() string {
	return "sales"
}

// SaleDetail represents a line item in a sale
type SaleDetail struct {
	DetailID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"detail_id"`
	SaleID         uuid.UUID `gorm:"type:uuid;not null" json:"sale_id"`
	ProductID      uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity       float64   `gorm:"type:decimal(15,3);not null" json:"quantity"`
	UnitPrice      float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	DiscountAmount float64   `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	Subtotal       float64   `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	TaxPercentage  float64   `gorm:"type:decimal(5,2);default:0" json:"tax_percentage"`
	TaxAmount      float64   `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	Total          float64   `gorm:"type:decimal(15,2);not null" json:"total"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Sale    *Sale    `gorm:"foreignKey:SaleID" json:"sale,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (SaleDetail) TableName() string {
	return "sale_details"
}

// AccountsReceivable represents money owed by customers
type AccountsReceivable struct {
	ReceivableID uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"receivable_id"`
	SaleID       *uuid.UUID    `gorm:"type:uuid" json:"sale_id,omitempty"`
	CustomerID   uuid.UUID     `gorm:"type:uuid;not null" json:"customer_id"`
	TotalAmount  float64       `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	PaidAmount   float64       `gorm:"type:decimal(15,2);default:0" json:"paid_amount"`
	Balance      float64       `gorm:"type:decimal(15,2);not null" json:"balance"`
	Currency     CurrencyCode  `gorm:"type:currency_code;default:'VES'" json:"currency"`
	DueDate      time.Time     `gorm:"type:date;not null" json:"due_date"`
	Status       AccountStatus `gorm:"type:account_status;default:'PENDING'" json:"status"`
	Notes        *string       `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Sale     *Sale     `gorm:"foreignKey:SaleID" json:"sale,omitempty"`
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (AccountsReceivable) TableName() string {
	return "accounts_receivable"
}

// CustomerPayment represents a payment made by a customer
type CustomerPayment struct {
	PaymentID    uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"payment_id"`
	ReceivableID uuid.UUID     `gorm:"type:uuid;not null" json:"receivable_id"`
	PaymentDate  time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"payment_date"`
	Amount       float64       `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency     CurrencyCode  `gorm:"type:currency_code;default:'VES'" json:"currency"`
	PaymentMethod PaymentMethod `gorm:"type:payment_method;not null" json:"payment_method"`
	Reference    *string       `gorm:"type:varchar(100)" json:"reference,omitempty"`
	Notes        *string       `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy    *uuid.UUID    `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relations
	Receivable *AccountsReceivable `gorm:"foreignKey:ReceivableID" json:"receivable,omitempty"`
}

func (CustomerPayment) TableName() string {
	return "customer_payments"
}
