package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate all tables
	err = db.AutoMigrate(
		&domain.AccountsReceivable{},
		&domain.CustomerPayment{},
		&domain.Sale{},
		&domain.Customer{},
	)
	require.NoError(t, err)

	return db
}

func TestAccountsReceivableRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountsReceivableRepository(db)
	ctx := context.Background()

	// Create test customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("John"),
		LastName:   stringPtr("Doe"),
		Email:      stringPtr("john@example.com"),
	}
	require.NoError(t, db.Create(customer).Error)

	// Create test sale
	sale := &domain.Sale{
		SaleID:         uuid.New(),
		CustomerID:     &customer.CustomerID,
		InvoiceNumber:  "INV-001",
		TotalAmount:    1000.0,
		Currency:       domain.CurrencyVES,
		SaleType:       domain.SaleTypeCash,
		Status:         domain.SaleStatusCompleted,
	}
	require.NoError(t, db.Create(sale).Error)

	// Create accounts receivable
	ar := &domain.AccountsReceivable{
		ReceivableID: uuid.New(),
		SaleID:       &sale.SaleID,
		CustomerID:   customer.CustomerID,
		TotalAmount:  1000.0,
		PaidAmount:   0,
		Balance:      1000.0,
		DueDate:      time.Now().AddDate(0, 0, 30),
		Status:       domain.AccountStatusPending,
		Currency:     domain.CurrencyVES,
	}

	err := repo.Create(ctx, ar)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, ar.ReceivableID)
}

func TestAccountsReceivableRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountsReceivableRepository(db)
	ctx := context.Background()

	// Create test data
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("John"),
		LastName:   stringPtr("Doe"),
	}
	require.NoError(t, db.Create(customer).Error)

	sale := &domain.Sale{
		SaleID:        uuid.New(),
		CustomerID:    &customer.CustomerID,
		InvoiceNumber: "INV-001",
		TotalAmount:   1000.0,
		Currency:      domain.CurrencyVES,
		SaleType:      domain.SaleTypeCash,
		Status:        domain.SaleStatusCompleted,
	}
	require.NoError(t, db.Create(sale).Error)

	ar := &domain.AccountsReceivable{
		ReceivableID: uuid.New(),
		SaleID:       &sale.SaleID,
		CustomerID:   customer.CustomerID,
		TotalAmount:  1000.0,
		PaidAmount:   0,
		Balance:      1000.0,
		DueDate:      time.Now().AddDate(0, 0, 30),
		Status:       domain.AccountStatusPending,
		Currency:     domain.CurrencyVES,
	}
	require.NoError(t, db.Create(ar).Error)

	// Test FindByID
	found, err := repo.FindByID(ctx, ar.ReceivableID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, ar.ReceivableID, found.ReceivableID)
	assert.Equal(t, ar.TotalAmount, found.TotalAmount)
}

func TestAccountsReceivableRepository_AddPayment(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountsReceivableRepository(db)
	ctx := context.Background()

	// Create test data
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("John"),
		LastName:   stringPtr("Doe"),
	}
	require.NoError(t, db.Create(customer).Error)

	sale := &domain.Sale{
		SaleID:        uuid.New(),
		CustomerID:    &customer.CustomerID,
		InvoiceNumber: "INV-001",
		TotalAmount:   1000.0,
		Currency:      domain.CurrencyVES,
		SaleType:      domain.SaleTypeCredit,
		Status:        domain.SaleStatusCompleted,
	}
	require.NoError(t, db.Create(sale).Error)

	ar := &domain.AccountsReceivable{
		ReceivableID: uuid.New(),
		SaleID:       &sale.SaleID,
		CustomerID:   customer.CustomerID,
		TotalAmount:  1000.0,
		PaidAmount:   0,
		Balance:      1000.0,
		DueDate:      time.Now().AddDate(0, 0, 30),
		Status:       domain.AccountStatusPending,
		Currency:     domain.CurrencyVES,
	}
	require.NoError(t, db.Create(ar).Error)

	// Test partial payment
	payment := &domain.CustomerPayment{
		PaymentID:     uuid.New(),
		ReceivableID:  ar.ReceivableID,
		Amount:        400.0,
		Currency:      domain.CurrencyVES,
		PaymentDate:   time.Now(),
		PaymentMethod: domain.PaymentMethodCash,
	}

	err := repo.AddPayment(ctx, payment)
	assert.NoError(t, err)

	// Verify updated status
	updated, err := repo.FindByID(ctx, ar.ReceivableID)
	assert.NoError(t, err)
	assert.Equal(t, 400.0, updated.PaidAmount)
	assert.Equal(t, 600.0, updated.Balance)
	assert.Equal(t, domain.AccountStatusPartiallyPaid, updated.Status)

	// Test full payment
	payment2 := &domain.CustomerPayment{
		PaymentID:     uuid.New(),
		ReceivableID:  ar.ReceivableID,
		Amount:        600.0,
		Currency:      domain.CurrencyVES,
		PaymentDate:   time.Now(),
		PaymentMethod: domain.PaymentMethodCash,
	}

	err = repo.AddPayment(ctx, payment2)
	assert.NoError(t, err)

	// Verify fully paid status
	updated, err = repo.FindByID(ctx, ar.ReceivableID)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, updated.PaidAmount)
	assert.Equal(t, 0.0, updated.Balance)
	assert.Equal(t, domain.AccountStatusPaid, updated.Status)
}

func TestAccountsReceivableRepository_GetOverdue(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountsReceivableRepository(db)
	ctx := context.Background()

	// Create test customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("John"),
		LastName:   stringPtr("Doe"),
	}
	require.NoError(t, db.Create(customer).Error)

	sale := &domain.Sale{
		SaleID:        uuid.New(),
		CustomerID:    &customer.CustomerID,
		InvoiceNumber: "INV-001",
		TotalAmount:   1000.0,
		Currency:      domain.CurrencyVES,
		SaleType:      domain.SaleTypeCredit,
		Status:        domain.SaleStatusCompleted,
	}
	require.NoError(t, db.Create(sale).Error)

	// Create overdue AR
	overdueAR := &domain.AccountsReceivable{
		ReceivableID: uuid.New(),
		SaleID:       &sale.SaleID,
		CustomerID:   customer.CustomerID,
		TotalAmount:  1000.0,
		PaidAmount:   0,
		Balance:      1000.0,
		DueDate:      time.Now().AddDate(0, 0, -10), // 10 days overdue
		Status:       domain.AccountStatusPending,
		Currency:     domain.CurrencyVES,
	}
	require.NoError(t, db.Create(overdueAR).Error)

	// Get overdue accounts
	overdue, err := repo.GetOverdue(ctx)
	assert.NoError(t, err)
	assert.Len(t, overdue, 1)
	assert.Equal(t, overdueAR.ReceivableID, overdue[0].ReceivableID)
}
