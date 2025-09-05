package domain

import (
	"github.com/google/uuid"
	"time"
)

// Subscription represents the core business entity for user server
type Subscription struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EndDate     *string // Format: MM-YYYY, nullable
	ServiceName string
	StartDate   string // Format: MM-YYYY
	Price       int
	ID          uuid.UUID
	UserID      uuid.UUID
}

// NewSubscription creates a new Subscription with validation
func NewSubscription(id uuid.UUID, serviceName string, price int, userID uuid.UUID, startDate string, endDate *string) (*Subscription, error) {
	sub := &Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := sub.Validate(); err != nil {
		return nil, err
	}

	return sub, nil
}

// Validate validates the subscription business rules
func (s *Subscription) Validate() error {
	if s.ServiceName == "" {
		return NewValidationError("service name", "service name is required")
	}

	if s.Price <= 0 {
		return NewValidationError("price", "price must be positive")
	}

	if s.UserID == uuid.Nil {
		return NewValidationError("user ID", "user ID is required")
	}

	if err := ValidateSubscriptionDates(s.StartDate, s.EndDate); err != nil {
		return err
	}

	return nil
}

// IsActive checks if the subscription is currently active based on the provided date
func (s *Subscription) IsActive(referenceDate string) (bool, error) {
	if err := ValidateDateFormat(referenceDate); err != nil {
		return false, err
	}

	// Subscription is active if the reference date is between start and end date
	// or if there's no end date and reference date is after start date
	refAfterStart, err := isDateAfterOrEqual(referenceDate, s.StartDate)
	if err != nil {
		return false, err
	}

	if s.EndDate == nil {
		return refAfterStart, nil
	}

	refBeforeEnd, err := isDateBeforeOrEqual(referenceDate, *s.EndDate)
	if err != nil {
		return false, err
	}

	return refAfterStart && refBeforeEnd, nil
}

// SubscriptionFactory creates server with generated ID
type SubscriptionFactory struct{}

func (f *SubscriptionFactory) CreateSubscription(serviceName string, price int, userID uuid.UUID, startDate string, endDate *string) (*Subscription, error) {
	id := uuid.New()

	sub, err := NewSubscription(id, serviceName, price, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
