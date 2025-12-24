package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/core/ports/services"
)

func getCustomerName(customer *domain.Customer) string {
	if customer.BusinessName != nil && *customer.BusinessName != "" {
		return *customer.BusinessName
	}
	name := ""
	if customer.FirstName != nil {
		name = *customer.FirstName
	}
	if customer.LastName != nil {
		if name != "" {
			name += " "
		}
		name += *customer.LastName
	}
	if name == "" {
		return "Cliente"
	}
	return name
}

type notificationService struct {
	reservationRepo repositories.ReservationRepository
	customerRepo    repositories.CustomerRepository
	db              *gorm.DB
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	reservationRepo repositories.ReservationRepository,
	customerRepo repositories.CustomerRepository,
	db *gorm.DB,
) services.NotificationService {
	return &notificationService{
		reservationRepo: reservationRepo,
		customerRepo:    customerRepo,
		db:              db,
	}
}

// SendReservationConfirmation sends confirmation when reservation is created
func (s *notificationService) SendReservationConfirmation(ctx context.Context, reservationID uuid.UUID) error {
	reservation, err := s.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return err
	}

	customer, err := s.customerRepo.FindByID(ctx, reservation.CustomerID)
	if err != nil {
		return err
	}

	// Build notification message
	message := fmt.Sprintf(
		"Estimado/a %s,\n\n"+
			"Su reserva %s ha sido creada exitosamente.\n\n"+
			"Detalles:\n"+
			"- Monto Total: %.2f %s\n"+
			"- Depósito: %.2f %s\n"+
			"- Saldo: %.2f %s\n"+
			"- Fecha de Vencimiento: %s\n\n"+
			"Gracias por su preferencia,\n"+
			"Bazaar Araira",
		getCustomerName(customer),
		reservation.ReservationNumber,
		reservation.TotalAmount,
		reservation.Currency,
		reservation.DepositAmount,
		reservation.Currency,
		reservation.Balance,
		reservation.Currency,
		reservation.ExpirationDate.Format("02/01/2006"),
	)

	// Create notification record
	notification := &domain.CustomerNotification{
		NotificationID:   uuid.New(),
		CustomerID:       customer.CustomerID,
		NotificationType: domain.NotificationTypeEmail,
		Status:           domain.NotificationStatusPending,
		Subject:          stringPtr("Confirmación de Reserva"),
		Message:          message,
		ReferenceType:    stringPtr("RESERVATION"),
		ReferenceID:      &reservationID,
		ScheduledAt:      timePtr(time.Now()),
	}

	if err := s.db.WithContext(ctx).Create(notification).Error; err != nil {
		return errors.WrapError(err, "failed to create notification")
	}

	// Log notification (in production, send via email/SMS)
	log.Printf("[NOTIFICATION] Reservation Confirmation - Customer: %s, Reservation: %s",
		getCustomerName(customer),
		reservation.ReservationNumber,
	)

	return nil
}

// SendReservationReminder sends reminder for expiring reservation
func (s *notificationService) SendReservationReminder(ctx context.Context, reservationID uuid.UUID) error {
	reservation, err := s.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return err
	}

	customer, err := s.customerRepo.FindByID(ctx, reservation.CustomerID)
	if err != nil {
		return err
	}

	hoursRemaining := time.Until(reservation.ExpirationDate).Hours()

	message := fmt.Sprintf(
		"Estimado/a %s,\n\n"+
			"RECORDATORIO: Su reserva %s vencerá pronto.\n\n"+
			"Tiempo restante: aproximadamente %.0f horas\n"+
			"Fecha de vencimiento: %s\n\n"+
			"Detalles:\n"+
			"- Monto Total: %.2f %s\n"+
			"- Saldo Pendiente: %.2f %s\n\n"+
			"Por favor, complete el pago antes de la fecha de vencimiento.\n\n"+
			"Gracias,\n"+
			"Bazaar Araira",
		getCustomerName(customer),
		reservation.ReservationNumber,
		hoursRemaining,
		reservation.ExpirationDate.Format("02/01/2006 15:04"),
		reservation.TotalAmount,
		reservation.Currency,
		reservation.Balance,
		reservation.Currency,
	)

	// Create notification record
	notification := &domain.CustomerNotification{
		NotificationID:   uuid.New(),
		CustomerID:       customer.CustomerID,
		NotificationType: domain.NotificationTypeEmail,
		Status:           domain.NotificationStatusPending,
		Subject:          stringPtr("Recordatorio de Reserva"),
		Message:          message,
		ReferenceType:    stringPtr("RESERVATION"),
		ReferenceID:      &reservationID,
		ScheduledAt:      timePtr(time.Now()),
	}

	if err := s.db.WithContext(ctx).Create(notification).Error; err != nil {
		return errors.WrapError(err, "failed to create notification")
	}

	log.Printf("[NOTIFICATION] Reservation Reminder - Customer: %s, Reservation: %s, Hours Remaining: %.0f",
		getCustomerName(customer),
		reservation.ReservationNumber,
		hoursRemaining,
	)

	return nil
}

