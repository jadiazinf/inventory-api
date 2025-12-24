package domain

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

// Location represents a hierarchical location (Country -> State -> City -> Neighborhood)
type Location struct {
	LocationID   uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"location_id"`
	Name         string       `gorm:"type:varchar(200);not null" json:"name"`
	LocationType LocationType `gorm:"type:location_type;not null" json:"location_type"`
	ParentID     *uuid.UUID   `gorm:"type:uuid" json:"parent_id,omitempty"`
	Code         *string      `gorm:"type:varchar(20)" json:"code,omitempty"`
	FullPath     string       `gorm:"type:text" json:"full_path"`
	Latitude     *float64     `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude    *float64     `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	IsActive     bool         `gorm:"default:true" json:"is_active"`
	Metadata     JSONB        `gorm:"type:jsonb" json:"metadata,omitempty"`
	BaseModel

	// Relations
	Parent   *Location  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Location `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Location) TableName() string {
	return "locations"
}

// JSONB is a custom type for JSONB columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
