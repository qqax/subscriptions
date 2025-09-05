package ports

import "github.com/google/uuid"

// CreateSubscriptionRequest represents the request to create a subscription
type CreateSubscriptionRequest struct {
	EndDate     *string   `json:"end_date" validate:"omitempty,mm_yyyy_format"`
	ServiceName string    `json:"service_name" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy_format"`
	Price       int       `json:"price" validate:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
}

// UpdateSubscriptionRequest represents the request to update a subscription
type UpdateSubscriptionRequest struct {
	EndDate     *string   `json:"end_date" validate:"omitempty,mm_yyyy_format"`
	ServiceName string    `json:"service_name" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy_format"`
	Price       int       `json:"price" validate:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
}

// PartialUpdateRequest represents the request for partial update
type PartialUpdateRequest struct {
	ServiceName *string    `json:"service_name" validate:"omitempty"`
	Price       *int       `json:"price" validate:"omitempty,min=1"`
	UserID      *uuid.UUID `json:"user_id" validate:"omitempty,uuid4"`
	StartDate   *string    `json:"start_date" validate:"omitempty,mm_yyyy_format"`
	EndDate     *string    `json:"end_date" validate:"omitempty,mm_yyyy_format"`
}

// TotalCostRequest represents the request for total cost calculation
type TotalCostRequest struct {
	StartDate    string      `json:"start_date" validate:"required,mm_yyyy_format"`
	EndDate      string      `json:"end_date" validate:"required,mm_yyyy_format"`
	UserIDs      []uuid.UUID `json:"user_ids" validate:"omitempty,dive,uuid4"`
	ServiceNames []string    `json:"service_names" validate:"omitempty"`
}

// TotalCostResponse represents the response for total cost calculation
type TotalCostResponse struct {
	Period         Period                  `json:"period"`
	FilterCriteria TotalCostFilterCriteria `json:"filter_criteria"`
	TotalCost      int                     `json:"total_cost"`
}

// Period represents a date period
type Period struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// TotalCostFilterCriteria represents filter criteria used in total cost calculation
type TotalCostFilterCriteria struct {
	UserIDs      []uuid.UUID `json:"user_ids"`
	ServiceNames []string    `json:"service_names"`
}
