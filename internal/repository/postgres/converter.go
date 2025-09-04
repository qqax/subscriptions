package postgres

import (
	"subscription/core/domain"
	"subscription/internal/repository/postgres/models"
)

// ToDBModel converts domain Subscription to DB model
func ToDBModel(domainSub *domain.Subscription) (*models.Subscription, error) {
	startMonth, startYear, err := parseMMYYYY(domainSub.StartDate)
	if err != nil {
		return nil, err
	}

	dbSub := &models.Subscription{
		ID:          domainSub.ID,
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

// ToDomain converts a DB model to domain Subscription
func ToDomain(dbSub *models.Subscription) (*domain.Subscription, error) {
	startDate := formatMMYYYY(dbSub.StartMonth, dbSub.StartYear)

	var endDate *string
	if dbSub.EndMonth != nil && dbSub.EndYear != nil {
		formatted := formatMMYYYY(*dbSub.EndMonth, *dbSub.EndYear)
		endDate = &formatted
	}

	return domain.NewSubscription(
		dbSub.ID,
		dbSub.ServiceName,
		dbSub.Price,
		dbSub.UserID,
		startDate,
		endDate,
	)
}
