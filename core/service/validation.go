package service

import (
	"github.com/google/uuid"
	"regexp"
	"subscription/core/domain"
	"subscription/core/ports"
)

// validateFilter валидация параметров фильтрации
func validateFilter(filter ports.SubscriptionFilter) error {
	// Валидация UUID в UserIDs
	for _, userID := range filter.UserIDs {
		if userID == uuid.Nil {
			return domain.NewValidationError("user_ids", "contains invalid UUID format")
		}
	}

	// Валидация дат если они указаны
	if filter.StartDateFrom != nil {
		if err := validateDateFormat(*filter.StartDateFrom); err != nil {
			return domain.NewValidationError("start_date_from", "invalid date format, expected MM-YYYY")
		}
	}

	if filter.StartDateTo != nil {
		if err := validateDateFormat(*filter.StartDateTo); err != nil {
			return domain.NewValidationError("start_date_to", "invalid date format, expected MM-YYYY")
		}
	}

	// Валидация диапазона дат если обе даты указаны
	if filter.StartDateFrom != nil && filter.StartDateTo != nil {
		if err := validateDateRange(*filter.StartDateFrom, *filter.StartDateTo); err != nil {
			return domain.NewValidationError("date_range", "start date cannot be after end date")
		}
	}

	return nil
}

func validateDateFormat(date string) error {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-20\d{2}$`, date)
	if !matched {
		return domain.ErrInvalidDateformat
	}
	return nil
}

func validateDateRange(startDate, endDate string) error {
	if startDate > endDate {
		return domain.ErrStartDateAfterEndDate
	}
	return nil
}

func validateSubscriptionDates(startDate string, endDate *string) error {
	if err := validateDateFormat(startDate); err != nil {
		return err
	}
	if endDate != nil {
		if err := validateDateFormat(*endDate); err != nil {
			return err
		}
		if err := validateDateRange(startDate, *endDate); err != nil {
			return err
		}
	}
	return nil
}
