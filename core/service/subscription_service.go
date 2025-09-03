// core/service/subscription_service.go
package service

import (
	"context"
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

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *ports.CreateSubscriptionRequest) (*ports.Subscription, error) {
	// Validate business rules
	if err := validateSubscriptionDates(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	subscription := &ports.Subscription{
		ID:          generateID(), // Implement UUID generation
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id string) (*ports.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, filter ports.SubscriptionFilter, pagination ports.Pagination) ([]*ports.Subscription, *ports.PaginationMetadata, error) {
	return s.repo.List(ctx, filter, pagination)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id string, req *ports.UpdateSubscriptionRequest) (*ports.Subscription, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validateSubscriptionDates(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	existing.ServiceName = req.ServiceName
	existing.Price = req.Price
	existing.UserID = req.UserID
	existing.StartDate = req.StartDate
	existing.EndDate = req.EndDate
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *subscriptionService) PartialUpdateSubscription(ctx context.Context, id string, req *ports.PartialUpdateRequest) (*ports.Subscription, error) {
	updates := make(map[string]interface{})

	if req.ServiceName != nil {
		updates["service_name"] = *req.ServiceName
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.UserID != nil {
		updates["user_id"] = *req.UserID
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}

	updates["updated_at"] = time.Now()

	if err := s.repo.PartialUpdate(ctx, id, updates); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) GetTotalCost(ctx context.Context, req *ports.TotalCostRequest) (*ports.TotalCostResponse, error) {
	filter := ports.SubscriptionFilter{
		UserIDs:      req.UserIDs,
		ServiceNames: req.ServiceNames,
	}

	totalCost, err := s.repo.GetTotalCost(ctx, req.StartDate, req.EndDate, filter)
	if err != nil {
		return nil, err
	}

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

// Helper functions
func validateSubscriptionDates(startDate string, endDate *string) error {
	// Implement date validation logic
	// Check if startDate is valid MM-YYYY
	// Check if endDate is after startDate if provided
	return nil
}

func generateID() string {
	// Implement UUID generation
	return "generated-uuid"
}
