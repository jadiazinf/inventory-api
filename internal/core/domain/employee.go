package domain

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents an employee
type Employee struct {
	EmployeeID      uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"employee_id"`
	NationalID      string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"national_id"`
	FirstName       string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName        string         `gorm:"type:varchar(100);not null" json:"last_name"`
	DateOfBirth     *time.Time     `gorm:"type:date" json:"date_of_birth,omitempty"`
	Gender          *GenderType    `gorm:"type:gender_type" json:"gender,omitempty"`
	LocationID      *uuid.UUID     `gorm:"type:uuid" json:"location_id,omitempty"`
	Address         *string        `gorm:"type:text" json:"address,omitempty"`
	Phone           *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email           *string        `gorm:"type:varchar(100)" json:"email,omitempty"`
	HireDate        time.Time      `gorm:"type:date;not null" json:"hire_date"`
	TerminationDate *time.Time     `gorm:"type:date" json:"termination_date,omitempty"`
	JobTitle        string         `gorm:"type:varchar(100);not null" json:"job_title"`
	Department      *string        `gorm:"type:varchar(100)" json:"department,omitempty"`
	StoreID         *uuid.UUID     `gorm:"type:uuid" json:"store_id,omitempty"`
	BaseSalary      float64        `gorm:"type:decimal(15,2);not null" json:"base_salary"`
	SalaryCurrency  CurrencyCode   `gorm:"type:currency_code;default:'VES'" json:"salary_currency"`
	Status          EmployeeStatus `gorm:"type:employee_status;default:'ACTIVE'" json:"status"`
	PhotoURL        *string        `gorm:"type:varchar(255)" json:"photo_url,omitempty"`
	UserID          *uuid.UUID     `gorm:"type:uuid" json:"user_id,omitempty"`
	BaseModelWithUser

	// Relations
	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Store    *Store    `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Employee) TableName() string {
	return "employees"
}
