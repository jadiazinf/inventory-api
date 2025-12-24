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

type saleService struct {
	saleRepo      repositories.SaleRepository
	productRepo   repositories.ProductRepository
	inventoryRepo repositories.InventoryRepository
	customerRepo  repositories.CustomerRepository
	db            *gorm.DB
}

// NewSaleService creates a new sale service
func NewSaleService(
	saleRepo repositories.SaleRepository,
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
	customerRepo repositories.CustomerRepository,
	db *gorm.DB,
) services.SaleService {
	return &saleService{
		saleRepo:      saleRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		customerRepo:  customerRepo,
		db:            db,
	}
}

// CreateSale creates a new sale
func (s *saleService) CreateSale(ctx context.Context, req services.CreateSaleRequest) (*domain.Sale, error) {
	// Validate customer if provided
	if req.CustomerID != nil {
		_, err := s.customerRepo.FindByID(ctx, *req.CustomerID)
		if err != nil {
			return nil, errors.NotFoundWithID("Customer", req.CustomerID.String())
		}
	}

	// Validate items
	if len(req.Items) == 0 {
		return nil, errors.InvalidInput("Sale must have at least one item")
	}

	// Build sale details and validate
	saleDetails := make([]domain.SaleDetail, 0, len(req.Items))

	for _, itemReq := range req.Items {
		// Validate product
		product, err := s.productRepo.FindByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, errors.NotFoundWithID("Product", itemReq.ProductID.String())
		}

		if product.Status != domain.ProductStatusActive {
			return nil, errors.InvalidInput(fmt.Sprintf("Product %s is not active", product.Name))
		}

		// Check inventory availability
		available, err := s.inventoryRepo.CheckAvailability(ctx, itemReq.ProductID, req.WarehouseID, itemReq.Quantity)
		if err != nil {
			return nil, err
		}

		if !available {
			inventory, _ := s.inventoryRepo.GetByProductAndWarehouse(ctx, itemReq.ProductID, req.WarehouseID)
			availableQty := 0.0
			if inventory != nil {
				availableQty = inventory.AvailableQuantity
			}
			return nil, errors.InsufficientStock(product.Name, availableQty, itemReq.Quantity)
		}

		// Determine unit price
		unitPrice := product.SellingPrice
		if itemReq.UnitPrice != nil {
			unitPrice = *itemReq.UnitPrice
		}

		saleDetail := domain.SaleDetail{
			ProductID:      itemReq.ProductID,
			Quantity:       itemReq.Quantity,
			UnitPrice:      unitPrice,
			DiscountAmount: itemReq.DiscountAmount,
		}

		saleDetails = append(saleDetails, saleDetail)
	}

	// Create sale
	sale := &domain.Sale{
		SaleID:           uuid.New(),
		CustomerID:       req.CustomerID,
		StoreID:          &req.StoreID,
		WarehouseID:      &req.WarehouseID,
		SaleType:         req.SaleType,
		Status:           domain.SaleStatusCompleted,
		Currency:         req.Currency,
		ExchangeRate:     req.ExchangeRate,
		DiscountAmount:   req.DiscountAmount,
		PaymentMethod:    req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
		SalespersonID:    &req.SalespersonID,
	}

	// Create sale with details (transaction handled in repository)
	if err := s.saleRepo.CreateWithDetails(ctx, sale, saleDetails); err != nil {
		return nil, err
	}

	// Reload with details
	return s.saleRepo.FindByID(ctx, sale.SaleID)
}

// CreateCreditSale creates a credit sale with accounts receivable
func (s *saleService) CreateCreditSale(ctx context.Context, req services.CreateSaleRequest, creditDays int) (*domain.Sale, *domain.AccountsReceivable, error) {
	// Customer is required for credit sales
	if req.CustomerID == nil {
		return nil, nil, errors.InvalidInput("Customer is required for credit sales")
	}

	// Validate customer
	_, err := s.customerRepo.FindByID(ctx, *req.CustomerID)
	if err != nil {
		return nil, nil, errors.NotFoundWithID("Customer", req.CustomerID.String())
	}

	// Set sale type to credit
	req.SaleType = domain.SaleTypeCredit

	// Create the sale
	sale, err := s.CreateSale(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	// Calculate due date
	dueDate := time.Now().AddDate(0, 0, creditDays)

	// Create accounts receivable
	ar := &domain.AccountsReceivable{
		ReceivableID:  uuid.New(),
		SaleID:        &sale.SaleID,
		CustomerID:    *req.CustomerID,
		TotalAmount:   sale.TotalAmount,
		PaidAmount:    0,
		Balance:      sale.TotalAmount,
		Currency:      sale.Currency,
		DueDate:       dueDate,
		Status:        domain.AccountStatusPending,
	}

	// Create AR record directly
	if err := s.db.WithContext(ctx).Create(ar).Error; err != nil {
		return sale, nil, errors.WrapError(err, "failed to create accounts receivable")
	}

	return sale, ar, nil
}

// GetSale retrieves a sale by ID
func (s *saleService) GetSale(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	return s.saleRepo.FindByID(ctx, id)
}

// GetSaleByInvoiceNumber retrieves a sale by invoice number
func (s *saleService) GetSaleByInvoiceNumber(ctx context.Context, invoiceNumber string) (*domain.Sale, error) {
	return s.saleRepo.FindByInvoiceNumber(ctx, invoiceNumber)
}

// ListSales lists sales with filters
func (s *saleService) ListSales(ctx context.Context, filters repositories.SaleFilters, limit, offset int) ([]domain.Sale, int64, error) {
	return s.saleRepo.List(ctx, filters, limit, offset)
}

// CancelSale cancels a sale and reverses inventory
func (s *saleService) CancelSale(ctx context.Context, id uuid.UUID, reason string) error {
	sale, err := s.saleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if sale.Status == domain.SaleStatusCancelled {
		return errors.InvalidInput("Sale is already cancelled")
	}

	if sale.Status != domain.SaleStatusCompleted {
		return errors.InvalidInput(fmt.Sprintf("Cannot cancel sale with status %s", sale.Status))
	}

	// Cancel sale (reverses inventory via repository)
	if err := s.saleRepo.Cancel(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetDailySales retrieves sales for a specific date
func (s *saleService) GetDailySales(ctx context.Context, storeID *uuid.UUID, date time.Time) ([]domain.Sale, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filters := repositories.SaleFilters{
		StoreID:   storeID,
		DateFrom: &startOfDay,
		DateTo:   &endOfDay,
	}

	sales, _, err := s.saleRepo.List(ctx, filters, 10000, 0)
	return sales, err
}

// GetSalesByPeriod retrieves sales for a period
func (s *saleService) GetSalesByPeriod(ctx context.Context, from, to time.Time) ([]domain.Sale, error) {
	filters := repositories.SaleFilters{
		DateFrom: &from,
		DateTo:   &to,
	}

	sales, _, err := s.saleRepo.List(ctx, filters, 10000, 0)
	return sales, err
}
