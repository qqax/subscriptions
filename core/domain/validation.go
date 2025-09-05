package domain

import (
	"regexp"
)

// ValidateDateFormat checks date format MM-YYYY
func ValidateDateFormat(date string) error {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-20\d{2}$`, date)
	if !matched {
		return NewValidationError("date", "must be in MM-YYYY format (e.g., 12-2024)")
	}
	return nil
}

// ValidateDateRange checks that startDate <= endDate
func ValidateDateRange(startDate, endDate string) error {
	startAfterEnd, err := isDateAfter(startDate, endDate)
	if err != nil {
		return NewValidationError("date_range", "invalid date comparison: "+err.Error())
	}
	if startAfterEnd {
		return ErrStartDateAfterEndDate
	}
	return nil
}

// ValidateSubscriptionDates validates startDate and endDate
func ValidateSubscriptionDates(startDate string, endDate *string) error {
	if err := ValidateDateFormat(startDate); err != nil {
		return err
	}

	if endDate != nil {
		if err := ValidateDateFormat(*endDate); err != nil {
			return err
		}
		if err := ValidateDateRange(startDate, *endDate); err != nil {
			return err
		}
	}

	return nil
}
