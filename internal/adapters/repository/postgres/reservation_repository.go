package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new reservation repository
func NewReservationRepository(db *gorm.DB) repositories.ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		return errors.WrapError(err, "failed to create reservation")
	}
	return nil
}

func (r *reservationRepository) CreateWithItems(ctx context.Context, reservation *domain.Reservation, items []domain.ReservationItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Generate unique reservation number if not provided
		if reservation.ReservationNumber == "" {
			resNum, err := r.generateReservationNumber(tx)
			if err != nil {
				return err
			}
			reservation.ReservationNumber = resNum
		}

		// 2. Create reservation record
		if err := tx.Create(reservation).Error; err != nil {
			return errors.WrapError(err, "failed to create reservation")
		}

		// 3. Get warehouse from store (assume first warehouse of store)
		var warehouse domain.Warehouse
		if err := tx.Where("store_id = ?", reservation.StoreID).First(&warehouse).Error; err != nil {
			return errors.WrapError(err, "failed to find warehouse for store")
		}

		// 4. Create reservation items and inventory movements
		for i := range items {
			items[i].ReservationID = reservation.ReservationID

			// Calculate amounts
			items[i].TotalAmount = items[i].Quantity * items[i].UnitPrice
			items[i].ReservedQuantity = items[i].Quantity
			items[i].FulfilledQuantity = 0
			items[i].IsFulfilled = false

			// Create reservation item
			if err := tx.Create(&items[i]).Error; err != nil {
				return errors.WrapError(err, "failed to create reservation item")
			}

			// Check inventory availability
			var inventory domain.Inventory
			err := tx.Where("product_id = ? AND warehouse_id = ?",
				items[i].ProductID, warehouse.WarehouseID).First(&inventory).Error

			if err == gorm.ErrRecordNotFound {
				return errors.InsufficientStock("Product", 0, items[i].Quantity)
			} else if err != nil {
				return errors.WrapError(err, "failed to check inventory")
			}

			// Verify available quantity
			if inventory.AvailableQuantity < items[i].Quantity {
				// Get product name for error message
				var product domain.Product
				tx.First(&product, "product_id = ?", items[i].ProductID)
				return errors.InsufficientStock(
					product.Name,
					inventory.AvailableQuantity,
					items[i].Quantity,
				)
			}

			// Create RESERVATION inventory movement (trigger updates inventory)
			movement := &domain.InventoryMovement{
				MovementID:    uuid.New(),
				ProductID:     items[i].ProductID,
				WarehouseID:   warehouse.WarehouseID,
				MovementType:  domain.MovementTypeReservation,
				Quantity:      items[i].Quantity,
				UnitCost:      &items[i].UnitPrice,
				Currency:      reservation.Currency,
				ReferenceType: stringPtr("RESERVATION"),
				ReferenceID:   &reservation.ReservationID,
				CreatedBy:     reservation.CreatedBy,
			}
			if err := tx.Create(movement).Error; err != nil {
				return errors.WrapError(err, "failed to create reservation inventory movement")
			}
		}

		return nil
	})
}

func (r *reservationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	var reservation domain.Reservation
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Child").
		Preload("List").
		Preload("Store").
		Preload("Items").
		Preload("Items.Product").
		First(&reservation, "reservation_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Reservation", id.String())
		}
		return nil, errors.WrapError(err, "failed to find reservation")
	}
	return &reservation, nil
}

func (r *reservationRepository) FindByNumber(ctx context.Context, reservationNumber string) (*domain.Reservation, error) {
	var reservation domain.Reservation
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Items").
		Preload("Items.Product").
		Where("reservation_number = ?", reservationNumber).
		First(&reservation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Reservation")
		}
		return nil, errors.WrapError(err, "failed to find reservation by number")
	}
	return &reservation, nil
}

func (r *reservationRepository) List(ctx context.Context, filters repositories.ReservationFilters, limit, offset int) ([]domain.Reservation, int64, error) {
	var reservations []domain.Reservation
	var total int64

	query := r.buildFilterQuery(r.db.WithContext(ctx).Model(&domain.Reservation{}), filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count reservations")
	}

	err := query.
		Preload("Customer").
		Preload("Store").
		Preload("Child").
		Order("reservation_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&reservations).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to list reservations")
	}

	return reservations, total, nil
}

func (r *reservationRepository) GetItems(ctx context.Context, reservationID uuid.UUID) ([]domain.ReservationItem, error) {
	var items []domain.ReservationItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Where("reservation_id = ?", reservationID).
		Find(&items).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get reservation items")
	}
	return items, nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *domain.Reservation) error {
	if err := r.db.WithContext(ctx).Save(reservation).Error; err != nil {
		return errors.WrapError(err, "failed to update reservation")
	}
	return nil
}

