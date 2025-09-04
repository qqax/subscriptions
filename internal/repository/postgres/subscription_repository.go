// internal/repository/postgres/subscription_repository.go
package postgres

import (
	"context"
	"fmt"
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

// Create создает новую подписку
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	dbSub, err := ToDBModel(subscription)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Create(dbSub)
	return result.Error
}

// GetByID возвращает подписку по ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var dbSub models.Subscription
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&dbSub)
	if result.Error != nil {
		return nil, result.Error
	}

	return ToDomain(&dbSub)
}

// List возвращает подписки с фильтрацией и пагинацией
func (r *SubscriptionRepository) List(ctx context.Context, filter ports.SubscriptionFilter, pagination ports.Pagination) ([]*domain.Subscription, *ports.PaginationMetadata, error) {
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
		return nil, nil, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	// Применяем пагинацию
	offset := (pagination.Page - 1) * pagination.Limit
	query = query.Offset(offset).Limit(pagination.Limit)

	var dbSubs []models.Subscription
	result := query.Find(&dbSubs)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to list subscriptions")
		return nil, nil, fmt.Errorf("failed to list subscriptions: %w", result.Error)
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
	totalPages := (int(total) + pagination.Limit - 1) / pagination.Limit
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
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	log := logger.WithRequestID(getRequestID(ctx))

	// Конвертируем domain модель в DB модель
	dbSub, err := ToDBModel(subscription)
	if err != nil {
		log.Error().Err(err).Str("subscription_id", subscription.ID.String()).Msg("Failed to convert domain model to DB model")
		return err
	}

	// Сохраняем DB модель, а не domain модель!
	result := r.db.WithContext(ctx).Save(dbSub)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("subscription_id", subscription.ID.String()).Msg("Failed to update subscription")
		return fmt.Errorf("failed to update subscription: %w", result.Error)
	}

	log.Info().Str("subscription_id", subscription.ID.String()).Msg("Subscription updated successfully")
	return nil
}

// PartialUpdate частично обновляет подписку
func (r *SubscriptionRepository) PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	log := logger.WithRequestID(getRequestID(ctx))

	updates["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&models.Subscription{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("subscription_id", id.String()).Msg("Failed to partially update subscription")
		return fmt.Errorf("failed to partially update subscription: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Debug().Str("subscription_id", id.String()).Msg("Subscription not found for partial update")
		return domain.ErrSubscriptionNotFound
	}

	log.Info().Str("subscription_id", id.String()).Msg("Subscription partially updated successfully")
	return nil
}

// Delete удаляет подписку
func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.WithRequestID(getRequestID(ctx))

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Subscription{})
	if result.Error != nil {
		log.Error().Err(result.Error).Str("subscription_id", id.String()).Msg("Failed to delete subscription")
		return fmt.Errorf("failed to delete subscription: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Debug().Str("subscription_id", id.String()).Msg("Subscription not found for deletion")
		return domain.ErrSubscriptionNotFound
	}

	log.Info().Str("subscription_id", id.String()).Msg("Subscription deleted successfully")
	return nil
}

// GetTotalCost вычисляет общую стоимость подписок
func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, startDate, endDate string, filter ports.SubscriptionFilter) (int, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	query := r.db.WithContext(ctx).Model(&models.Subscription{}).
		Select("COALESCE(SUM(price), 0) as total_cost")

	// Применяем фильтры
	if len(filter.UserIDs) > 0 {
		query = query.Where("user_id IN ?", filter.UserIDs)
	}
	if len(filter.ServiceNames) > 0 {
		query = query.Where("service_name IN ?", filter.ServiceNames)
	}

	// Фильтр по датам (упрощенная реализация)
	query = applyDateRangeFilter(query, startDate, endDate)

	var totalCost int
	result := query.Scan(&totalCost)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to calculate total cost")
		return 0, fmt.Errorf("failed to calculate total cost: %w", result.Error)
	}

	log.Debug().Int("total_cost", totalCost).Msg("Total cost calculated successfully")
	return totalCost, nil
}

// SubscriptionExists проверяет существование подписки
func (r *SubscriptionRepository) SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, error) {
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
		return false, fmt.Errorf("failed to check subscription existence: %w", result.Error)
	}

	exists := count > 0
	log.Debug().
		Str("user_id", userID.String()).
		Str("service_name", serviceName).
		Bool("exists", exists).
		Msg("Subscription existence checked")

	return exists, nil
}

// GetByUserAndService возвращает подписку по user ID и service name
func (r *SubscriptionRepository) GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	var dbSub models.Subscription
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND service_name = ?", userID, serviceName).
		First(&dbSub)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
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
		return nil, fmt.Errorf("failed to get subscription by user and service: %w", result.Error)
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

	return domainSub, nil // Теперь возвращаем domain.Subscription
}
