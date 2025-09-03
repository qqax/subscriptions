// core/ports/errors.go
package ports

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
)

// DomainError represents a domain-specific error (optional)
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}
