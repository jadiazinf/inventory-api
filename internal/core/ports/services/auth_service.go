package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	VerifyFirebaseToken(ctx context.Context, idToken string) (*FirebaseUser, error)
	GetOrCreateUser(ctx context.Context, firebaseUID string) (*domain.User, error)
	HasPermission(ctx context.Context, userID uuid.UUID, module, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error)
}

// FirebaseUser represents user information from Firebase
type FirebaseUser struct {
	UID         string
	Email       string
	DisplayName string
	PhotoURL    string
	PhoneNumber string
}
