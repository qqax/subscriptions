package ports

import (
	"context"
	"github.com/google/uuid"
	"subscription/core/domain"
	"time"
)

// Subscription represents a user subscription entity
type Subscription struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"` // Format: MM-YYYY
	EndDate     *string   `json:"end_date"`   // Format: MM-YYYY, nullable
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//// SubscriptionRepository defines the interface for subscription data operations
//type SubscriptionRepository interface {
//	// Create creates a new subscription
//	Create(ctx context.Context, subscription *Subscription) error
//
//	// GetByID returns a subscription by its ID
//	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
//
//	// List returns server with filtering and pagination
//	List(ctx context.Context, filter SubscriptionFilter, pagination Pagination) ([]*Subscription, *PaginationMetadata, error)
//
//	// Update fully updates a subscription
//	Update(ctx context.Context, subscription *Subscription) error
//
//	// PartialUpdate partially updates a subscription
//	PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
//
//	// Delete removes a subscription by ID
//	Delete(ctx context.Context, id uuid.UUID) error
//
//	// GetTotalCost calculates total cost for period with filters
//	GetTotalCost(ctx context.Context, startDate, endDate string, filter SubscriptionFilter) (int, error)
//
//	// SubscriptionExists проверяет существование подписки
//	SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, error)
//
//	// GetByUserAndService возвращает подписку по user ID и service name
//	GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*Subscription, error)
//}

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, filter SubscriptionFilter, pagination Pagination) ([]*domain.Subscription, *PaginationMetadata, error)
	Update(ctx context.Context, subscription *domain.Subscription) error
	PartialUpdate(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, startDate, endDate string, filter SubscriptionFilter) (int, error)
	SubscriptionExists(ctx context.Context, userID uuid.UUID, serviceName string) (bool, error)
	GetByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string) (*domain.Subscription, error)
}

// SubscriptionFilter contains filtering criteria for server
type SubscriptionFilter struct {
	UserIDs       []uuid.UUID `json:"user_ids"`
	ServiceNames  []string    `json:"service_names"`
	StartDateFrom *string     `json:"start_date_from"` // MM-YYYY
	StartDateTo   *string     `json:"start_date_to"`   // MM-YYYY
	EndDateNull   *bool       `json:"end_date_null"`   // Filter by null end date
}

// Pagination contains pagination parameters
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginationMetadata contains pagination metadata
type PaginationMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
