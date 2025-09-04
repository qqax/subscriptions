package postgres

import (
	"context"
	"strconv"
	"strings"
	"subscription/internal/logger"

	"gorm.io/gorm"
	"subscription/core/domain"
)

const requestIdKey = "requestId"

// applyDateFilter применяет фильтрацию по датам для подписок
func applyDateFilter(query *gorm.DB, startDateFrom string, startDateTo *string) *gorm.DB {
	startMonth, startYear, err := parseMMYYYY(startDateFrom)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse start date")
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
		logger.Error().Str("date", date).Err(err).Msg("failed to parse date")
		return 0, 0, domain.ErrInvalidDateformat
	}

	month, e := strconv.Atoi(parts[0])
	if e != nil || month < 1 || month > 12 {
		logger.Error().Str("date", date).Err(err).Msg("failed to parse date")
		return 0, 0, domain.ErrInvalidDateformat
	}

	year, e = strconv.Atoi(parts[1])
	if e != nil {
		logger.Error().Str("date", date).Err(err).Msg("failed to parse date")
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
	if requestID, ok := ctx.Value(requestIdKey).(string); ok && requestID != "" {
		return requestID
	}

	return "unknown"
}

// buildWhereINCondition строит условие IN для массивов
func buildWhereINCondition[T any](query *gorm.DB, field string, values []T) *gorm.DB {
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
