package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

type accountsReceivableService struct {
	arRepo repositories.AccountsReceivableRepository
	db     *gorm.DB
}

// NewAccountsReceivableService creates a new accounts receivable service
func NewAccountsReceivableService(
	arRepo repositories.AccountsReceivableRepository,
	db *gorm.DB,
) services.AccountsReceivableService {
	return &accountsReceivableService{
		arRepo: arRepo,
		db:     db,
	}
}

// GetAccountsReceivable retrieves an accounts receivable by ID
func (s *accountsReceivableService) GetAccountsReceivable(ctx context.Context, id uuid.UUID) (*domain.AccountsReceivable, error) {
	return s.arRepo.FindByID(ctx, id)
}

// GetCustomerReceivables retrieves all receivables for a customer
func (s *accountsReceivableService) GetCustomerReceivables(ctx context.Context, customerID uuid.UUID) ([]domain.AccountsReceivable, error) {
	return s.arRepo.FindByCustomer(ctx, customerID)
}

// GetOverdueReceivables retrieves all overdue receivables
func (s *accountsReceivableService) GetOverdueReceivables(ctx context.Context) ([]domain.AccountsReceivable, error) {
	return s.arRepo.GetOverdue(ctx)
}

// RegisterPayment registers a payment for an accounts receivable
func (s *accountsReceivableService) RegisterPayment(
	ctx context.Context,
	receivableID uuid.UUID,
	amount float64,
	currency domain.CurrencyCode,
	paymentMethod domain.PaymentMethod,
	reference, notes *string,
	userID uuid.UUID,
) error {
	// Get accounts receivable
	ar, err := s.arRepo.FindByID(ctx, receivableID)
	if err != nil {
		return err
	}

	// Validate payment
	if ar.Status == domain.AccountStatusPaid {
		return errors.InvalidInput("Receivable is already fully paid")
	}

	if ar.Status == domain.AccountStatusCancelled {
		return errors.InvalidInput("Receivable is cancelled")
	}

	if amount <= 0 {
		return errors.InvalidInput("Payment amount must be positive")
	}

	if amount > ar.Balance {
		return errors.InvalidInput(fmt.Sprintf(
			"Payment amount (%.2f) exceeds balance (%.2f)",
			amount, ar.Balance,
		))
	}

	// Validate currency matches
	if currency != ar.Currency {
		return errors.InvalidInput(fmt.Sprintf(
			"Payment currency (%s) does not match receivable currency (%s)",
			currency, ar.Currency,
		))
	}

	// Create payment record
	payment := &domain.CustomerPayment{
		PaymentID:     uuid.New(),
		ReceivableID:  receivableID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: paymentMethod,
		PaymentDate:   time.Now(),
		Reference:     reference,
		Notes:         notes,
		CreatedBy:     &userID,
	}

	// Add payment (updates AR balance)
	if err := s.arRepo.AddPayment(ctx, payment); err != nil {
		return err
	}

	return nil
}

// GetPaymentHistory retrieves payment history for a receivable
func (s *accountsReceivableService) GetPaymentHistory(ctx context.Context, receivableID uuid.UUID) ([]domain.CustomerPayment, error) {
	return s.arRepo.GetPayments(ctx, receivableID)
}
