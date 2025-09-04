package service

import (
	"github.com/google/uuid"
	"subscription/core/domain"
	"subscription/core/ports"
)

// validateFilter validates filter parameters
func validateFilter(filter ports.SubscriptionFilter) error {
	for _, userID := range filter.UserIDs {
		if userID == uuid.Nil {
			return domain.NewValidationError("user_ids", "contains invalid UUID format")
		}
	}

	if filter.StartDateFrom != nil {
		if err := domain.ValidateDateFormat(*filter.StartDateFrom); err != nil {
			return domain.NewValidationError("start_date_from", "invalid date format, expected MM-YYYY")
		}
	}

	if filter.StartDateTo != nil {
		if err := domain.ValidateDateFormat(*filter.StartDateTo); err != nil {
			return domain.NewValidationError("start_date_to", "invalid date format, expected MM-YYYY")
		}
	}

	if filter.StartDateFrom != nil && filter.StartDateTo != nil {
		if err := domain.ValidateDateRange(*filter.StartDateFrom, *filter.StartDateTo); err != nil {
			return domain.NewValidationError("date_range", "start date cannot be after end date")
		}
	}

	return nil
}
