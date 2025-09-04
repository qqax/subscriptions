package ogen

import (
	"context"
	"errors"
	"subscription/core/domain"
	api "subscription/internal/api/generated"
	"time"
)

// Error conversion functions for each operation

func convertSubscriptionsPostError(err error) *api.SubscriptionsPostBadRequest {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsPostBadRequest)(&errorResponse)
}

func convertSubscriptionsIDGetError(err error) *api.SubscriptionsIDGetNotFound {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsIDGetNotFound)(&errorResponse)
}

func convertSubscriptionsGetError(err error) *api.SubscriptionsGetBadRequest {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsGetBadRequest)(&errorResponse)
}

func convertSubscriptionsIDPutError(err error) *api.SubscriptionsIDPutBadRequest {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsIDPutBadRequest)(&errorResponse)
}

func convertSubscriptionsIDPatchError(err error) *api.SubscriptionsIDPatchBadRequest {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsIDPatchBadRequest)(&errorResponse)
}

func convertSubscriptionsIDDeleteError(err error) *api.SubscriptionsIDDeleteNotFound {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsIDDeleteNotFound)(&errorResponse)
}

func convertSubscriptionsSummaryTotalCostGetError(err error) *api.SubscriptionsSummaryTotalCostGetBadRequest {
	errorResponse := createErrorResponse(err)
	return (*api.SubscriptionsSummaryTotalCostGetBadRequest)(&errorResponse)
}

// Helper functions for creating error responses

func createErrorResponse(err error) api.Error {
	statusCode := getStatusCode(err)
	errorCode := getErrorCode(err)

	return api.Error{
		Error:     errorCode,
		Message:   getErrorMessage(err),
		Code:      api.NewOptInt32(int32(statusCode)),
		Timestamp: api.NewOptDateTime(time.Now()),
		Details:   api.OptErrorDetails{},
	}
}

func getStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		return 404
	case errors.Is(err, domain.ErrInvalidDateformat),
		errors.Is(err, domain.ErrInvalidUUID),
		errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrStartDateAfterEndDate),
		errors.Is(err, domain.ErrInvalidDateRange):
		return 400
	case errors.Is(err, domain.ErrDuplicateSubscription):
		return 409
	default:
		return 500
	}
}

func getErrorCode(err error) string {
	switch {
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		return "not_found"
	case errors.Is(err, domain.ErrInvalidDateformat):
		return "invalid_date_format"
	case errors.Is(err, domain.ErrInvalidUUID):
		return "invalid_uuid"
	case errors.Is(err, domain.ErrInvalidPrice):
		return "invalid_price"
	case errors.Is(err, domain.ErrStartDateAfterEndDate):
		return "invalid_date_range"
	case errors.Is(err, domain.ErrInvalidDateRange):
		return "invalid_date_range"
	case errors.Is(err, domain.ErrDuplicateSubscription):
		return "duplicate_subscription"
	default:
		return "internal_error"
	}
}

func getErrorMessage(err error) string {
	// Для стандартных ошибок возвращаем их текст
	// Для кастомных можно добавать дополнительную информацию
	return err.Error()
}

// NewError implements api.Handler.
func (h *OgenAdapter) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: getStatusCode(err),
		Response:   createErrorResponse(err),
	}
}
func getStatusCodeFromDomainError(err *domain.DomainError) int {
	switch err.Code {
	case "validation_error", "invalid_input":
		return 400
	case "not_found":
		return 404
	case "conflict", "duplicate":
		return 409
	case "unauthorized", "forbidden":
		return 403
	default:
		return 500
	}
}
