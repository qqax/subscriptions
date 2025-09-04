package domain

import (
	"regexp"
	"strconv"
	"strings"
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

// isDateAfter checks that date1 > date2
func isDateAfter(date1, date2 string) (bool, error) {
	year1, month1, err := ParseDate(date1)
	if err != nil {
		return false, NewValidationError("date", "invalid first date: "+err.Error())
	}

	year2, month2, err := ParseDate(date2)
	if err != nil {
		return false, NewValidationError("date", "invalid second date: "+err.Error())
	}

	if year1 > year2 {
		return true, nil
	}
	if year1 == year2 && month1 > month2 {
		return true, nil
	}
	return false, nil
}

// isDateAfterOrEqual checks that date1 >= date2
func isDateAfterOrEqual(date1, date2 string) (bool, error) {
	after, err := isDateAfter(date1, date2)
	if err != nil {
		return false, err
	}
	if after {
		return true, nil
	}

	// Check equality
	year1, month1, err := ParseDate(date1)
	if err != nil {
		return false, err
	}

	year2, month2, err := ParseDate(date2)
	if err != nil {
		return false, err
	}

	return year1 == year2 && month1 == month2, nil
}

// isDateBeforeOrEqual checks that date1 <= date2
func isDateBeforeOrEqual(date1, date2 string) (bool, error) {
	after, err := isDateAfter(date1, date2)
	if err != nil {
		return false, err
	}
	return !after, nil
}

// ParseDate parses string MM-YYYY into year and month
func ParseDate(date string) (year, month int, err error) {
	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return 0, 0, NewValidationError("date_format", "invalid format, expected MM-YYYY")
	}

	month, e := strconv.Atoi(parts[0])
	if e != nil || month < 1 || month > 12 {
		return 0, 0, NewValidationError("date_month", "month must be between 01 and 12")
	}

	year, e = strconv.Atoi(parts[1])
	if e != nil || year < 2000 || year > 2100 {
		return 0, 0, NewValidationError("date_year", "year must be between 2000 and 2100")
	}

	return year, month, nil
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
