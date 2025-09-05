package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"subscription/core/domain"
	"subscription/internal/repository/postgres/model"
	"time"

	"gorm.io/gorm"

	"subscription/core/ports"
	"subscription/internal/logger"

	"github.com/google/uuid"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) ports.SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) (uuid.UUID, error) {
	dbSub, err := ToDBModel(subscription)
	if err != nil {
		return uuid.Nil, err
	}

	result := r.db.WithContext(ctx).Create(dbSub)
	if result.Error != nil {
		logger.Error().Err(result.Error)

		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, domain.ErrDuplicateSubscription
		}

		return uuid.Nil, domain.ErrInternal
	}

	return dbSub.ID, nil
}

// GetByID returns subscription by ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var dbSub model.Subscription
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&dbSub)
	if result.Error != nil {
		logger.Error().Err(result.Error)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSubscriptionNotFound
		}

		return nil, domain.ErrInternal
	}

	return ToDomain(&dbSub)
}

// List returns subscriptions with filtering and pagination
func (r *SubscriptionRepository) List(ctx context.Context, filter ports.SubscriptionFilter, pagination ports.Pagination) ([]*domain.Subscription, *ports.PaginationMetadata, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	query := r.db.WithContext(ctx).Model(&model.Subscription{})

	if len(filter.UserIDs) > 0 {
		query = query.Where("user_id IN ?", filter.UserIDs)
	}
	if len(filter.ServiceNames) > 0 {
		query = query.Where("service_name IN ?", filter.ServiceNames)
	}
	if filter.StartDateFrom != nil {
		query = applyDateFilter(query, *filter.StartDateFrom, filter.StartDateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count subscriptions")
		return nil, nil, domain.ErrInternal
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query = applyPagination(query, offset, pagination.Limit)

	var dbSubs []model.Subscription
	result := query.Find(&dbSubs)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to list subscriptions")
		return nil, nil, domain.ErrInternal
	}

	domainSubs := make([]*domain.Subscription, len(dbSubs))
	for i, dbSub := range dbSubs {
		domainSub, err := ToDomain(&dbSub)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert DB model to domain model")
			return nil, nil, err
		}
		domainSubs[i] = domainSub
	}

	totalPages := calculateTotalPages(int(total), pagination.Limit)

	paginationMeta := &ports.PaginationMetadata{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	log.Debug().Int("count", len(domainSubs)).Msg("Subscriptions listed successfully")
	return domainSubs, paginationMeta, nil
}

// Update renews subscription
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	log := logger.WithRequestID(getRequestID(ctx))

	dbSub, err := ToDBModel(subscription)
	if err != nil {
		log.Error().Err(err).Str("subscription_id", subscription.ID.String()).Msg("Failed to convert domain model to DB model")
		return err
	}

	result := r.db.WithContext(ctx).Save(dbSub)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("subscription_id", subscription.ID.String()).Msg("Failed to update subscription")
		return domain.ErrInternal
	}

	log.Info().Str("subscription_id", subscription.ID.String()).Msg("Subscription updated successfully")
	return nil
}

// PartialUpdate partially renews subscription
func (r *SubscriptionRepository) PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	log := logger.WithRequestID(getRequestID(ctx))

	updates["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&model.Subscription{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("subscription_id", id.String()).Msg("Failed to partially update subscription")
		return domain.ErrInternal
	}

	if result.RowsAffected == 0 {
		log.Debug().Str("subscription_id", id.String()).Msg("Subscription not found for partial update")
		return domain.ErrSubscriptionNotFound
	}

	log.Info().Str("subscription_id", id.String()).Msg("Subscription partially updated successfully")
	return nil
}

// Delete deletes subscription
func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.WithRequestID(getRequestID(ctx))

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Subscription{})
	if result.Error != nil {
		log.Error().Err(result.Error).Str("subscription_id", id.String()).Msg("Failed to delete subscription")
		return domain.ErrInternal
	}

	if result.RowsAffected == 0 {
		log.Debug().Str("subscription_id", id.String()).Msg("Subscription not found for deletion")
		return domain.ErrSubscriptionNotFound
	}

	log.Info().Str("subscription_id", id.String()).Msg("Subscription deleted successfully")
	return nil
}

// GetTotalCost calculates the total cost of subscriptions
func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, startDate, endDate string, filter ports.SubscriptionFilter) (int, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	startMonth, startYear, err := parseMMYYYY(startDate)
	if err != nil {
		return 0, err
	}

	endMonth, endYear, err := parseMMYYYY(endDate)
	if err != nil {
		return 0, err
	}

	startMonths := startYear*12 + startMonth
	endMonths := endYear*12 + endMonth + 1

	query := r.db.WithContext(ctx).Model(&model.Subscription{}).
		Select(`
			SUM(
					CASE
						WHEN (end_year IS NOT NULL AND ? > end_year * 12 + end_month) OR (? < start_year * 12 + start_month)
							THEN 0
						ELSE
							(
								LEAST(COALESCE(end_year * 12 + end_month, ?), ?::bigint) -
								GREATEST(start_year * 12 + start_month, ?::bigint)
								) * price
						END
			) AS total_cost`,
			startMonths, endMonths, endMonths, endMonths, startMonths)

	query = buildWhereINCondition(query, "user_id", filter.UserIDs)
	query = buildWhereINCondition(query, "service_name", filter.ServiceNames)

	var totalCost int
	result := query.Scan(&totalCost)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to calculate total cost")
		return 0, domain.ErrInternal
	}

	log.Debug().Int("total_cost", totalCost).Msg("Total cost calculated successfully")
	return totalCost, nil
}

// SubscriptionExists checks for the existence of a subscription
func (r *SubscriptionRepository) SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	var count int64
	result := r.db.WithContext(ctx).Model(&model.Subscription{}).
		Where("user_id = ? AND service_name = ?", userID, serviceName).
		Count(&count)

	if result.Error != nil {
		log.Error().Err(result.Error).
			Str("user_id", userID.String()).
			Str("service_name", serviceName).
			Msg("Failed to check subscription existence")
		return false, domain.ErrInternal
	}

	exists := count > 0
	log.Debug().
		Str("user_id", userID.String()).
		Str("service_name", serviceName).
		Bool("exists", exists).
		Msg("Subscription existence checked")

	return exists, nil
}

// GetByUserAndService returns subscription by user ID and service name
func (r *SubscriptionRepository) GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	var dbSub model.Subscription
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND service_name = ?", userID, serviceName).
		First(&dbSub)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("user_id", userID.String()).
				Str("service_name", serviceName).
				Msg("Subscription not found")
			return nil, domain.ErrSubscriptionNotFound
		}
		log.Error().Err(result.Error).
			Str("user_id", userID.String()).
			Str("service_name", serviceName).
			Msg("Failed to get subscription by user and service")
		return nil, domain.ErrInternal
	}

	domainSub, err := ToDomain(&dbSub)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID.String()).
			Str("service_name", serviceName).
			Msg("Failed to convert DB model to domain model")
		return nil, err
	}

	log.Debug().
		Str("user_id", userID.String()).
		Str("service_name", serviceName).
		Msg("Subscription found by user and service")

	return domainSub, nil
}
