package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ValidationErrors holds multiple validation errors
type ValidationErrors map[string][]string

func (v ValidationErrors) Add(field, message string) {
	if v[field] == nil {
		v[field] = []string{}
	}
	v[field] = append(v[field], message)
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func (v ValidationErrors) Error() string {
	var messages []string
	for field, errs := range v {
		messages = append(messages, fmt.Sprintf("%s: %s", field, strings.Join(errs, ", ")))
	}
	return strings.Join(messages, "; ")
}

// Validator provides validation utilities
type Validator struct {
	errors ValidationErrors
}

func New() *Validator {
	return &Validator{
		errors: make(ValidationErrors),
	}
}

func (v *Validator) IsValid() bool {
	return !v.errors.HasErrors()
}

func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// Required checks if a value is not empty
func (v *Validator) Required(field string, value interface{}) {
	switch val := value.(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			v.errors.Add(field, "is required")
		}
	case *string:
		if val == nil || strings.TrimSpace(*val) == "" {
			v.errors.Add(field, "is required")
		}
	case uuid.UUID:
		if val == uuid.Nil {
			v.errors.Add(field, "is required")
		}
	case *uuid.UUID:
		if val == nil || *val == uuid.Nil {
			v.errors.Add(field, "is required")
		}
	case nil:
		v.errors.Add(field, "is required")
	}
}

// MinLength checks if a string has minimum length
func (v *Validator) MinLength(field, value string, min int) {
	if len(value) < min {
		v.errors.Add(field, fmt.Sprintf("must be at least %d characters", min))
	}
}

// MaxLength checks if a string doesn't exceed maximum length
func (v *Validator) MaxLength(field, value string, max int) {
	if len(value) > max {
		v.errors.Add(field, fmt.Sprintf("must not exceed %d characters", max))
	}
}

// Email validates an email address
func (v *Validator) Email(field, value string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid email address")
	}
}

// Min checks if a number is greater than or equal to minimum
func (v *Validator) Min(field string, value, min float64) {
	if value < min {
		v.errors.Add(field, fmt.Sprintf("must be at least %.2f", min))
	}
}

// Max checks if a number is less than or equal to maximum
func (v *Validator) Max(field string, value, max float64) {
	if value > max {
		v.errors.Add(field, fmt.Sprintf("must not exceed %.2f", max))
	}
}

// Positive checks if a number is positive
func (v *Validator) Positive(field string, value float64) {
	if value <= 0 {
		v.errors.Add(field, "must be positive")
	}
}

// NonNegative checks if a number is non-negative
func (v *Validator) NonNegative(field string, value float64) {
	if value < 0 {
		v.errors.Add(field, "must be non-negative")
	}
}

// UUID validates a UUID string
func (v *Validator) UUID(field, value string) {
	if _, err := uuid.Parse(value); err != nil {
		v.errors.Add(field, "must be a valid UUID")
	}
}

// In checks if a value is in a list of allowed values
func (v *Validator) In(field string, value interface{}, allowed []interface{}) {
	found := false
	for _, a := range allowed {
		if value == a {
			found = true
			break
		}
	}
	if !found {
		v.errors.Add(field, "is not a valid value")
	}
}

// Custom allows adding custom validation errors
func (v *Validator) Custom(field, message string, condition bool) {
	if !condition {
		v.errors.Add(field, message)
	}
}
