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

type reservationService struct {
	reservationRepo repositories.ReservationRepository
	customerRepo    repositories.CustomerRepository
	productRepo     repositories.ProductRepository
	inventoryRepo   repositories.InventoryRepository
	saleRepo        repositories.SaleRepository
	notificationSvc services.NotificationService
	db              *gorm.DB
}

// NewReservationService creates a new reservation service
func NewReservationService(
	reservationRepo repositories.ReservationRepository,
	customerRepo repositories.CustomerRepository,
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
	saleRepo repositories.SaleRepository,
	notificationSvc services.NotificationService,
	db *gorm.DB,
) services.ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		customerRepo:    customerRepo,
		productRepo:     productRepo,
		inventoryRepo:   inventoryRepo,
		saleRepo:        saleRepo,
		notificationSvc: notificationSvc,
		db:              db,
	}
}

// CreateReservation creates a new reservation with validation
func (s *reservationService) CreateReservation(ctx context.Context, req services.CreateReservationRequest) (*domain.Reservation, error) {
	// Validate customer exists
	_, err := s.customerRepo.FindByID(ctx, req.CustomerID)
	if err != nil {
		return nil, errors.NotFoundWithID("Customer", req.CustomerID.String())
	}

	// Validate child if provided
	if req.ChildID != nil {
		child, err := s.customerRepo.FindByID(ctx, *req.ChildID)
		if err != nil {
			return nil, errors.NotFoundWithID("Child", req.ChildID.String())
		}
		// Verify child belongs to customer
		if child.CustomerID != req.CustomerID {
			return nil, errors.InvalidInput("Child does not belong to the specified customer")
		}
	}

	// Validate items
	if len(req.Items) == 0 {
		return nil, errors.InvalidInput("Reservation must have at least one item")
	}

	// Build reservation items
	reservationItems := make([]domain.ReservationItem, 0, len(req.Items))
	totalAmount := 0.0

	for _, itemReq := range req.Items {
		// Validate product exists
		product, err := s.productRepo.FindByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, errors.NotFoundWithID("Product", itemReq.ProductID.String())
		}

		if product.Status != domain.ProductStatusActive {
			return nil, errors.InvalidInput(fmt.Sprintf("Product %s is not active", product.Name))
		}

		// Use current sale price
		unitPrice := product.SellingPrice
		itemTotal := unitPrice * itemReq.Quantity

		reservationItem := domain.ReservationItem{
			ReservationItemID: uuid.New(),
			ProductID:         itemReq.ProductID,
			Quantity:          itemReq.Quantity,
			UnitPrice:         unitPrice,
			TotalAmount:       itemTotal,
		}

		reservationItems = append(reservationItems, reservationItem)
		totalAmount += itemTotal
	}

	// Calculate expiration date
	expirationDate := time.Now().AddDate(0, 0, req.ExpirationDays)

	// Validate deposit amount
	if req.DepositAmount < 0 {
		return nil, errors.InvalidInput("Deposit amount cannot be negative")
	}

	if req.DepositAmount > totalAmount {
		return nil, errors.InvalidInput("Deposit amount cannot exceed total amount")
	}

	// Create reservation
	reservation := &domain.Reservation{
		ReservationID:   uuid.New(),
		CustomerID:      req.CustomerID,
		ChildID:         req.ChildID,
		ListID:          req.ListID,
		StoreID:         &req.StoreID,
		Status:          domain.ReservationStatusPending,
		ReservationDate: time.Now(),
		ExpirationDate:  expirationDate,
		TotalAmount:     totalAmount,
		DepositAmount:   req.DepositAmount,
		Balance:         totalAmount - req.DepositAmount,
		Currency:        req.Currency,
		Notes:           req.Notes,
		CreatedBy:       &req.UserID,
	}

	// Create reservation with items (transaction handled in repository)
	if err := s.reservationRepo.CreateWithItems(ctx, reservation, reservationItems); err != nil {
		return nil, err
	}

	// Send confirmation notification
	go func() {
		_ = s.notificationSvc.SendReservationConfirmation(context.Background(), reservation.ReservationID)
	}()

	// Reload with all relations
	return s.reservationRepo.FindByID(ctx, reservation.ReservationID)
}

// GetReservation retrieves a reservation by ID
func (s *reservationService) GetReservation(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	return s.reservationRepo.FindByID(ctx, id)
}

// GetReservationByNumber retrieves a reservation by number
func (s *reservationService) GetReservationByNumber(ctx context.Context, reservationNumber string) (*domain.Reservation, error) {
	return s.reservationRepo.FindByNumber(ctx, reservationNumber)
}

// ListReservations lists reservations with filters
func (s *reservationService) ListReservations(ctx context.Context, filters repositories.ReservationFilters, limit, offset int) ([]domain.Reservation, int64, error) {
	return s.reservationRepo.List(ctx, filters, limit, offset)
}

