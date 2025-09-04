package postgres

import (
	"context"
	"errors"
	"subscription/core/domain"
	"subscription/internal/repository/postgres/models"
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
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) *domain.DomainError {
	dbSub, err := ToDBModel(subscription)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(dbSub)
	if result.Error != nil {
		logger.Error().Err(result.Error)

		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicateSubscription
		}

		return domain.ErrInternal
	}
	return nil
}

// GetByID returns subscription by ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, *domain.DomainError) {
	var dbSub models.Subscription
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

// List возвращает подписки с фильтрацией и пагинацией
func (r *SubscriptionRepository) List(ctx context.Context, filter ports.SubscriptionFilter, pagination ports.Pagination) ([]*domain.Subscription, *ports.PaginationMetadata, *domain.DomainError) {
	log := logger.WithRequestID(getRequestID(ctx))

	query := r.db.WithContext(ctx).Model(&models.Subscription{})

	// Применяем фильтры
	if len(filter.UserIDs) > 0 {
		query = query.Where("user_id IN ?", filter.UserIDs)
	}
	if len(filter.ServiceNames) > 0 {
		query = query.Where("service_name IN ?", filter.ServiceNames)
	}
	if filter.StartDateFrom != nil {
		query = applyDateFilter(query, *filter.StartDateFrom, filter.StartDateTo)
	}

	// Получаем общее количество для пагинации
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count subscriptions")
		return nil, nil, domain.ErrInternal
	}

	// Применяем пагинацию

	offset := (pagination.Page - 1) * pagination.Limit
	query = applyPagination(query, offset, pagination.Limit)

	var dbSubs []models.Subscription
	result := query.Find(&dbSubs)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to list subscriptions")
		return nil, nil, domain.ErrInternal
	}

	// Конвертируем в domain модели (изменили тип возвращаемого значения)
	domainSubs := make([]*domain.Subscription, len(dbSubs))
	for i, dbSub := range dbSubs {
		domainSub, err := ToDomain(&dbSub)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert DB model to domain model")
			return nil, nil, err
		}
		domainSubs[i] = domainSub
	}

	// Метаданные пагинации
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

// Update обновляет подписку
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) *domain.DomainError {
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

// PartialUpdate частично обновляет подписку
func (r *SubscriptionRepository) PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) *domain.DomainError {
	log := logger.WithRequestID(getRequestID(ctx))

	updates["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&models.Subscription{}).Where("id = ?", id).Updates(updates)
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

// Delete удаляет подписку
func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) *domain.DomainError {
	log := logger.WithRequestID(getRequestID(ctx))

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Subscription{})
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

// GetTotalCost вычисляет общую стоимость подписок
func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, startDate, endDate string, filter ports.SubscriptionFilter) (int, *domain.DomainError) {
	log := logger.WithRequestID(getRequestID(ctx))

	query := r.db.WithContext(ctx).Model(&models.Subscription{}).
		Select("COALESCE(SUM(price), 0) as total_cost")

	query = buildWhereINCondition(query, "user_id", filter.UserIDs)
	query = buildWhereINCondition(query, "service_name", filter.ServiceNames)

	query = applyDateRangeFilter(query, startDate, endDate)

	var totalCost int
	result := query.Scan(&totalCost)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to calculate total cost")
		return 0, domain.ErrInternal
	}

	log.Debug().Int("total_cost", totalCost).Msg("Total cost calculated successfully")
	return totalCost, nil
}

// SubscriptionExists проверяет существование подписки
func (r *SubscriptionRepository) SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, *domain.DomainError) {
	log := logger.WithRequestID(getRequestID(ctx))

	var count int64
	result := r.db.WithContext(ctx).Model(&models.Subscription{}).
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
func (r *SubscriptionRepository) GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, *domain.DomainError) {
	log := logger.WithRequestID(getRequestID(ctx))

	var dbSub models.Subscription
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
