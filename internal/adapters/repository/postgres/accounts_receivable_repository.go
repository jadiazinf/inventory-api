package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"gorm.io/gorm"
)

type accountsReceivableRepository struct {
	db *gorm.DB
}

// NewAccountsReceivableRepository creates a new accounts receivable repository
func NewAccountsReceivableRepository(db *gorm.DB) repositories.AccountsReceivableRepository {
	return &accountsReceivableRepository{db: db}
}

func (r *accountsReceivableRepository) Create(ctx context.Context, receivable *domain.AccountsReceivable) error {
	if err := r.db.WithContext(ctx).Create(receivable).Error; err != nil {
		return errors.WrapError(err, "failed to create accounts receivable")
	}
	return nil
}

func (r *accountsReceivableRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.AccountsReceivable, error) {
	var receivable domain.AccountsReceivable
	err := r.db.WithContext(ctx).
		Preload("Sale").
		Preload("Customer").
		First(&receivable, "account_id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("AccountsReceivable", id.String())
		}
		return nil, errors.WrapError(err, "failed to find accounts receivable")
	}
	return &receivable, nil
}

func (r *accountsReceivableRepository) FindBySale(ctx context.Context, saleID uuid.UUID) (*domain.AccountsReceivable, error) {
	var receivable domain.AccountsReceivable
	err := r.db.WithContext(ctx).
		Preload("Sale").
		Preload("Customer").
		Where("sale_id = ?", saleID).
		First(&receivable).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("AccountsReceivable")
		}
		return nil, errors.WrapError(err, "failed to find accounts receivable by sale")
	}
	return &receivable, nil
}

func (r *accountsReceivableRepository) FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]domain.AccountsReceivable, error) {
	var receivables []domain.AccountsReceivable
	err := r.db.WithContext(ctx).
		Preload("Sale").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&receivables).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to find accounts receivable by customer")
	}
	return receivables, nil
}

func (r *accountsReceivableRepository) GetOverdue(ctx context.Context) ([]domain.AccountsReceivable, error) {
	var receivables []domain.AccountsReceivable
	now := time.Now()

	err := r.db.WithContext(ctx).
		Preload("Sale").
		Preload("Customer").
		Where("status IN (?, ?)", domain.AccountStatusPending, domain.AccountStatusPartiallyPaid).
		Where("due_date < ?", now).
		Order("due_date ASC").
		Find(&receivables).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get overdue accounts")
	}
	return receivables, nil
}

func (r *accountsReceivableRepository) List(ctx context.Context, status *domain.AccountStatus, limit, offset int) ([]domain.AccountsReceivable, int64, error) {
	var receivables []domain.AccountsReceivable
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.AccountsReceivable{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapError(err, "failed to count accounts receivable")
	}

	err := query.
		Preload("Sale").
		Preload("Customer").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&receivables).Error

	if err != nil {
		return nil, 0, errors.WrapError(err, "failed to list accounts receivable")
	}

	return receivables, total, nil
}

func (r *accountsReceivableRepository) Update(ctx context.Context, receivable *domain.AccountsReceivable) error {
	if err := r.db.WithContext(ctx).Save(receivable).Error; err != nil {
		return errors.WrapError(err, "failed to update accounts receivable")
	}
	return nil
}

func (r *accountsReceivableRepository) AddPayment(ctx context.Context, payment *domain.CustomerPayment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create payment record
		if err := tx.Create(payment).Error; err != nil {
			return errors.WrapError(err, "failed to create payment")
		}

		// Get the accounts receivable record
		var receivable domain.AccountsReceivable
		if err := tx.First(&receivable, "account_id = ?", payment.ReceivableID).Error; err != nil {
			return errors.WrapError(err, "failed to find accounts receivable")
		}

		// Update paid amount
		receivable.PaidAmount += payment.Amount

		// Update status based on paid amount
		if receivable.PaidAmount >= receivable.TotalAmount {
			receivable.Status = domain.AccountStatusPaid
			receivable.PaidAmount = receivable.TotalAmount // Ensure we don't overpay
		} else if receivable.PaidAmount > 0 {
			receivable.Status = domain.AccountStatusPartiallyPaid
		}

		// Update balance
		receivable.Balance = receivable.TotalAmount - receivable.PaidAmount

		// Save updated receivable
		if err := tx.Save(&receivable).Error; err != nil {
			return errors.WrapError(err, "failed to update accounts receivable")
		}

		return nil
	})
}

func (r *accountsReceivableRepository) GetPayments(ctx context.Context, receivableID uuid.UUID) ([]domain.CustomerPayment, error) {
	var payments []domain.CustomerPayment
	err := r.db.WithContext(ctx).
		Where("account_id = ?", receivableID).
		Order("payment_date DESC").
		Find(&payments).Error

	if err != nil {
		return nil, errors.WrapError(err, "failed to get payments")
	}
	return payments, nil
}
