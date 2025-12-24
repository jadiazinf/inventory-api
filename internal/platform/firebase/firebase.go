package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jadiazinf/inventory/internal/config"
	"google.golang.org/api/option"
)

var (
	firebaseApp    *firebase.App
	firebaseAuth   *auth.Client
)

// InitFirebase initializes Firebase app and auth client
func InitFirebase(cfg *config.Config) (*firebase.App, *auth.Client, error) {
	opt := option.WithCredentialsFile(cfg.FirebaseCred)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing firebase auth: %v", err)
	}

	firebaseApp = app
	firebaseAuth = authClient

	return app, authClient, nil
}

// NewFirebaseApp creates a new Firebase app (legacy function for compatibility)
func NewFirebaseApp(cfg *config.Config) (*firebase.App, error) {
	app, _, err := InitFirebase(cfg)
	return app, err
}

// GetAuthClient returns the Firebase Auth client
func GetAuthClient() *auth.Client {
	return firebaseAuth
}

// VerifyIDToken verifies a Firebase ID token
func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if firebaseAuth == nil {
		return nil, fmt.Errorf("firebase auth not initialized")
	}
	return firebaseAuth.VerifyIDToken(ctx, idToken)
}