// ConfirmReservation confirms a reservation (marks deposit as paid)
func (s *reservationService) ConfirmReservation(ctx context.Context, id uuid.UUID) error {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.Status != domain.ReservationStatusPending {
		return errors.InvalidInput(fmt.Sprintf("Cannot confirm reservation with status %s", reservation.Status))
	}

	// Update status to confirmed
	reservation.Status = domain.ReservationStatusConfirmed
	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		return err
	}

	return nil
}

// FulfillReservation converts a reservation into a sale
func (s *reservationService) FulfillReservation(ctx context.Context, req services.FulfillReservationRequest) (*domain.Sale, error) {
	// Get reservation with items
	reservation, err := s.reservationRepo.FindByID(ctx, req.ReservationID)
	if err != nil {
		return nil, err
	}

	if reservation.Status != domain.ReservationStatusConfirmed {
		return nil, errors.InvalidInput(fmt.Sprintf("Cannot fulfill reservation with status %s. Must be CONFIRMED", reservation.Status))
	}

	// Get reservation items
	items, err := s.reservationRepo.GetItems(ctx, req.ReservationID)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.InvalidInput("Reservation has no items")
	}

	// Get warehouse from store
	var warehouse domain.Warehouse
	if err := s.db.Where("store_id = ?", reservation.StoreID).First(&warehouse).Error; err != nil {
		return nil, errors.WrapError(err, "failed to find warehouse for store")
	}

	// Create sale from reservation
	sale := &domain.Sale{
		SaleID:           uuid.New(),
		CustomerID:       &reservation.CustomerID,
		StoreID:          reservation.StoreID,
		WarehouseID:      &warehouse.WarehouseID,
		SaleType:         domain.SaleTypeCash,
		Status:           domain.SaleStatusCompleted,
		Currency:         reservation.Currency,
		ExchangeRate:     req.ExchangeRate,
		PaymentMethod:    &req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		Notes:            stringPtr(fmt.Sprintf("Fulfillment of reservation %s", reservation.ReservationNumber)),
		SalespersonID:    &req.UserID,
	}

	// Convert reservation items to sale items
	saleDetails := make([]domain.SaleDetail, 0, len(items))
	for _, item := range items {
		saleDetail := domain.SaleDetail{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
		}
		saleDetails = append(saleDetails, saleDetail)
	}

	// Create sale with details in transaction
	if err := s.saleRepo.CreateWithDetails(ctx, sale, saleDetails); err != nil {
		return nil, err
	}

	// Mark reservation as fulfilled
	now := time.Now()
	reservation.Status = domain.ReservationStatusFulfilled
	reservation.FulfilledAt = &now
	reservation.FulfilledBy = &req.UserID

	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		// Sale was created but reservation update failed - log this
		return sale, errors.WrapError(err, "sale created but failed to update reservation status")
	}

	return sale, nil
}

// CancelReservation cancels a reservation and releases inventory
func (s *reservationService) CancelReservation(ctx context.Context, id, userID uuid.UUID, reason string) error {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.Status == domain.ReservationStatusFulfilled {
		return errors.InvalidInput("Cannot cancel a fulfilled reservation")
	}

	if reservation.Status == domain.ReservationStatusCancelled {
		return errors.InvalidInput("Reservation is already cancelled")
	}

	// Cancel reservation (releases inventory via repository)
	if err := s.reservationRepo.Cancel(ctx, id); err != nil {
		return err
	}

	// TODO: Record cancellation reason and user

	return nil
}

// ExpireReservations expires all pending/confirmed reservations past their expiration date
func (s *reservationService) ExpireReservations(ctx context.Context) (int, error) {
	expired, err := s.reservationRepo.GetExpired(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, reservation := range expired {
		if err := s.reservationRepo.Cancel(ctx, reservation.ReservationID); err != nil {
			// Log error but continue processing
			continue
		}
		count++
	}

	return count, nil
}

// SendReminders sends reminders for reservations expiring soon
func (s *reservationService) SendReminders(ctx context.Context, hoursBeforeExpiration int) (int, error) {

	reservations, err := s.reservationRepo.GetExpiringFor(ctx, time.Duration(hoursBeforeExpiration)*time.Hour)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, reservation := range reservations {
		// Skip if reminder already sent recently
		if reservation.ReminderSentAt != nil && time.Since(*reservation.ReminderSentAt) < 24*time.Hour {
			continue
		}

		if err := s.notificationSvc.SendReservationReminder(ctx, reservation.ReservationID); err != nil {
			// Log error but continue
			continue
		}

		// Update reminder sent time
		now := time.Now()
		reservation.ReminderSentAt = &now
		if err := s.reservationRepo.Update(ctx, &reservation); err != nil {
			// Log error but continue
			continue
		}

		count++
	}

	return count, nil
}
