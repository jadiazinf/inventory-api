package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// LocationRepository defines the interface for location data access
type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	FindByCode(ctx context.Context, code string) (*domain.Location, error)
	FindByType(ctx context.Context, locationType domain.LocationType) ([]domain.Location, error)
	FindChildren(ctx context.Context, parentID uuid.UUID) ([]domain.Location, error)
	List(ctx context.Context, limit, offset int) ([]domain.Location, int64, error)
	Update(ctx context.Context, location *domain.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
}
