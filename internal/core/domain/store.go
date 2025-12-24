package domain

import (
	"github.com/google/uuid"
)

// Store represents a physical store/branch
type Store struct {
	StoreID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"store_id"`
	Code         string    `gorm:"type:varchar(20);not null;uniqueIndex" json:"code"`
	Name         string    `gorm:"type:varchar(200);not null" json:"name"`
	LocationID   *uuid.UUID `gorm:"type:uuid" json:"location_id,omitempty"`
	Address      *string   `gorm:"type:text" json:"address,omitempty"`
	Phone        *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email        *string   `gorm:"type:varchar(100)" json:"email,omitempty"`
	ManagerID    *uuid.UUID `gorm:"type:uuid" json:"manager_id,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	OpeningHours JSONB     `gorm:"type:jsonb" json:"opening_hours,omitempty"`
	Metadata     JSONB     `gorm:"type:jsonb" json:"metadata,omitempty"`
	BaseModel

	// Relations
	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Manager  *User     `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
}

func (Store) TableName() string {
	return "stores"
}
