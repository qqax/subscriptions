package postgres

import (
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
