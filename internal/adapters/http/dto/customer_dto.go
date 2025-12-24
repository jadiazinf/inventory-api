package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// CustomerRequest represents the request to create/update a customer
type CustomerRequest struct {
	FirebaseUID    *string               `json:"firebase_uid,omitempty"`
	TaxID          string                `json:"tax_id" validate:"required"`
	CustomerType   domain.CustomerType   `json:"customer_type" validate:"required"`
	BusinessName   *string               `json:"business_name,omitempty"`
	FirstName      *string               `json:"first_name,omitempty"`
	LastName       *string               `json:"last_name,omitempty"`
	Email          *string               `json:"email,omitempty" validate:"omitempty,email"`
	Phone          *string               `json:"phone,omitempty"`
	LocationID     *uuid.UUID            `json:"location_id,omitempty"`
	Address        *string               `json:"address,omitempty"`
	CreditLimit    float64               `json:"credit_limit,omitempty"`
	CreditDays     int                   `json:"credit_days,omitempty"`
	Status         domain.CustomerStatus `json:"status,omitempty"`
}

// CustomerResponse represents a customer in API responses
type CustomerResponse struct {
	CustomerID      uuid.UUID             `json:"customer_id"`
	FirebaseUID     *string               `json:"firebase_uid,omitempty"`
	TaxID           string                `json:"tax_id"`
	CustomerType    domain.CustomerType   `json:"customer_type"`
	BusinessName    *string               `json:"business_name,omitempty"`
	FirstName       *string               `json:"first_name,omitempty"`
	LastName        *string               `json:"last_name,omitempty"`
	Email           *string               `json:"email,omitempty"`
	Phone           *string               `json:"phone,omitempty"`
	LocationID      *uuid.UUID            `json:"location_id,omitempty"`
	Address         *string               `json:"address,omitempty"`
	CreditLimit     float64               `json:"credit_limit"`
	CreditDays      int                   `json:"credit_days"`
	LoyaltyPoints   int                   `json:"loyalty_points"`
	Status          domain.CustomerStatus `json:"status"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// CustomerChildRequest represents a child of a customer
type CustomerChildRequest struct {
	FirstName    string              `json:"first_name" validate:"required"`
	LastName     string              `json:"last_name" validate:"required"`
	DateOfBirth  *time.Time          `json:"date_of_birth,omitempty"`
	SchoolLevel  *domain.SchoolLevel `json:"school_level,omitempty"`
	Grade        *string             `json:"grade,omitempty"`
	SchoolName   *string             `json:"school_name,omitempty"`
	Notes        *string             `json:"notes,omitempty"`
}

// CustomerChildResponse represents a child in API responses
type CustomerChildResponse struct {
	ChildID      uuid.UUID           `json:"child_id"`
	CustomerID   uuid.UUID           `json:"customer_id"`
	FirstName    string              `json:"first_name"`
	LastName     string              `json:"last_name"`
	DateOfBirth  *time.Time          `json:"date_of_birth,omitempty"`
	SchoolLevel  *domain.SchoolLevel `json:"school_level,omitempty"`
	Grade        *string             `json:"grade,omitempty"`
	SchoolName   *string             `json:"school_name,omitempty"`
	Notes        *string             `json:"notes,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// CustomerListResponse represents paginated customer list
type CustomerListResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int64              `json:"total"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
}

// ToCustomerDomain converts CustomerRequest to domain.Customer
func (r *CustomerRequest) ToCustomerDomain() *domain.Customer {
	status := r.Status
	if status == "" {
		status = domain.CustomerStatusActive
	}

	return &domain.Customer{
		CustomerID:   uuid.New(),
		FirebaseUID:  r.FirebaseUID,
		TaxID:        r.TaxID,
		CustomerType: r.CustomerType,
		BusinessName: r.BusinessName,
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        r.Email,
		Phone:        r.Phone,
		LocationID:   r.LocationID,
		Address:      r.Address,
		CreditLimit:  r.CreditLimit,
		CreditDays:   r.CreditDays,
		LoyaltyPoints: 0,
		Status:       status,
	}
}

// ToCustomerResponse converts domain.Customer to CustomerResponse
func ToCustomerResponse(c *domain.Customer) CustomerResponse {
	return CustomerResponse{
		CustomerID:   c.CustomerID,
		FirebaseUID:  c.FirebaseUID,
		TaxID:        c.TaxID,
		CustomerType: c.CustomerType,
		BusinessName: c.BusinessName,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		Email:        c.Email,
		Phone:        c.Phone,
		LocationID:   c.LocationID,
		Address:      c.Address,
		CreditLimit:  c.CreditLimit,
		CreditDays:   c.CreditDays,
		LoyaltyPoints: c.LoyaltyPoints,
		Status:       c.Status,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

// ToCustomerChildDomain converts CustomerChildRequest to domain.CustomerChild
func (r *CustomerChildRequest) ToCustomerChildDomain(customerID uuid.UUID) *domain.CustomerChild {
	return &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customerID,
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		DateOfBirth: r.DateOfBirth,
		SchoolLevel: r.SchoolLevel,
		Grade:       r.Grade,
		SchoolName:  r.SchoolName,
		Notes:       r.Notes,
	}
}

// ToCustomerChildResponse converts domain.CustomerChild to CustomerChildResponse
func ToCustomerChildResponse(c *domain.CustomerChild) CustomerChildResponse {
	return CustomerChildResponse{
		ChildID:     c.ChildID,
		CustomerID:  c.CustomerID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		DateOfBirth: c.DateOfBirth,
		SchoolLevel: c.SchoolLevel,
		Grade:       c.Grade,
		SchoolName:  c.SchoolName,
		Notes:       c.Notes,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ToCustomerListResponse converts customer slice to list response
func ToCustomerListResponse(customers []domain.Customer, total int64, limit, offset int) CustomerListResponse {
	responses := make([]CustomerResponse, len(customers))
	for i, c := range customers {
		responses[i] = ToCustomerResponse(&c)
	}
	return CustomerListResponse{
		Customers: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}
}
