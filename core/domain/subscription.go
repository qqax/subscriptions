package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Subscription represents the core business entity for user server
type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   string  // Format: MM-YYYY
	EndDate     *string // Format: MM-YYYY, nullable
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewSubscription creates a new Subscription with validation
func NewSubscription(serviceName string, price int, userID, startDate string, endDate *string) (*Subscription, error) {
	sub := &Subscription{
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
		return errors.New("service name is required")
	}

	if s.Price <= 0 {
		return errors.New("price must be positive")
	}

	if s.UserID == "" {
		return errors.New("user ID is required")
	}

	if err := validateDateFormat(s.StartDate); err != nil {
		return fmt.Errorf("invalid start date: %w", err)
	}

	if s.EndDate != nil {
		if err := validateDateFormat(*s.EndDate); err != nil {
			return fmt.Errorf("invalid end date: %w", err)
		}
	}

	if s.EndDate != nil {
		if err := validateDateRange(s.StartDate, *s.EndDate); err != nil {
			return err
		}
	}

	return nil
}

// IsActive checks if the subscription is currently active based on provided date
func (s *Subscription) IsActive(referenceDate string) (bool, error) {
	if err := validateDateFormat(referenceDate); err != nil {
		return false, err
	}

	// Subscription is active if reference date is between start and end date
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
	if err := validateDateFormat(startPeriod); err != nil {
		return 0, err
	}
	if err := validateDateFormat(endPeriod); err != nil {
		return 0, err
	}

	// Check if subscription was active during any part of the period
	subStart := s.StartDate
	subEnd := s.EndDate
	if subEnd == nil {
		// If no end date, subscription is ongoing
		subEnd = &endPeriod
	}

	// Find overlapping months between subscription period and requested period
	overlapMonths, err := calculateOverlapMonths(subStart, *subEnd, startPeriod, endPeriod)
	if err != nil {
		return 0, err
	}

	return overlapMonths * s.Price, nil
}

// Helper functions
func validateDateFormat(date string) error {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-20\d{2}$`, date)
	if !matched {
		return errors.New("date must be in MM-YYYY format (e.g., 12-2024)")
	}
	return nil
}

func validateDateRange(startDate, endDate string) error {
	startAfterEnd, err := isDateAfter(startDate, endDate)
	if err != nil {
		return err
	}
	if startAfterEnd {
		return errors.New("start date cannot be after end date")
	}
	return nil
}

func isDateAfter(date1, date2 string) (bool, error) {
	// Convert MM-YYYY to comparable format (YYYYMM)
	// Implementation depends on date comparison logic
	return false, nil
}

func isDateAfterOrEqual(date1, date2 string) (bool, error) {
	// Implementation for date comparison
	return false, nil
}

func isDateBeforeOrEqual(date1, date2 string) (bool, error) {
	// Implementation for date comparison
	return false, nil
}

func calculateOverlapMonths(subStart, subEnd, periodStart, periodEnd string) (int, error) {
	// Calculate number of overlapping months between two periods
	// This is a simplified implementation
	return 1, nil
}

// SubscriptionFactory creates server with generated ID
type SubscriptionFactory struct{}

func (f *SubscriptionFactory) CreateSubscription(serviceName string, price int, userID, startDate string, endDate *string) (*Subscription, error) {
	sub, err := NewSubscription(serviceName, price, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	sub.ID = generateUUID()
	return sub, nil
}

func generateUUID() string {
	// This would typically use github.com/google/uuid
	// For now, return a placeholder
	return "uuid-placeholder"
}
