package domain

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Category represents a product category
type Category struct {
	CategoryID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"category_id"`
	Name             string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description      *string   `gorm:"type:text" json:"description,omitempty"`
	ParentCategoryID *uuid.UUID `gorm:"type:uuid" json:"parent_category_id,omitempty"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	Icon             *string   `gorm:"type:varchar(50)" json:"icon,omitempty"`
	SortOrder        *int      `json:"sort_order,omitempty"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	ParentCategory *Category  `gorm:"foreignKey:ParentCategoryID" json:"parent_category,omitempty"`
	SubCategories  []Category `gorm:"foreignKey:ParentCategoryID" json:"sub_categories,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

// UnitOfMeasure represents a unit of measure
type UnitOfMeasure struct {
	UnitID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"unit_id"`
	Code         string    `gorm:"type:varchar(10);not null;uniqueIndex" json:"code"`
	Name         string    `gorm:"type:varchar(50);not null" json:"name"`
	Abbreviation *string   `gorm:"type:varchar(10)" json:"abbreviation,omitempty"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
}

func (UnitOfMeasure) TableName() string {
	return "units_of_measure"
}

// SchoolLevelArray is a custom type for PostgreSQL array of school_level enum
type SchoolLevelArray []SchoolLevel

func (a SchoolLevelArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	strArray := make([]string, len(a))
	for i, v := range a {
		strArray[i] = string(v)
	}
	return pq.Array(strArray).Value()
}

func (a *SchoolLevelArray) Scan(value interface{}) error {
	var strArray pq.StringArray
	if err := strArray.Scan(value); err != nil {
		return err
	}
	*a = make(SchoolLevelArray, len(strArray))
	for i, v := range strArray {
		(*a)[i] = SchoolLevel(v)
	}
	return nil
}

// Product represents a product
type Product struct {
	ProductID      uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"product_id"`
	Barcode        *string          `gorm:"type:varchar(50);uniqueIndex" json:"barcode,omitempty"`
	SKU            string           `gorm:"type:varchar(50);not null;uniqueIndex" json:"sku"`
	Name           string           `gorm:"type:varchar(200);not null" json:"name"`
	Description    *string          `gorm:"type:text" json:"description,omitempty"`
	CategoryID     *uuid.UUID       `gorm:"type:uuid" json:"category_id,omitempty"`
	UnitID         *uuid.UUID       `gorm:"type:uuid" json:"unit_id,omitempty"`
	CostPrice      *float64         `gorm:"type:decimal(15,2)" json:"cost_price,omitempty"`
	SellingPrice   float64          `gorm:"type:decimal(15,2);not null" json:"selling_price"`
	PriceCurrency  CurrencyCode     `gorm:"type:currency_code;default:'VES'" json:"price_currency"`
	MinStock       int              `gorm:"default:0" json:"min_stock"`
	MaxStock       int              `gorm:"default:0" json:"max_stock"`
	HasTax         bool             `gorm:"default:true" json:"has_tax"`
	TaxPercentage  float64          `gorm:"type:decimal(5,2);default:16.00" json:"tax_percentage"`
	Status         ProductStatus    `gorm:"type:product_status;default:'ACTIVE'" json:"status"`
	ImageURL       *string          `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Weight         *float64         `gorm:"type:decimal(10,3)" json:"weight,omitempty"`
	Dimensions     *string          `gorm:"type:varchar(50)" json:"dimensions,omitempty"`
	IsSchoolSupply bool             `gorm:"default:false" json:"is_school_supply"`
	GradeLevels    SchoolLevelArray `gorm:"type:school_level[]" json:"grade_levels,omitempty"`
	SeasonalDemand bool             `gorm:"default:false" json:"seasonal_demand"`
	ReorderPoint   *int             `json:"reorder_point,omitempty"`
	SupplierID     *uuid.UUID       `gorm:"type:uuid" json:"supplier_id,omitempty"`
	BaseModelWithUser

	// Relations
	Category *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Unit     *UnitOfMeasure `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Supplier *Supplier      `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// ProductPriceHistory tracks price changes
type ProductPriceHistory struct {
	PriceHistoryID uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"price_history_id"`
	ProductID      uuid.UUID    `gorm:"type:uuid;not null" json:"product_id"`
	OldPrice       *float64     `gorm:"type:decimal(15,2)" json:"old_price,omitempty"`
	NewPrice       float64      `gorm:"type:decimal(15,2)" json:"new_price"`
	Currency       CurrencyCode `gorm:"type:currency_code;default:'VES'" json:"currency"`
	Reason         *string      `gorm:"type:text" json:"reason,omitempty"`
	EffectiveDate  time.Time    `gorm:"type:date;not null" json:"effective_date"`
	CreatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy      *uuid.UUID   `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ProductPriceHistory) TableName() string {
	return "product_price_history"
}

// Warehouse represents a warehouse
type Warehouse struct {
	WarehouseID uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"warehouse_id"`
	Code        string     `gorm:"type:varchar(20);not null;uniqueIndex" json:"code"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	StoreID     *uuid.UUID `gorm:"type:uuid" json:"store_id,omitempty"`
	LocationID  *uuid.UUID `gorm:"type:uuid" json:"location_id,omitempty"`
	Address     *string    `gorm:"type:text" json:"address,omitempty"`
	ManagerID   *uuid.UUID `gorm:"type:uuid" json:"manager_id,omitempty"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Store    *Store    `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Manager  *Employee `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

// Inventory represents product inventory in a warehouse
type Inventory struct {
	InventoryID      uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"inventory_id"`
	ProductID        uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	WarehouseID      uuid.UUID  `gorm:"type:uuid;not null" json:"warehouse_id"`
	AvailableQuantity float64   `gorm:"type:decimal(15,3);default:0" json:"available_quantity"`
	ReservedQuantity  float64   `gorm:"type:decimal(15,3);default:0" json:"reserved_quantity"`
	InTransitQuantity float64   `gorm:"type:decimal(15,3);default:0" json:"in_transit_quantity"`
	LastMovementDate  *time.Time `json:"last_movement_date,omitempty"`
	LastCountDate     *time.Time `json:"last_count_date,omitempty"`

	// Relations
	Product   *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Warehouse *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
}

func (Inventory) TableName() string {
	return "inventory"
}

// InventoryMovement represents inventory movements
type InventoryMovement struct {
	MovementID    uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"movement_id"`
	ProductID     uuid.UUID     `gorm:"type:uuid;not null" json:"product_id"`
	WarehouseID   uuid.UUID     `gorm:"type:uuid;not null" json:"warehouse_id"`
	MovementType  MovementType  `gorm:"type:movement_type;not null" json:"movement_type"`
	Quantity      float64       `gorm:"type:decimal(15,3);not null" json:"quantity"`
	UnitCost      *float64      `gorm:"type:decimal(15,2)" json:"unit_cost,omitempty"`
	Currency      CurrencyCode  `gorm:"type:currency_code;default:'VES'" json:"currency"`
	ReferenceType *string       `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID    `gorm:"type:uuid" json:"reference_id,omitempty"`
	Notes         *string       `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy     *uuid.UUID    `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relations
	Product   *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Warehouse *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
}

func (InventoryMovement) TableName() string {
	return "inventory_movements"
}

// Supplier represents a supplier
type Supplier struct {
	SupplierID       uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"supplier_id"`
	TaxID            string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"tax_id"`
	BusinessName     string         `gorm:"type:varchar(200);not null" json:"business_name"`
	TradeName        *string        `gorm:"type:varchar(200)" json:"trade_name,omitempty"`
	Email            *string        `gorm:"type:varchar(100)" json:"email,omitempty"`
	Phone            *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	LocationID       *uuid.UUID     `gorm:"type:uuid" json:"location_id,omitempty"`
	Address          *string        `gorm:"type:text" json:"address,omitempty"`
	ContactPerson    *string        `gorm:"type:varchar(100)" json:"contact_person,omitempty"`
	CreditDays       int            `gorm:"default:0" json:"credit_days"`
	Status           CustomerStatus `gorm:"type:customer_status;default:'ACTIVE'" json:"status"`
	Notes            *string        `gorm:"type:text" json:"notes,omitempty"`
	Rating           *int           `json:"rating,omitempty"`
	TotalPurchases   float64        `gorm:"type:decimal(15,2);default:0" json:"total_purchases"`
	LastPurchaseDate *time.Time     `json:"last_purchase_date,omitempty"`
	BaseModelWithUser

	// Relations
	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
