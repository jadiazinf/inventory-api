package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role
type Role struct {
	RoleID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"role_id"`
	RoleName    string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"role_name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	BaseModel

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission represents a granular permission
type Permission struct {
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"permission_id"`
	Module       string    `gorm:"type:varchar(50);not null" json:"module"`
	Action       string    `gorm:"type:varchar(50);not null" json:"action"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission represents the many-to-many relationship
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"permission_id"`
	AssignedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// User represents a system user synchronized with Firebase
type User struct {
	UserID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"user_id"`
	FirebaseUID      string      `gorm:"type:varchar(128);not null;uniqueIndex" json:"firebase_uid"`
	NationalID       *string     `gorm:"type:varchar(20);uniqueIndex" json:"national_id,omitempty"`
	FirstName        string      `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName         string      `gorm:"type:varchar(100);not null" json:"last_name"`
	Email            string      `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	Phone            *string     `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Username         *string     `gorm:"type:varchar(50);uniqueIndex" json:"username,omitempty"`
	PhotoURL         *string     `gorm:"type:varchar(500)" json:"photo_url,omitempty"`
	RoleID           *uuid.UUID  `gorm:"type:uuid" json:"role_id,omitempty"`
	LocationID       *uuid.UUID  `gorm:"type:uuid" json:"location_id,omitempty"`
	Status           UserStatus  `gorm:"type:user_status;default:'ACTIVE'" json:"status"`
	LastLogin        *time.Time  `json:"last_login,omitempty"`
	FirebaseMetadata JSONB       `gorm:"type:jsonb" json:"firebase_metadata,omitempty"`
	BaseModelWithUser

	// Relations
	Role     *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// AuditLog represents an audit trail record
type AuditLog struct {
	AuditID      uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"audit_id"`
	UserID       *uuid.UUID      `gorm:"type:uuid" json:"user_id,omitempty"`
	Action       AuditActionType `gorm:"type:audit_action_type;not null" json:"action"`
	Module       *string         `gorm:"type:varchar(50)" json:"module,omitempty"`
	RecordID     *uuid.UUID      `gorm:"type:uuid" json:"record_id,omitempty"`
	Description  *string         `gorm:"type:text" json:"description,omitempty"`
	IPAddress    *string         `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent    *string         `gorm:"type:text" json:"user_agent,omitempty"`
	RequestData  JSONB           `gorm:"type:jsonb" json:"request_data,omitempty"`
	ResponseData JSONB           `gorm:"type:jsonb" json:"response_data,omitempty"`
	CreatedAt    time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// UserSession represents a user session for tracking beyond Firebase
type UserSession struct {
	SessionID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"session_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	FirebaseTokenID *string   `gorm:"type:varchar(500)" json:"firebase_token_id,omitempty"`
	RefreshToken   *string   `gorm:"type:varchar(500)" json:"refresh_token,omitempty"`
	IPAddress      *string   `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent      *string   `gorm:"type:text" json:"user_agent,omitempty"`
	ExpiresAt      time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	LastActivity   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_activity"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
