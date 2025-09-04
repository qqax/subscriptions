package ports

import (
	"context"
	"github.com/google/uuid"
	"subscription/core/domain"
)

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	// Create creates a new subscription
	Create(ctx context.Context, subscription *domain.Subscription) (uuid.UUID, error)

	// GetByID returns a subscription by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)

	// List returns server with filtering and pagination
	List(ctx context.Context, filter SubscriptionFilter, pagination Pagination) ([]*domain.Subscription, *PaginationMetadata, error)

	// Update fully updates a subscription
	Update(ctx context.Context, subscription *domain.Subscription) error

	// PartialUpdate partially updates a subscription
	PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error

	// Delete removes a subscription by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// GetTotalCost calculates total cost for a period with filters
	GetTotalCost(ctx context.Context, startDate, endDate string, filter SubscriptionFilter) (int, error)

	// SubscriptionExists проверяет существование подписки
	SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, error)

	// GetByUserAndService возвращает подписку по user ID и service name
	GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, error)
}
