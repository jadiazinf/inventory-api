package domain

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer
type Customer struct {
	CustomerID            uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"customer_id"`
	CustomerType          CustomerType       `gorm:"type:customer_type;default:'INDIVIDUAL'" json:"customer_type"`
	TaxID                 string             `gorm:"type:varchar(20);not null;uniqueIndex" json:"tax_id"`
	BusinessName          *string            `gorm:"type:varchar(200)" json:"business_name,omitempty"`
	FirstName             *string            `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName              *string            `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	Email                 *string            `gorm:"type:varchar(100)" json:"email,omitempty"`
	Phone                 *string            `gorm:"type:varchar(20)" json:"phone,omitempty"`
	LocationID            *uuid.UUID         `gorm:"type:uuid" json:"location_id,omitempty"`
	Address               *string            `gorm:"type:text" json:"address,omitempty"`
	CreditLimit           float64            `gorm:"type:decimal(15,2);default:0" json:"credit_limit"`
	CreditDays            int                `gorm:"default:0" json:"credit_days"`
	Status                CustomerStatus     `gorm:"type:customer_status;default:'ACTIVE'" json:"status"`
	Notes                 *string            `gorm:"type:text" json:"notes,omitempty"`
	LoyaltyPoints         int                `gorm:"default:0" json:"loyalty_points"`
	TotalPurchases        float64            `gorm:"type:decimal(15,2);default:0" json:"total_purchases"`
	LastPurchaseDate      *time.Time         `json:"last_purchase_date,omitempty"`
	PreferredContactMethod *NotificationType `gorm:"type:notification_type" json:"preferred_contact_method,omitempty"`
	FirebaseUID           *string            `gorm:"type:varchar(128)" json:"firebase_uid,omitempty"`
	BaseModelWithUser

	// Relations
	Location *Location        `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Children []CustomerChild  `gorm:"foreignKey:CustomerID" json:"children,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

// CustomerChild represents a customer's child for back-to-school module
type CustomerChild struct {
	ChildID     uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"child_id"`
	CustomerID  uuid.UUID    `gorm:"type:uuid;not null" json:"customer_id"`
	FirstName   string       `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName    string       `gorm:"type:varchar(100);not null" json:"last_name"`
	DateOfBirth *time.Time   `gorm:"type:date" json:"date_of_birth,omitempty"`
	SchoolLevel *SchoolLevel `gorm:"type:school_level" json:"school_level,omitempty"`
	Grade       *string      `gorm:"type:varchar(20)" json:"grade,omitempty"`
	SchoolName  *string      `gorm:"type:varchar(200)" json:"school_name,omitempty"`
	Notes       *string      `gorm:"type:text" json:"notes,omitempty"`
	BaseModel

	// Relations
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (CustomerChild) TableName() string {
	return "customer_children"
}
