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

func setupCustomerChildTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate tables
	err = db.AutoMigrate(
		&domain.Customer{},
		&domain.CustomerChild{},
	)
	require.NoError(t, err)

	return db
}

func TestCustomerChildRepository_Create(t *testing.T) {
	db := setupCustomerChildTestDB(t)
	repo := NewCustomerChildRepository(db)
	ctx := context.Background()

	// Create parent customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("Parent"),
		LastName:   stringPtr("Customer"),
	}
	require.NoError(t, db.Create(customer).Error)

	// Create child
	birthDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	schoolLevel := domain.SchoolLevelPrimary
	child := &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customer.CustomerID,
		FirstName:   "Child",
		LastName:    "Doe",
		DateOfBirth: &birthDate,
		SchoolLevel: &schoolLevel,
	}

	err := repo.Create(ctx, child)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, child.ChildID)
}

func TestCustomerChildRepository_FindByID(t *testing.T) {
	db := setupCustomerChildTestDB(t)
	repo := NewCustomerChildRepository(db)
	ctx := context.Background()

	// Create parent customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("Parent"),
		LastName:   stringPtr("Customer"),
	}
	require.NoError(t, db.Create(customer).Error)

	// Create child
	birthDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	schoolLevel := domain.SchoolLevelPrimary
	child := &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customer.CustomerID,
		FirstName:   "Child",
		LastName:    "Doe",
		DateOfBirth: &birthDate,
		SchoolLevel: &schoolLevel,
	}
	require.NoError(t, db.Create(child).Error)

	// Test FindByID
	found, err := repo.FindByID(ctx, child.ChildID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, child.ChildID, found.ChildID)
	assert.Equal(t, child.FirstName, found.FirstName)
	assert.Equal(t, child.LastName, found.LastName)
}

func TestCustomerChildRepository_FindByCustomer(t *testing.T) {
	db := setupCustomerChildTestDB(t)
	repo := NewCustomerChildRepository(db)
	ctx := context.Background()

	// Create parent customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("Parent"),
		LastName:   stringPtr("Customer"),
	}
	require.NoError(t, db.Create(customer).Error)

	// Create multiple children
	birthDate1 := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	schoolLevel1 := domain.SchoolLevelPrimary
	child1 := &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customer.CustomerID,
		FirstName:   "Child1",
		LastName:    "Doe",
		DateOfBirth: &birthDate1,
		SchoolLevel: &schoolLevel1,
	}
	birthDate2 := time.Date(2017, 6, 15, 0, 0, 0, 0, time.UTC)
	schoolLevel2 := domain.SchoolLevelPreschool
	child2 := &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customer.CustomerID,
		FirstName:   "Child2",
		LastName:    "Doe",
		DateOfBirth: &birthDate2,
		SchoolLevel: &schoolLevel2,
	}
	require.NoError(t, db.Create(child1).Error)
	require.NoError(t, db.Create(child2).Error)

	// Test FindByCustomer
	children, err := repo.FindByCustomer(ctx, customer.CustomerID)
	assert.NoError(t, err)
	assert.Len(t, children, 2)
}

func TestCustomerChildRepository_Update(t *testing.T) {
	db := setupCustomerChildTestDB(t)
	repo := NewCustomerChildRepository(db)
	ctx := context.Background()

	// Create parent customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("Parent"),
		LastName:   stringPtr("Customer"),
	}
	require.NoError(t, db.Create(customer).Error)

	// Create child
	birthDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	schoolLevel := domain.SchoolLevelPrimary
	child := &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customer.CustomerID,
		FirstName:   "Child",
		LastName:    "Doe",
		DateOfBirth: &birthDate,
		SchoolLevel: &schoolLevel,
	}
	require.NoError(t, db.Create(child).Error)

	// Update child
	newSchoolLevel := domain.SchoolLevelHighSchool
	child.SchoolLevel = &newSchoolLevel
	child.FirstName = "UpdatedChild"

	err := repo.Update(ctx, child)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.FindByID(ctx, child.ChildID)
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedChild", updated.FirstName)
	assert.Equal(t, domain.SchoolLevelHighSchool, *updated.SchoolLevel)
}

func TestCustomerChildRepository_Delete(t *testing.T) {
	db := setupCustomerChildTestDB(t)
	repo := NewCustomerChildRepository(db)
	ctx := context.Background()

	// Create parent customer
	customer := &domain.Customer{
		CustomerID: uuid.New(),
		TaxID:      "V12345678",
		FirstName:  stringPtr("Parent"),
		LastName:   stringPtr("Customer"),
	}
	require.NoError(t, db.Create(customer).Error)

	// Create child
	birthDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	schoolLevel := domain.SchoolLevelPrimary
	child := &domain.CustomerChild{
		ChildID:     uuid.New(),
		CustomerID:  customer.CustomerID,
		FirstName:   "Child",
		LastName:    "Doe",
		DateOfBirth: &birthDate,
		SchoolLevel: &schoolLevel,
	}
	require.NoError(t, db.Create(child).Error)

	// Delete child
	err := repo.Delete(ctx, child.ChildID)
	assert.NoError(t, err)

	// Verify deletion (soft delete)
	_, err = repo.FindByID(ctx, child.ChildID)
	assert.Error(t, err) // Should return not found error
}
