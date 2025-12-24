package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Details    any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Common error codes
const (
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeAlreadyExists     = "ALREADY_EXISTS"
	ErrCodeInvalidInput      = "INVALID_INPUT"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeInternalError     = "INTERNAL_ERROR"
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInsufficientStock = "INSUFFICIENT_STOCK"
	ErrCodeExpired           = "EXPIRED"
)

// Constructor functions for common errors

func NotFound(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func NotFoundWithID(resource, id string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s with id %s not found", resource, id),
		StatusCode: http.StatusNotFound,
	}
}

func AlreadyExists(resource, field, value string) *AppError {
	return &AppError{
		Code:       ErrCodeAlreadyExists,
		Message:    fmt.Sprintf("%s with %s '%s' already exists", resource, field, value),
		StatusCode: http.StatusConflict,
	}
}

func InvalidInput(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func InvalidInputWithDetails(message string, details any) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

func Unauthorized(message string) *AppError {
	if message == "" {
		message = "Unauthorized"
	}
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func Forbidden(message string) *AppError {
	if message == "" {
		message = "Forbidden"
	}
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func InternalError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInternalError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequest(message string) *AppError{
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func InsufficientStock(productName string, available, requested float64) *AppError {
	return &AppError{
		Code:       ErrCodeInsufficientStock,
		Message:    fmt.Sprintf("Insufficient stock for %s. Available: %.2f, Requested: %.2f", productName, available, requested),
		StatusCode: http.StatusConflict,
		Details: map[string]interface{}{
			"available": available,
			"requested": requested,
		},
	}
}

func Expired(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeExpired,
		Message:    fmt.Sprintf("%s has expired", resource),
		StatusCode: http.StatusGone,
	}
}

// WrapError wraps a generic error into an AppError
func WrapError(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return &AppError{
		Code:       ErrCodeInternalError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Details:    err.Error(),
	}
}