// SendPreOrderReady sends notification when pre-order is ready
func (s *notificationService) SendPreOrderReady(ctx context.Context, preOrderID uuid.UUID) error {
	var preOrder domain.PreOrder
	if err := s.db.WithContext(ctx).First(&preOrder, "pre_order_id = ?", preOrderID).Error; err != nil {
		return errors.NotFoundWithID("PreOrder", preOrderID.String())
	}

	customer, err := s.customerRepo.FindByID(ctx, preOrder.CustomerID)
	if err != nil {
		return err
	}

	message := fmt.Sprintf(
		"Estimado/a %s,\n\n"+
			"¡Buenas noticias! Su pre-orden %s está lista para retirar.\n\n"+
			"Detalles:\n"+
			"- Monto Total: %.2f %s\n"+
			"- Depósito Pagado: %.2f %s\n"+
			"- Saldo: %.2f %s\n\n"+
			"Por favor, pase por nuestra tienda para completar su compra.\n\n"+
			"Gracias,\n"+
			"Bazaar Araira",
		getCustomerName(customer),
		preOrder.PreOrderNumber,
		preOrder.TotalAmount,
		preOrder.Currency,
		preOrder.DepositPaid,
		preOrder.Currency,
		preOrder.TotalAmount-preOrder.DepositPaid,
		preOrder.Currency,
	)

	// Create notification record
	notification := &domain.CustomerNotification{
		NotificationID:   uuid.New(),
		CustomerID:       customer.CustomerID,
		NotificationType: domain.NotificationTypeEmail,
		Status:           domain.NotificationStatusPending,
		Subject:          stringPtr("Pre-Orden Lista"),
		Message:          message,
		ReferenceType:    stringPtr("PRE_ORDER"),
		ReferenceID:      &preOrderID,
		ScheduledAt:      timePtr(time.Now()),
	}

	if err := s.db.WithContext(ctx).Create(notification).Error; err != nil {
		return errors.WrapError(err, "failed to create notification")
	}

	log.Printf("[NOTIFICATION] Pre-Order Ready - Customer: %s, PreOrder: %s",
		getCustomerName(customer),
		preOrder.PreOrderNumber,
	)

	return nil
}

// SendCustomNotification sends a custom notification
func (s *notificationService) SendCustomNotification(
	ctx context.Context,
	customerID uuid.UUID,
	notificationType domain.NotificationType,
	subject, message string,
) error {
	customer, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return errors.NotFoundWithID("Customer", customerID.String())
	}

	// Create notification record
	notification := &domain.CustomerNotification{
		NotificationID:   uuid.New(),
		CustomerID:       customerID,
		NotificationType: notificationType,
		Status:           domain.NotificationStatusPending,
		Subject:          &subject,
		Message:          message,
		ScheduledAt:      timePtr(time.Now()),
	}

	if err := s.db.WithContext(ctx).Create(notification).Error; err != nil {
		return errors.WrapError(err, "failed to create notification")
	}

	log.Printf("[NOTIFICATION] Custom - Customer: %s, Type: %s, Subject: %s",
		getCustomerName(customer),
		notificationType,
		subject,
	)

	return nil
}

// ProcessPendingNotifications processes pending notifications
func (s *notificationService) ProcessPendingNotifications(ctx context.Context) (int, error) {
	var notifications []domain.CustomerNotification

	// Get pending notifications
	err := s.db.WithContext(ctx).
		Where("status = ?", domain.NotificationStatusPending).
		Where("scheduled_at <= ?", time.Now()).
		Limit(100).
		Find(&notifications).Error

	if err != nil {
		return 0, errors.WrapError(err, "failed to fetch pending notifications")
	}

	count := 0
	for _, notification := range notifications {
		// TODO: Implement actual sending logic (email, SMS, push)
		// For now, just mark as sent

		now := time.Now()
		notification.Status = domain.NotificationStatusSent
		notification.SentAt = &now

		if err := s.db.WithContext(ctx).Save(&notification).Error; err != nil {
			log.Printf("[ERROR] Failed to update notification %s: %v", notification.NotificationID, err)
			continue
		}

		count++
	}

	return count, nil
}
