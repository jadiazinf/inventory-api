package dto

import "github.com/gofiber/fiber/v2"

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(code int, message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// SendSuccess sends a success response
func SendSuccess(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(NewSuccessResponse(data, message))
}

// SendError sends an error response
func SendError(c *fiber.Ctx, statusCode int, message string, details interface{}) error {
	return c.Status(statusCode).JSON(NewErrorResponse(statusCode, message, details))
}

// PaginationParams represents common pagination parameters
type PaginationParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// GetPaginationParams extracts and validates pagination parameters
func GetPaginationParams(c *fiber.Ctx) PaginationParams {
	params := PaginationParams{
		Limit:  20,  // default
		Offset: 0,   // default
	}

	if err := c.QueryParser(&params); err == nil {
		if params.Limit <= 0 || params.Limit > 100 {
			params.Limit = 20
		}
		if params.Offset < 0 {
			params.Offset = 0
		}
	}

	return params
}