func (r *reservationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReservationStatus) error {
	err := r.db.WithContext(ctx).
		Model(&domain.Reservation{}).
		Where("reservation_id = ?", id).
		Update("status", status).Error

	if err != nil {
		return errors.WrapError(err, "failed to update reservation status")
	}
	return nil
}

func (r *reservationRepository) MarkAsFulfilled(ctx context.Context, id uuid.UUID, fulfilledBy uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&domain.Reservation{}).
		Where("reservation_id = ?", id).
		Updates(map[string]interface{}{
			"status":       domain.ReservationStatusFulfilled,
			"fulfilled_at": now,
			"fulfilled_by": fulfilledBy,
		}).Error

	if err != nil {
		return errors.WrapError(err, "failed to mark reservation as fulfilled")
	}
	return nil
}

func (r *reservationRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get reservation with items
		var reservation domain.Reservation
		if err := tx.Preload("Items").Preload("Store").
			First(&reservation, "reservation_id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NotFoundWithID("Reservation", id.String())
			}
			return errors.WrapError(err, "failed to find reservation")
		}

		// Can only cancel pending or confirmed reservations
		if reservation.Status != domain.ReservationStatusPending &&
		   reservation.Status != domain.ReservationStatusConfirmed {
			return errors.BadRequest("Can only cancel pending or confirmed reservations")
		}

		// Get warehouse
		var warehouse domain.Warehouse
		if err := tx.Where("store_id = ?", reservation.StoreID).First(&warehouse).Error; err != nil {
			return errors.WrapError(err, "failed to find warehouse")
		}

		// Update reservation status
		if err := tx.Model(&reservation).Update("status", domain.ReservationStatusCancelled).Error; err != nil {
			return errors.WrapError(err, "failed to cancel reservation")
		}

		// Create RESERVATION_RELEASE movements to free up inventory
		for _, item := range reservation.Items {
			movement := &domain.InventoryMovement{
				MovementID:    uuid.New(),
				ProductID:     item.ProductID,
				WarehouseID:   warehouse.WarehouseID,
				MovementType:  domain.MovementTypeReservationRelease,
				Quantity:      item.ReservedQuantity,
				ReferenceType: stringPtr("RESERVATION_CANCELLATION"),
				ReferenceID:   &reservation.ReservationID,
				Notes:         stringPtr("Release from cancelled reservation"),
			}
			if err := tx.Create(movement).Error; err != nil {
				return errors.WrapError(err, "failed to create release inventory movement")
			}
		}

		return nil
	})
}

func (r *reservationRepository) GetExpired(ctx context.Context) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	now := time.Now()

	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Items").
		Where("expiration_date < ?", now).
		Where("status IN ?", []domain.ReservationStatus{
			domain.ReservationStatusPending,
			domain.ReservationStatusConfirmed,
		}).
		Find(&reservations).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get expired reservations")
	}

	return reservations, nil
}

func (r *reservationRepository) GetExpiringFor(ctx context.Context, within time.Duration) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	now := time.Now()
	expiryThreshold := now.Add(within)

	err := r.db.WithContext(ctx).
		Preload("Customer").
		Where("expiration_date BETWEEN ? AND ?", now, expiryThreshold).
		Where("status IN ?", []domain.ReservationStatus{
			domain.ReservationStatusPending,
			domain.ReservationStatusConfirmed,
		}).
		Where("reminder_sent_at IS NULL").
		Find(&reservations).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get expiring reservations")
	}

	return reservations, nil
}

// Helper functions

func (r *reservationRepository) buildFilterQuery(query *gorm.DB, filters repositories.ReservationFilters) *gorm.DB {
	if filters.CustomerID != nil {
		query = query.Where("customer_id = ?", *filters.CustomerID)
	}

	if filters.StoreID != nil {
		query = query.Where("store_id = ?", *filters.StoreID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.DateFrom != nil {
		query = query.Where("reservation_date >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("reservation_date <= ?", *filters.DateTo)
	}

	return query
}

func (r *reservationRepository) generateReservationNumber(tx *gorm.DB) (string, error) {
	// Get current year and month
	now := time.Now()
	prefix := now.Format("RES-2006-01")

	// Get count of reservations this month
	var count int64
	if err := tx.Model(&domain.Reservation{}).
		Where("reservation_number LIKE ?", prefix+"%").
		Count(&count).Error; err != nil {
		return "", errors.WrapError(err, "failed to count reservations for number generation")
	}

	// Generate reservation number: RES-YYYY-MM-NNNN
	reservationNumber := fmt.Sprintf("%s-%04d", prefix, count+1)
	return reservationNumber, nil
}

func stringPtr(s string) *string {
	return &s
}
