package domain

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UUIDArray is a custom type for PostgreSQL array of UUIDs
type UUIDArray []uuid.UUID

func (a UUIDArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	strArray := make([]string, len(a))
	for i, v := range a {
		strArray[i] = v.String()
	}
	return pq.Array(strArray).Value()
}

func (a *UUIDArray) Scan(value interface{}) error {
	var strArray pq.StringArray
	if err := strArray.Scan(value); err != nil {
		return err
	}
	*a = make(UUIDArray, len(strArray))
	for i, v := range strArray {
		id, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		(*a)[i] = id
	}
	return nil
}

// SchoolSupplyList represents a list of school supplies for a specific grade
type SchoolSupplyList struct {
	ListID             uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"list_id"`
	ListName           string           `gorm:"type:varchar(200);not null" json:"list_name"`
	SchoolLevel        SchoolLevel      `gorm:"type:school_level;not null" json:"school_level"`
	Grade              *string          `gorm:"type:varchar(20)" json:"grade,omitempty"`
	SchoolYear         string           `gorm:"type:varchar(9);not null" json:"school_year"`
	Status             SchoolListStatus `gorm:"type:school_list_status;default:'DRAFT'" json:"status"`
	Description        *string          `gorm:"type:text" json:"description,omitempty"`
	PublishDate        *time.Time       `gorm:"type:date" json:"publish_date,omitempty"`
	ExpirationDate     *time.Time       `gorm:"type:date" json:"expiration_date,omitempty"`
	TotalEstimatedCost *float64         `gorm:"type:decimal(15,2)" json:"total_estimated_cost,omitempty"`
	IsTemplate         bool             `gorm:"default:false" json:"is_template"`
	BaseModelWithUser

	// Relations
	Items []SchoolSupplyListItem `gorm:"foreignKey:ListID" json:"items,omitempty"`
}

func (SchoolSupplyList) TableName() string {
	return "school_supply_lists"
}

// SchoolSupplyListItem represents an item in a school supply list
type SchoolSupplyListItem struct {
	ListItemID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"list_item_id"`
	ListID             uuid.UUID `gorm:"type:uuid;not null" json:"list_id"`
	ProductID          uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity           int       `gorm:"not null" json:"quantity"`
	IsRequired         bool      `gorm:"default:true" json:"is_required"`
	IsOptional         bool      `gorm:"default:false" json:"is_optional"`
	AlternativesAllowed bool     `gorm:"default:true" json:"alternatives_allowed"`
	Notes              *string   `gorm:"type:text" json:"notes,omitempty"`
	DisplayOrder       *int      `json:"display_order,omitempty"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	List         *SchoolSupplyList       `gorm:"foreignKey:ListID" json:"list,omitempty"`
	Product      *Product                `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Alternatives []ListItemAlternative   `gorm:"foreignKey:ListItemID" json:"alternatives,omitempty"`
}

func (SchoolSupplyListItem) TableName() string {
	return "school_supply_list_items"
}

// ListItemAlternative represents an alternative product for a list item
type ListItemAlternative struct {
	AlternativeID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"alternative_id"`
	ListItemID           uuid.UUID `gorm:"type:uuid;not null" json:"list_item_id"`
	AlternativeProductID uuid.UUID `gorm:"type:uuid;not null" json:"alternative_product_id"`
	IsRecommended        bool      `gorm:"default:false" json:"is_recommended"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	ListItem           *SchoolSupplyListItem `gorm:"foreignKey:ListItemID" json:"list_item,omitempty"`
	AlternativeProduct *Product              `gorm:"foreignKey:AlternativeProductID" json:"alternative_product,omitempty"`
}

func (ListItemAlternative) TableName() string {
	return "list_item_alternatives"
}

