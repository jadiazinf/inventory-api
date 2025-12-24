package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/jadiazinf/inventory/internal/adapters/http/dto"
	"github.com/jadiazinf/inventory/internal/common/errors"
)

// HandleServiceError converts service errors to HTTP responses
func HandleServiceError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *errors.AppError:
		return dto.SendError(c, e.StatusCode, e.Message, e.Details)
	default:
		return dto.SendError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}
}

// GetUserID extracts user ID from context (set by auth middleware)
func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	return userID, ok
}

// ParseUUID parses and validates a UUID from request params
func ParseUUID(c *fiber.Ctx, param string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(param))
	if err != nil {
		return uuid.Nil, dto.SendError(c, fiber.StatusBadRequest, "Invalid "+param, err.Error())
	}
	return id, nil
}
