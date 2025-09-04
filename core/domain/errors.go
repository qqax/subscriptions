package domain

type domainErrorCodes int

const (
	_                   domainErrorCodes = iota
	ValidationError                      = 400
	NotFoundError                        = 404
	DuplicateError                       = 422
	InternalServerError                  = 500
)

// Error definitions for the core domain
var (
	ErrSubscriptionNotFound  = NewDomainError(NotFoundError, "subscription not found")
	ErrDuplicateSubscription = NewDomainError(DuplicateError, "DuplicateError subscription")
	ErrInvalidDateformat     = NewDomainError(ValidationError, "invalid date format, expected MM-YYYY")
	ErrInvalidUUID           = NewDomainError(ValidationError, "invalid UUID format")
	ErrInvalidPrice          = NewDomainError(ValidationError, "price must be positive integer")
	ErrStartDateAfterEndDate = NewDomainError(ValidationError, "start date cannot be after end date")
	ErrInvalidDateRange      = NewDomainError(ValidationError, "invalid date range")
	ErrValidationFailed      = NewDomainError(ValidationError, "validation failed")
	ErrInternal              = NewDomainError(InternalServerError, "internal server error")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// NewValidationError creates error with details
func NewValidationError(field, reason string) *DomainError {
	return &DomainError{
		Code:    ValidationError,
		Message: "Validation failed for field: " + field,
		Details: map[string]interface{}{
			"field":  field,
			"reason": reason,
			"type":   "validation",
		},
	}
}

// NewValidationErrorWithDetails creates validation error with additional details
func NewValidationErrorWithDetails(field, reason string, details map[string]interface{}) *DomainError {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["field"] = field
	details["reason"] = reason
	details["type"] = "validation"

	return &DomainError{
		Code:    ValidationError,
		Message: "Validation failed for field: " + field,
		Details: details,
	}
}