// Reservation represents a reservation for school supplies
type Reservation struct {
	ReservationID     uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"reservation_id"`
	ReservationNumber string            `gorm:"type:varchar(50);not null;uniqueIndex" json:"reservation_number"`
	CustomerID        uuid.UUID         `gorm:"type:uuid;not null" json:"customer_id"`
	ChildID           *uuid.UUID        `gorm:"type:uuid" json:"child_id,omitempty"`
	ListID            *uuid.UUID        `gorm:"type:uuid" json:"list_id,omitempty"`
	StoreID           *uuid.UUID        `gorm:"type:uuid" json:"store_id,omitempty"`
	Status            ReservationStatus `gorm:"type:reservation_status;default:'PENDING'" json:"status"`
	ReservationDate   time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"reservation_date"`
	ExpirationDate    time.Time         `gorm:"not null" json:"expiration_date"`
	PickupDate        *time.Time        `json:"pickup_date,omitempty"`
	TotalAmount       float64           `gorm:"type:decimal(15,2);default:0" json:"total_amount"`
	DepositAmount     float64           `gorm:"type:decimal(15,2);default:0" json:"deposit_amount"`
	Balance           float64           `gorm:"type:decimal(15,2);default:0" json:"balance"`
	Currency          CurrencyCode      `gorm:"type:currency_code;default:'VES'" json:"currency"`
	Notes             *string           `gorm:"type:text" json:"notes,omitempty"`
	ReminderSentAt    *time.Time        `json:"reminder_sent_at,omitempty"`
	CreatedAt         time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy         *uuid.UUID        `gorm:"type:uuid" json:"created_by,omitempty"`
	FulfilledAt       *time.Time        `json:"fulfilled_at,omitempty"`
	FulfilledBy       *uuid.UUID        `gorm:"type:uuid" json:"fulfilled_by,omitempty"`

	// Relations
	Customer *Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Child    *CustomerChild     `gorm:"foreignKey:ChildID" json:"child,omitempty"`
	List     *SchoolSupplyList  `gorm:"foreignKey:ListID" json:"list,omitempty"`
	Store    *Store             `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Items    []ReservationItem  `gorm:"foreignKey:ReservationID" json:"items,omitempty"`
}

func (Reservation) TableName() string {
	return "reservations"
}

// ReservationItem represents an item in a reservation
type ReservationItem struct {
	ReservationItemID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"reservation_item_id"`
	ReservationID     uuid.UUID `gorm:"type:uuid;not null" json:"reservation_id"`
	ProductID         uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity          float64   `gorm:"type:decimal(15,3);not null" json:"quantity"`
	ReservedQuantity  float64   `gorm:"type:decimal(15,3);not null" json:"reserved_quantity"`
	FulfilledQuantity float64   `gorm:"type:decimal(15,3);default:0" json:"fulfilled_quantity"`
	UnitPrice         float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	TotalAmount       float64   `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	IsFulfilled       bool      `gorm:"default:false" json:"is_fulfilled"`
	Notes             *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Reservation *Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Product     *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ReservationItem) TableName() string {
	return "reservation_items"
}

// PreOrder represents a pre-order for out-of-stock items
type PreOrder struct {
	PreOrderID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"pre_order_id"`
	PreOrderNumber     string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"pre_order_number"`
	CustomerID         uuid.UUID      `gorm:"type:uuid;not null" json:"customer_id"`
	StoreID            *uuid.UUID     `gorm:"type:uuid" json:"store_id,omitempty"`
	Status             PreOrderStatus `gorm:"type:pre_order_status;default:'PENDING'" json:"status"`
	OrderDate          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"order_date"`
	ExpectedReadyDate  *time.Time     `gorm:"type:date" json:"expected_ready_date,omitempty"`
	ReadyDate          *time.Time     `gorm:"type:date" json:"ready_date,omitempty"`
	NotificationSentAt *time.Time     `json:"notification_sent_at,omitempty"`
	PickupDeadline     *time.Time     `gorm:"type:date" json:"pickup_deadline,omitempty"`
	TotalAmount        float64        `gorm:"type:decimal(15,2);default:0" json:"total_amount"`
	DepositPaid        float64        `gorm:"type:decimal(15,2);default:0" json:"deposit_paid"`
	Currency           CurrencyCode   `gorm:"type:currency_code;default:'VES'" json:"currency"`
	Notes              *string        `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy          *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	ConfirmedAt        *time.Time     `json:"confirmed_at,omitempty"`
	ConfirmedBy        *uuid.UUID     `gorm:"type:uuid" json:"confirmed_by,omitempty"`

	// Relations
	Customer *Customer       `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Store    *Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Items    []PreOrderItem  `gorm:"foreignKey:PreOrderID" json:"items,omitempty"`
}

