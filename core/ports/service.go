package ports

import (
	"context"
	"github.com/google/uuid"
	"subscription/core/domain"
)

// SubscriptionService defines the business logic operations for server
type SubscriptionService interface {
	// CreateSubscription creates a new subscription
	CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*domain.Subscription, *domain.DomainError)

	// GetSubscription returns a subscription by ID
	GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, *domain.DomainError)

	// ListSubscriptions returns server with filtering and pagination
	ListSubscriptions(ctx context.Context, filter SubscriptionFilter, pagination Pagination) ([]*domain.Subscription, *PaginationMetadata, *domain.DomainError)

	// UpdateSubscription fully updates a subscription
	UpdateSubscription(ctx context.Context, id uuid.UUID, req *UpdateSubscriptionRequest) (*domain.Subscription, *domain.DomainError)

	// PartialUpdateSubscription partially updates a subscription
	PartialUpdateSubscription(ctx context.Context, id uuid.UUID, req *PartialUpdateRequest) (*domain.Subscription, *domain.DomainError)

	// DeleteSubscription removes a subscription by ID
	DeleteSubscription(ctx context.Context, id uuid.UUID) *domain.DomainError

	// GetTotalCost calculates total subscription cost for period
	GetTotalCost(ctx context.Context, req *TotalCostRequest) (*TotalCostResponse, *domain.DomainError)
}
