package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	List(ctx context.Context) ([]domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error)
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	RevokePermission(ctx context.Context, roleID, permissionID uuid.UUID) error
}

// PermissionRepository defines the interface for permission data access
type PermissionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error)
	FindByModuleAndAction(ctx context.Context, module, action string) (*domain.Permission, error)
	List(ctx context.Context) ([]domain.Permission, error)
	HasPermission(ctx context.Context, userID uuid.UUID, module, action string) (bool, error)
}
