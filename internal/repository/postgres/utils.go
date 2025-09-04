// internal/repository/postgres/utils.go
package postgres

import (
	"context"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"subscription/core/domain"
)

// applyDateFilter применяет фильтрацию по датам для подписок
func applyDateFilter(query *gorm.DB, startDateFrom string, startDateTo *string) *gorm.DB {
	// Парсим начальную дату
	startMonth, startYear, err := parseMMYYYY(startDateFrom)
	if err != nil {
		// В продакшене лучше логировать ошибку и возвращать query без фильтра
		return query
	}

	// Базовый фильтр от начальной даты
	query = query.Where("(start_year > ? OR (start_year = ? AND start_month >= ?))",
		startYear, startYear, startMonth)

	// Если указана конечная дата
	if startDateTo != nil {
		endMonth, endYear, err := parseMMYYYY(*startDateTo)
		if err == nil {
			query = query.Where("(start_year < ? OR (start_year = ? AND start_month <= ?))",
				endYear, endYear, endMonth)
		}
	}

	return query
}

// applyDateRangeFilter применяет фильтрацию по диапазону дат
func applyDateRangeFilter(query *gorm.DB, startDate, endDate string) *gorm.DB {
	startMonth, startYear, err := parseMMYYYY(startDate)
	if err != nil {
		return query
	}

	endMonth, endYear, err := parseMMYYYY(endDate)
	if err != nil {
		return query
	}

	// Фильтр: подписки, которые были активны в любой части периода
	// Подписка активна если её период пересекается с запрашиваемым периодом
	return query.Where(`
		(start_year < ? OR (start_year = ? AND start_month <= ?)) AND
		(end_year IS NULL OR end_year > ? OR (end_year = ? AND end_month >= ?))
	`, endYear, endYear, endMonth, startYear, startYear, startMonth)
}

// applyDateFilterForCost применяет фильтр дат для расчета стоимости
func applyDateFilterForCost(query *gorm.DB, startDate, endDate string) *gorm.DB {
	startMonth, startYear, err := parseMMYYYY(startDate)
	if err != nil {
		return query
	}

	endMonth, endYear, err := parseMMYYYY(endDate)
	if err != nil {
		return query
	}

	// Фильтр для подписок, активных в течение всего или части периода
	return query.Where(`
		(start_year < ? OR (start_year = ? AND start_month <= ?)) AND
		(end_year IS NULL OR end_year > ? OR (end_year = ? AND end_month >= ?))
	`, endYear, endYear, endMonth, startYear, startYear, startMonth)
}

// parseMMYYYY парсит строку формата MM-YYYY
func parseMMYYYY(date string) (month, year int, err error) {
	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return 0, 0, domain.ErrInvalidDateformat
	}

	month, err = strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, domain.ErrInvalidDateformat
	}

	year, err = strconv.Atoi(parts[1])
	if err != nil || year < 2020 {
		return 0, 0, domain.ErrInvalidDateformat
	}

	return month, year, nil
}

// formatMMYYYY форматирует месяц и год в MM-YYYY
func formatMMYYYY(month, year int) string {
	return strings.Join([]string{
		strconv.Itoa(month),
		strconv.Itoa(year),
	}, "-")
}

// getRequestID извлекает request ID из контекста
func getRequestID(ctx context.Context) string {
	// Тип для ключа контекста
	type contextKey string

	// Пробуем получить request ID из контекста
	if requestID, ok := ctx.Value(contextKey("request_id")).(string); ok && requestID != "" {
		return requestID
	}

	// Альтернативные ключи, которые могут использоваться
	keys := []contextKey{"x-request-id", "requestId", "reqId"}
	for _, key := range keys {
		if value, ok := ctx.Value(key).(string); ok && value != "" {
			return value
		}
	}

	return "unknown"
}

// buildWhereINCondition строит условие IN для массивов
func buildWhereINCondition(query *gorm.DB, field string, values []string) *gorm.DB {
	if len(values) == 0 {
		return query
	}
	return query.Where(field+" IN (?)", values)
}

// applyPagination применяет пагинацию к запросу
func applyPagination(query *gorm.DB, offset, limit int) *gorm.DB {
	return query.Offset(offset).Limit(limit)
}

// calculateTotalPages вычисляет общее количество страниц
func calculateTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
