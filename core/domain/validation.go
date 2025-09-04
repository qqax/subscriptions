// core/domain/validation.go
package domain

import (
	"regexp"
	"strconv"
	"strings"
)

// validateDateFormat проверяет формат даты MM-YYYY
func validateDateFormat(date string) *DomainError {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-20\d{2}$`, date)
	if !matched {
		return NewValidationError("date", "must be in MM-YYYY format (e.g., 12-2024)")
	}
	return nil
}

// validateDateRange проверяет что startDate <= endDate
func validateDateRange(startDate, endDate string) *DomainError {
	startAfterEnd, err := isDateAfter(startDate, endDate)
	if err != nil {
		return NewValidationError("date_range", "invalid date comparison: "+err.Error())
	}
	if startAfterEnd {
		return ErrStartDateAfterEndDate
	}
	return nil
}

// isDateAfter проверяет что date1 > date2
func isDateAfter(date1, date2 string) (bool, *DomainError) {
	year1, month1, err := parseDate(date1)
	if err != nil {
		return false, NewValidationError("date", "invalid first date: "+err.Error())
	}

	year2, month2, err := parseDate(date2)
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

// isDateAfterOrEqual проверяет что date1 >= date2
func isDateAfterOrEqual(date1, date2 string) (bool, *DomainError) {
	after, err := isDateAfter(date1, date2)
	if err != nil {
		return false, err
	}
	if after {
		return true, nil
	}

	// Проверяем равенство
	year1, month1, err := parseDate(date1)
	if err != nil {
		return false, err
	}

	year2, month2, err := parseDate(date2)
	if err != nil {
		return false, err
	}

	return year1 == year2 && month1 == month2, nil
}

// isDateBeforeOrEqual проверяет что date1 <= date2
func isDateBeforeOrEqual(date1, date2 string) (bool, *DomainError) {
	after, err := isDateAfter(date1, date2)
	if err != nil {
		return false, err
	}
	return !after, nil
}

// parseDate парсит строку MM-YYYY в год и месяц
func parseDate(date string) (year, month int, err *DomainError) {
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

// validateSubscriptionDates валидация дат подписки для использования в сервисе
func validateSubscriptionDates(startDate string, endDate *string) *DomainError {
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
