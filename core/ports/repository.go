package ports

import (
	"context"
	"github.com/google/uuid"
	"subscription/core/domain"
)

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	// Create creates a new subscription
	Create(ctx context.Context, subscription *domain.Subscription) *domain.DomainError

	// GetByID returns a subscription by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, *domain.DomainError)

	// List returns server with filtering and pagination
	List(ctx context.Context, filter SubscriptionFilter, pagination Pagination) ([]*domain.Subscription, *PaginationMetadata, *domain.DomainError)

	// Update fully updates a subscription
	Update(ctx context.Context, subscription *domain.Subscription) *domain.DomainError

	// PartialUpdate partially updates a subscription
	PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) *domain.DomainError

	// Delete removes a subscription by ID
	Delete(ctx context.Context, id uuid.UUID) *domain.DomainError

	// GetTotalCost calculates total cost for a period with filters
	GetTotalCost(ctx context.Context, startDate, endDate string, filter SubscriptionFilter) (int, *domain.DomainError)

	// SubscriptionExists проверяет существование подписки
	SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, *domain.DomainError)

	// GetByUserAndService возвращает подписку по user ID и service name
	GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, *domain.DomainError)
}
