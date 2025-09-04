// core/domain/validation.go
package domain

import (
	"regexp"
	"strconv"
	"strings"
)

// validateDateFormat проверяет формат даты MM-YYYY
func validateDateFormat(date string) error {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-20\d{2}$`, date)
	if !matched {
		return NewValidationError("date", "must be in MM-YYYY format (e.g., 12-2024)")
	}
	return nil
}

// validateDateRange проверяет что startDate <= endDate
func validateDateRange(startDate, endDate string) error {
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
func isDateAfter(date1, date2 string) (bool, error) {
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
func isDateAfterOrEqual(date1, date2 string) (bool, error) {
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
func isDateBeforeOrEqual(date1, date2 string) (bool, error) {
	after, err := isDateAfter(date1, date2)
	if err != nil {
		return false, err
	}
	return !after, nil
}

// calculateOverlapMonths вычисляет количество пересекающихся месяцев
func calculateOverlapMonths(subStart, subEnd, periodStart, periodEnd string) (int, error) {
	// Конвертируем даты в числовой формат (YYYYMM)
	start1, err := dateToNumber(subStart)
	if err != nil {
		return 0, NewValidationError("subscription_start", "invalid date: "+err.Error())
	}

	end1, err := dateToNumber(subEnd)
	if err != nil {
		return 0, NewValidationError("subscription_end", "invalid date: "+err.Error())
	}

	start2, err := dateToNumber(periodStart)
	if err != nil {
		return 0, NewValidationError("period_start", "invalid date: "+err.Error())
	}

	end2, err := dateToNumber(periodEnd)
	if err != nil {
		return 0, NewValidationError("period_end", "invalid date: "+err.Error())
	}

	// Проверяем есть ли пересечение
	if end1 < start2 || start1 > end2 {
		return 0, nil // Нет пересечения
	}

	// Находим начало и конец пересечения
	overlapStart := max(start1, start2)
	overlapEnd := min(end1, end2)

	// Вычисляем количество месяцев
	startYear := overlapStart / 100
	startMonth := overlapStart % 100
	endYear := overlapEnd / 100
	endMonth := overlapEnd % 100

	months := (endYear-startYear)*12 + (endMonth - startMonth) + 1
	return months, nil
}

// parseDate парсит строку MM-YYYY в год и месяц
func parseDate(date string) (year, month int, err error) {
	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return 0, 0, NewValidationError("date_format", "invalid format, expected MM-YYYY")
	}

	month, err = strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, NewValidationError("date_month", "month must be between 01 and 12")
	}

	year, err = strconv.Atoi(parts[1])
	if err != nil || year < 2000 || year > 2100 {
		return 0, 0, NewValidationError("date_year", "year must be between 2000 and 2100")
	}

	return year, month, nil
}

// dateToNumber конвертирует MM-YYYY в числовой формат YYYYMM
func dateToNumber(date string) (int, error) {
	year, month, err := parseDate(date)
	if err != nil {
		return 0, err
	}
	return year*100 + month, nil
}

// ValidateSubscriptionDates валидация дат подписки для использования в сервисе
func ValidateSubscriptionDates(startDate string, endDate *string) error {
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
