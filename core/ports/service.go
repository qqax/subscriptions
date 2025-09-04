package ports

import (
	"context"
	"github.com/google/uuid"
	"subscription/core/domain"
)

// SubscriptionService defines the business logic operations for server
type SubscriptionService interface {
	// CreateSubscription creates a new subscription
	CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*domain.Subscription, error)

	// GetSubscription returns a subscription by ID
	GetSubscription(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)

	// ListSubscriptions returns server with filtering and pagination
	ListSubscriptions(ctx context.Context, filter SubscriptionFilter, pagination Pagination) ([]*domain.Subscription, *PaginationMetadata, error)

	// UpdateSubscription fully updates a subscription
	UpdateSubscription(ctx context.Context, id uuid.UUID, req *UpdateSubscriptionRequest) (*domain.Subscription, error)

	// PartialUpdateSubscription partially updates a subscription
	PartialUpdateSubscription(ctx context.Context, id uuid.UUID, req *PartialUpdateRequest) (*domain.Subscription, error)

	// DeleteSubscription removes a subscription by ID
	DeleteSubscription(ctx context.Context, id uuid.UUID) error

	// GetTotalCost calculates total subscription cost for period
	GetTotalCost(ctx context.Context, req *TotalCostRequest) (*TotalCostResponse, error)
}