func (PreOrder) TableName() string {
	return "pre_orders"
}

// PreOrderItem represents an item in a pre-order
type PreOrderItem struct {
	PreOrderItemID      uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"pre_order_item_id"`
	PreOrderID          uuid.UUID  `gorm:"type:uuid;not null" json:"pre_order_id"`
	ProductID           uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	Quantity            float64    `gorm:"type:decimal(15,3);not null" json:"quantity"`
	UnitPrice           float64    `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	TotalAmount         float64    `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	IsAvailable         bool       `gorm:"default:false" json:"is_available"`
	ExpectedArrivalDate *time.Time `gorm:"type:date" json:"expected_arrival_date,omitempty"`
	Notes               *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	PreOrder *PreOrder `gorm:"foreignKey:PreOrderID" json:"pre_order,omitempty"`
	Product  *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (PreOrderItem) TableName() string {
	return "pre_order_items"
}

// CustomerNotification represents a notification sent to a customer
type CustomerNotification struct {
	NotificationID   uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"notification_id"`
	CustomerID       uuid.UUID          `gorm:"type:uuid;not null" json:"customer_id"`
	NotificationType NotificationType   `gorm:"type:notification_type;not null" json:"notification_type"`
	Status           NotificationStatus `gorm:"type:notification_status;default:'PENDING'" json:"status"`
	Subject          *string            `gorm:"type:varchar(200)" json:"subject,omitempty"`
	Message          string             `gorm:"type:text;not null" json:"message"`
	ReferenceType    *string            `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	ReferenceID      *uuid.UUID         `gorm:"type:uuid" json:"reference_id,omitempty"`
	ScheduledAt      *time.Time         `json:"scheduled_at,omitempty"`
	SentAt           *time.Time         `json:"sent_at,omitempty"`
	ReadAt           *time.Time         `json:"read_at,omitempty"`
	CreatedAt        time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (CustomerNotification) TableName() string {
	return "customer_notifications"
}

// Campaign represents a marketing campaign
type Campaign struct {
	CampaignID          uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"campaign_id"`
	CampaignName        string           `gorm:"type:varchar(200);not null" json:"campaign_name"`
	Description         *string          `gorm:"type:text" json:"description,omitempty"`
	StartDate           time.Time        `gorm:"type:date;not null" json:"start_date"`
	EndDate             time.Time        `gorm:"type:date;not null" json:"end_date"`
	DiscountPercentage  *float64         `gorm:"type:decimal(5,2)" json:"discount_percentage,omitempty"`
	TargetSchoolLevels  SchoolLevelArray `gorm:"type:school_level[]" json:"target_school_levels,omitempty"`
	TargetCategories    UUIDArray        `gorm:"type:uuid[]" json:"target_categories,omitempty"`
	TargetLocations     UUIDArray        `gorm:"type:uuid[]" json:"target_locations,omitempty"`
	IsActive            bool             `gorm:"default:true" json:"is_active"`
	Budget              *float64         `gorm:"type:decimal(15,2)" json:"budget,omitempty"`
	ActualSales         float64          `gorm:"type:decimal(15,2);default:0" json:"actual_sales"`
	CreatedAt           time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy           *uuid.UUID       `gorm:"type:uuid" json:"created_by,omitempty"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

// DemandForecast represents a demand forecast for a product
type DemandForecast struct {
	ForecastID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"forecast_id"`
	ProductID          uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	SchoolYear         string     `gorm:"type:varchar(9);not null" json:"school_year"`
	ForecastedQuantity int        `gorm:"not null" json:"forecasted_quantity"`
	ConfidenceLevel    *float64   `gorm:"type:decimal(5,2)" json:"confidence_level,omitempty"`
	ForecastMethod     *string    `gorm:"type:varchar(50)" json:"forecast_method,omitempty"`
	HistoricalData     JSONB      `gorm:"type:jsonb" json:"historical_data,omitempty"`
	CreatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy          *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (DemandForecast) TableName() string {
	return "demand_forecasts"
}
