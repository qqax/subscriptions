package postgres

import (
	"fmt"
	"strconv"
	"strings"
	"subscription/internal/repository/postgres/models"

	"subscription/core/domain"
)

// ToDBModel converts domain Subscription to DB model
func ToDBModel(domainSub *domain.Subscription) (*models.Subscription, error) {
	startMonth, startYear, err := parseMMYYYY(domainSub.StartDate)
	if err != nil {
		return nil, err
	}

	dbSub := &models.Subscription{
		ServiceName: domainSub.ServiceName,
		Price:       domainSub.Price,
		UserID:      domainSub.UserID,
		StartMonth:  startMonth,
		StartYear:   startYear,
	}

	if domainSub.EndDate != nil {
		endMonth, endYear, err := parseMMYYYY(*domainSub.EndDate)
		if err != nil {
			return nil, err
		}
		dbSub.EndMonth = &endMonth
		dbSub.EndYear = &endYear
	}

	return dbSub, nil
}

// ToDomain converts DB model to domain Subscription

func ToDomain(dbSub *models.Subscription) (*domain.Subscription, error) {
	startDate := formatMMYYYY(dbSub.StartMonth, dbSub.StartYear)

	var endDate *string
	if dbSub.EndMonth != nil && dbSub.EndYear != nil {
		formatted := formatMMYYYY(*dbSub.EndMonth, *dbSub.EndYear)
		endDate = &formatted
	}

	// Создаем domain модель
	return domain.NewSubscription(
		dbSub.ServiceName,
		dbSub.Price,
		dbSub.UserID,
		startDate,
		endDate,
	)
}

// Helper functions
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

func formatMMYYYY(month, year int) string {
	return fmt.Sprintf("%02d-%d", month, year)
}
