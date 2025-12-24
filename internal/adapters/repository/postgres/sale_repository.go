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

type saleRepository struct {
	db *gorm.DB
}

// NewSaleRepository creates a new sale repository
func NewSaleRepository(db *gorm.DB) repositories.SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	if err := r.db.WithContext(ctx).Create(sale).Error; err != nil {
		return errors.WrapError(err, "failed to create sale")
	}
	return nil
}

func (r *saleRepository) CreateWithDetails(ctx context.Context, sale *domain.Sale, details []domain.SaleDetail) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Generate unique invoice number if not provided
		if sale.InvoiceNumber == "" {
			invoiceNum, err := r.generateInvoiceNumber(tx)
			if err != nil {
				return err
			}
			sale.InvoiceNumber = invoiceNum
		}

		// 2. Create sale record
		if err := tx.Create(sale).Error; err != nil {
			return errors.WrapError(err, "failed to create sale")
		}

		// 3. Create sale details (trigger will calculate totals)
		for i := range details {
			details[i].SaleID = sale.SaleID

			// Calculate subtotal and total for detail
			details[i].Subtotal = details[i].Quantity * details[i].UnitPrice - details[i].DiscountAmount
			details[i].TaxAmount = details[i].Subtotal * (details[i].TaxPercentage / 100)
			details[i].Total = details[i].Subtotal + details[i].TaxAmount

			if err := tx.Create(&details[i]).Error; err != nil {
				return errors.WrapError(err, "failed to create sale detail")
			}
		}

		// 4. Create inventory movements for completed sales
		if sale.Status == domain.SaleStatusCompleted && sale.WarehouseID != nil {
			for _, detail := range details {
				movement := &domain.InventoryMovement{
					MovementID:    uuid.New(),
					ProductID:     detail.ProductID,
					WarehouseID:   *sale.WarehouseID,
					MovementType:  domain.MovementTypeOut,
					Quantity:      detail.Quantity,
					UnitCost:      &detail.UnitPrice,
					Currency:      sale.Currency,
					ReferenceType: stringPtr("SALE"),
					ReferenceID:   &sale.SaleID,
					CreatedBy:     sale.CreatedBy,
				}
				if err := tx.Create(movement).Error; err != nil {
					return errors.WrapError(err, "failed to create inventory movement")
				}
			}
		}

		return nil
	})
}

func (r *saleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	var sale domain.Sale
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Store").
		Preload("Warehouse").
		Preload("Salesperson").
		Preload("Details").
		Preload("Details.Product").
		First(&sale, "sale_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Sale", id.String())
		}
		return nil, errors.WrapError(err, "failed to find sale")
	}
	return &sale, nil
}

func (r *saleRepository) FindByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Sale, error) {
	var sale domain.Sale
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Store").
		Preload("Details").
		Preload("Details.Product").
		Where("invoice_number = ?", invoiceNumber).
		First(&sale).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Sale")
		}
		return nil, errors.WrapError(err, "failed to find sale by invoice number")
	}
	return &sale, nil
}

func (r *saleRepository) List(ctx context.Context, filters repositories.SaleFilters, limit, offset int) ([]domain.Sale, int64, error) {
	var sales []domain.Sale
	var total int64

	query := r.buildFilterQuery(r.db.WithContext(ctx).Model(&domain.Sale{}), filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count sales")
	}

	err := query.
		Preload("Customer").
		Preload("Store").
		Preload("Salesperson").
		Order("sale_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&sales).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to list sales")
	}

	return sales, total, nil
}

func (r *saleRepository) GetDetails(ctx context.Context, saleID uuid.UUID) ([]domain.SaleDetail, error) {
	var details []domain.SaleDetail
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Where("sale_id = ?", saleID).
		Find(&details).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get sale details")
	}
	return details, nil
}

func (r *saleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	if err := r.db.WithContext(ctx).Save(sale).Error; err != nil {
		return errors.WrapError(err, "failed to update sale")
	}
	return nil
}

func (r *saleRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get sale with details
		var sale domain.Sale
		if err := tx.Preload("Details").First(&sale, "sale_id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NotFoundWithID("Sale", id.String())
			}
			return errors.WrapError(err, "failed to find sale")
		}

		// Can only cancel completed sales
		if sale.Status != domain.SaleStatusCompleted {
			return errors.BadRequest("Can only cancel completed sales")
		}

		// Update sale status
		if err := tx.Model(&sale).Update("status", domain.SaleStatusCancelled).Error; err != nil {
			return errors.WrapError(err, "failed to cancel sale")
		}

		// Create reverse inventory movements (IN)
		if sale.WarehouseID != nil {
			for _, detail := range sale.Details {
				movement := &domain.InventoryMovement{
					MovementID:    uuid.New(),
					ProductID:     detail.ProductID,
					WarehouseID:   *sale.WarehouseID,
					MovementType:  domain.MovementTypeIn,
					Quantity:      detail.Quantity,
					ReferenceType: stringPtr("SALE_CANCELLATION"),
					ReferenceID:   &sale.SaleID,
					Notes:         stringPtr("Reversal from cancelled sale"),
				}
				if err := tx.Create(movement).Error; err != nil {
					return errors.WrapError(err, "failed to create reverse inventory movement")
				}
			}
		}

		return nil
	})
}

func (r *saleRepository) GetDailySales(ctx context.Context, storeID *uuid.UUID, date time.Time) ([]domain.Sale, error) {
	var sales []domain.Sale

	// Get start and end of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Salesperson").
		Where("sale_date >= ? AND sale_date < ?", startOfDay, endOfDay).
		Where("status = ?", domain.SaleStatusCompleted)

	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	err := query.Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, errors.WrapError(err, "failed to get daily sales")
	}

	return sales, nil
}

func (r *saleRepository) GetSalesByPeriod(ctx context.Context, from, to time.Time) ([]domain.Sale, error) {
	var sales []domain.Sale
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Store").
		Where("sale_date >= ? AND sale_date <= ?", from, to).
		Where("status = ?", domain.SaleStatusCompleted).
		Order("sale_date DESC").
		Find(&sales).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get sales by period")
	}

	return sales, nil
}

// Helper functions

func (r *saleRepository) buildFilterQuery(query *gorm.DB, filters repositories.SaleFilters) *gorm.DB {
	if filters.CustomerID != nil {
		query = query.Where("customer_id = ?", *filters.CustomerID)
	}

	if filters.StoreID != nil {
		query = query.Where("store_id = ?", *filters.StoreID)
	}

	if filters.SalespersonID != nil {
		query = query.Where("salesperson_id = ?", *filters.SalespersonID)
	}

	if filters.SaleType != nil {
		query = query.Where("sale_type = ?", *filters.SaleType)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.DateFrom != nil {
		query = query.Where("sale_date >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("sale_date <= ?", *filters.DateTo)
	}

	return query
}

func (r *saleRepository) generateInvoiceNumber(tx *gorm.DB) (string, error) {
	// Get current year and month
	now := time.Now()
	prefix := now.Format("2006-01")

	// Get count of sales this month
	var count int64
	if err := tx.Model(&domain.Sale{}).
		Where("invoice_number LIKE ?", prefix+"%").
		Count(&count).Error; err != nil {
		return "", errors.WrapError(err, "failed to count sales for invoice number")
	}

	// Generate invoice number: YYYY-MM-NNNN
	invoiceNumber := fmt.Sprintf("%s-%04d", prefix, count+1)
	return invoiceNumber, nil
}
