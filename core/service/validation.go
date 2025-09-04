package service

import (
	"regexp"
	"subscription/core/domain"
)

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
