package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/common/errors"
	"github.com/jadiazinf/inventory/internal/core/domain"
	"github.com/jadiazinf/inventory/internal/core/ports/repositories"
	"github.com/jadiazinf/inventory/internal/platform/firebase"
)

// UserContextKey is the key for storing user in context
const UserContextKey = "user"

// AuthMiddleware creates a middleware that authenticates requests using Firebase
type AuthMiddleware struct {
	userRepo repositories.UserRepository
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
	}
}

// Authenticate is the middleware function
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.Unauthorized("Missing authorization header"))
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.Unauthorized("Invalid authorization header format"))
		}

		idToken := parts[1]

		// Verify token with Firebase
		token, err := firebase.VerifyIDToken(c.Context(), idToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.Unauthorized("Invalid or expired token"))
		}

		// Get or create user in database using firebase_uid
		user, err := m.userRepo.FindByFirebaseUID(c.Context(), token.UID)
		if err != nil {
			// If user doesn't exist, create a minimal user record
			// In production, you might want to require user registration first
			user = &domain.User{
				UserID:      uuid.New(),
				FirebaseUID: token.UID,
				Email:       token.Claims["email"].(string),
				FirstName:   getStringClaim(token.Claims, "name", "User"),
				LastName:    "",
				Status:      domain.UserStatusActive,
			}

			if err := m.userRepo.Create(c.Context(), user); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.InternalError("Failed to create user"))
			}
		}

		// Check if user is active
		if user.Status != domain.UserStatusActive {
			return c.Status(fiber.StatusForbidden).JSON(errors.Forbidden("User account is not active"))
		}

		// Update last login
		_ = m.userRepo.UpdateLastLogin(c.Context(), user.UserID)

		// Store user in context
		c.Locals(UserContextKey, user)

		return c.Next()
	}
}

// RequirePermission creates a middleware that checks for specific permission
func (m *AuthMiddleware) RequirePermission(module, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.Unauthorized("User not authenticated"))
		}

		// TODO: Implement permission checking with permission repository
		// For now, we'll allow all authenticated users
		// In production, implement proper RBAC:
		// hasPermission, err := m.permissionRepo.HasPermission(c.Context(), user.UserID, module, action)
		// if err != nil || !hasPermission {
		//     return c.Status(fiber.StatusForbidden).JSON(errors.Forbidden("Insufficient permissions"))
		// }

		return c.Next()
	}
}

// Optional makes authentication optional but still extracts user if present
func (m *AuthMiddleware) Optional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		idToken := parts[1]
		token, err := firebase.VerifyIDToken(c.Context(), idToken)
		if err != nil {
			return c.Next()
		}

		user, err := m.userRepo.FindByFirebaseUID(c.Context(), token.UID)
		if err == nil {
			c.Locals(UserContextKey, user)
			_ = m.userRepo.UpdateLastLogin(c.Context(), user.UserID)
		}

		return c.Next()
	}
}

// GetUserFromContext retrieves the authenticated user from context
func GetUserFromContext(c *fiber.Ctx) *domain.User {
	user, ok := c.Locals(UserContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}

// getStringClaim safely extracts a string claim from Firebase token
func getStringClaim(claims map[string]interface{}, key, defaultValue string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}
