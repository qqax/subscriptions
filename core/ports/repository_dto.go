package ports

import "github.com/google/uuid"

// SubscriptionFilter contains filtering criteria for server
type SubscriptionFilter struct {
	UserIDs       []uuid.UUID `json:"user_ids" validate:"omitempty,dive,uuid4"`
	ServiceNames  []string    `json:"service_names" validate:"omitempty"`
	StartDateFrom *string     `json:"start_date_from" validate:"omitempty,mm_yyyy_format"`
	StartDateTo   *string     `json:"start_date_to" validate:"omitempty,mm_yyyy_format"`
	EndDateNull   *bool       `json:"end_date_null" validate:"omitempty"`
}

// Pagination contains pagination parameters
type Pagination struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// PaginationMetadata contains pagination metadata
type PaginationMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
