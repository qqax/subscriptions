package domain

import "errors"

// Error definitions for the core domain
var (
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrInvalidDateformat     = errors.New("invalid date format, expected MM-YYYY")
	ErrInvalidUUID           = errors.New("invalid UUID format")
	ErrInvalidPrice          = errors.New("price must be positive integer")
	ErrStartDateAfterEndDate = errors.New("start date cannot be after end date")
	ErrInvalidDateRange      = errors.New("invalid date range")
	ErrDuplicateSubscription = errors.New("duplicate subscription")
	ErrValidationFailed      = errors.New("validation failed")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// NewValidationError создает ошибку валидации с деталями
func NewValidationError(field, reason string) *DomainError {
	return &DomainError{
		Code:    "validation_error",
		Message: "Validation failed for field: " + field,
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
			"type":   "validation",
		},
	}
}

// NewValidationErrorWithDetails создает ошибку валидации с дополнительными деталями
func NewValidationErrorWithDetails(field, reason string, details map[string]interface{}) *DomainError {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["field"] = field
	details["reason"] = reason
	details["type"] = "validation"

	return &DomainError{
		Code:    "validation_error",
		Message: "Validation failed for field: " + field,
		Details: details,
	}
}

// NewNotFoundError создает ошибку "не найдено"
func NewNotFoundError(resourceType, resourceID string) *DomainError {
	return &DomainError{
		Code:    "not_found",
		Message: resourceType + " not found",
		Details: map[string]interface{}{
			"resource_type": resourceType,
			"resource_id":   resourceID,
			"type":          "not_found",
		},
	}
}

// NewDuplicateError создает ошибку дубликата
func NewDuplicateError(resourceType, field, value string) *DomainError {
	return &DomainError{
		Code:    "duplicate",
		Message: resourceType + " with " + field + " '" + value + "' already exists",
		Details: map[string]interface{}{
			"resource_type": resourceType,
			"field":         field,
			"value":         value,
			"type":          "duplicate",
		},
	}
}
