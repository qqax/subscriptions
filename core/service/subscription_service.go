package service

import (
	"context"
	"github.com/google/uuid"
	"subscription/core/domain"
	"time"

	"subscription/core/ports"
)

type subscriptionService struct {
	repo ports.SubscriptionRepository
	// validator could be added here
}

func NewSubscriptionService(repo ports.SubscriptionRepository) ports.SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *ports.CreateSubscriptionRequest) (*domain.Subscription, *domain.DomainError) {
	subscription, err := domain.NewSubscription(
		req.ServiceName,
		req.Price,
		req.UserID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, *domain.DomainError) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, filter ports.SubscriptionFilter, pagination ports.Pagination) ([]*domain.Subscription, *ports.PaginationMetadata, *domain.DomainError) {
	// Валидация пагинации
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.Limit < 1 {
		pagination.Limit = 20
	} else if pagination.Limit > 100 {
		pagination.Limit = 100
	}

	// Валидация фильтров
	if err := validateFilter(filter); err != nil {
		return nil, nil, domain.ErrValidationFailed
	}

	// Вызов репозитория
	return s.repo.List(ctx, filter, pagination)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, req *ports.UpdateSubscriptionRequest) (*domain.Subscription, *domain.DomainError) {
	// Получаем существующую подписку
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Валидация данных
	if err = validateSubscriptionDates(req.StartDate, req.EndDate); err != nil {
		return nil, domain.ErrInvalidDateRange
	}

	// Обновляем поля
	existing.ServiceName = req.ServiceName
	existing.Price = req.Price
	existing.UserID = req.UserID
	existing.StartDate = req.StartDate
	existing.EndDate = req.EndDate

	// Сохраняем изменения
	if err = s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *subscriptionService) PartialUpdateSubscription(ctx context.Context, id uuid.UUID, req *ports.PartialUpdateRequest) (*domain.Subscription, *domain.DomainError) {
	// Создаем map для обновлений
	updates := make(map[string]interface{})

	if req.Price != nil {
		updates["price"] = *req.Price
	}

	if req.EndDate != nil {
		if *req.EndDate != "" {
			if err := validateDateFormat(*req.EndDate); err != nil {
				return nil, domain.ErrInvalidDateformat
			}
		}
		updates["end_date"] = *req.EndDate
	}

	// Добавляем updated_at
	updates["updated_at"] = time.Now()

	// Вызываем репозиторий для частичного обновления
	if err := s.repo.PartialUpdate(ctx, id, updates); err != nil {
		return nil, err
	}

	// Возвращаем обновленную подписку
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) *domain.DomainError {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) GetTotalCost(ctx context.Context, req *ports.TotalCostRequest) (*ports.TotalCostResponse, *domain.DomainError) {
	// Валидация дат
	if err := validateDateFormat(req.StartDate); err != nil {
		return nil, domain.ErrInvalidDateformat
	}
	if err := validateDateFormat(req.EndDate); err != nil {
		return nil, domain.ErrInvalidDateformat
	}

	// Проверка корректности диапазона дат
	if err := validateDateRange(req.StartDate, req.EndDate); err != nil {
		return nil, domain.ErrInvalidDateRange
	}

	// Создаем фильтр для репозитория
	filter := ports.SubscriptionFilter{
		UserIDs:      req.UserIDs,
		ServiceNames: req.ServiceNames,
	}

	// Вычисляем общую стоимость
	totalCost, err := s.repo.GetTotalCost(ctx, req.StartDate, req.EndDate, filter)
	if err != nil {
		return nil, err
	}

	// Формируем ответ
	return &ports.TotalCostResponse{
		TotalCost: totalCost,
		Period: ports.Period{
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		},
		FilterCriteria: ports.TotalCostFilterCriteria{
			UserIDs:      req.UserIDs,
			ServiceNames: req.ServiceNames,
		},
	}, nil
}
