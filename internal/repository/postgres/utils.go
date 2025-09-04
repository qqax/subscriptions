package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"subscription/internal/logger"

	"gorm.io/gorm"
	"subscription/core/domain"
)

const requestIdKey = "request_id"

// applyDateFilter applies date filtering to subscriptions
func applyDateFilter(query *gorm.DB, startDateFrom string, startDateTo *string) *gorm.DB {
	startMonth, startYear, err := parseMMYYYY(startDateFrom)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse start date")
		return query
	}

	query = query.Where("(start_year > ? OR (start_year = ? AND start_month >= ?))",
		startYear, startYear, startMonth)

	if startDateTo != nil {
		endMonth, endYear, err := parseMMYYYY(*startDateTo)
		if err == nil {
			query = query.Where("(start_year < ? OR (start_year = ? AND start_month <= ?))",
				endYear, endYear, endMonth)
		}
	}

	return query
}

// parseMMYYYY parses a string of format MM-YYYY
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

// formatMMYYYY formats month and year to MM-YYYY
func formatMMYYYY(month, year int) string {
	return strings.Join([]string{
		fmt.Sprintf("%02d", month),
		strconv.Itoa(year),
	}, "-")
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIdKey).(string); ok && requestID != "" {
		return requestID
	}

	return "unknown"
}

// buildWhereINCondition builds IN condition for arrays
func buildWhereINCondition[T any](query *gorm.DB, field string, values []T) *gorm.DB {
	if len(values) == 0 {
		return query
	}
	return query.Where(field+" IN (?)", values)
}

// applyPagination applies pagination to a query
func applyPagination(query *gorm.DB, offset, limit int) *gorm.DB {
	return query.Offset(offset).Limit(limit)
}

// calculateTotalPages calculates the total number of pages
func calculateTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
