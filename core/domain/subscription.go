package domain

import (
	"github.com/google/uuid"
	"time"
)

// Subscription represents the core business entity for user server
type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   string  // Format: MM-YYYY
	EndDate     *string // Format: MM-YYYY, nullable
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

// CalculateCostForPeriod calculates the cost for a given period
func (s *Subscription) CalculateCostForPeriod(startPeriod, endPeriod string) (int, error) {
	if err := ValidateDateFormat(startPeriod); err != nil {
		return 0, err
	}
	if err := ValidateDateFormat(endPeriod); err != nil {
		return 0, err
	}

	// Check if subscription was active during any part of the period
	subStart := s.StartDate
	subEnd := s.EndDate
	if subEnd == nil {
		// If no end date, subscription is ongoing
		subEnd = &endPeriod
	}

	// Find overlapping months between a subscription period and requested period
	overlapMonths, err := calculateOverlapMonths(subStart, *subEnd, startPeriod, endPeriod)
	if err != nil {
		return 0, err
	}

	return overlapMonths * s.Price, nil
}

// SubscriptionFactory creates server with generated ID
type SubscriptionFactory struct{}

func (f *SubscriptionFactory) CreateSubscription(serviceName string, price int, userID uuid.UUID, startDate string, endDate *string) (*Subscription, error) {
	id := GenerateUUID()

	sub, err := NewSubscription(id, serviceName, price, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
