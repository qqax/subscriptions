package ogen

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"subscription/core/domain"
	"subscription/core/ports"
	"subscription/internal/logger"
	"time"

	api "subscription/internal/api/generated" // сгенерированный ogen код
)

type OgenAdapter struct {
	service ports.SubscriptionService
}

func NewOgenAdapter(service ports.SubscriptionService) *OgenAdapter {
	return &OgenAdapter{service: service}
}

// Ensure interface implementation
var _ api.Handler = (*OgenAdapter)(nil)

// SubscriptionsGet implements api.Handler.
func (h *OgenAdapter) SubscriptionsGet(ctx context.Context, params api.SubscriptionsGetParams) (api.SubscriptionsGetRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	// Convert ogen params to domain filter/pagination
	filter := convertFilterParams(params)

	pagination := ports.Pagination{
		Page:  getIntOrDefault(params.Page.Get, 1),
		Limit: getIntOrDefault(params.Limit.Get, 20),
	}

	// Call domain service
	subscriptions, paginationMeta, err := h.service.ListSubscriptions(ctx, filter, pagination)
	if err != nil {
		msg := "Failed to list server"
		log.Error().Err(err).Msg(msg)
		return &api.SubscriptionsGetInternalServerError{
			Error:   err.Error(),
			Message: msg,
		}, err
	}

	// Convert to ogen response
	response := &api.SubscriptionsGetOK{
		Data:       convertSubscriptionsToOgen(subscriptions),
		Pagination: convertPaginationToOgen(paginationMeta),
	}

	return response, nil
}

// SubscriptionsPost implements api.Handler.
func (h *OgenAdapter) SubscriptionsPost(ctx context.Context, req *api.SubscriptionCreate) (api.SubscriptionsPostRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	userID, err := uuid.Parse(req.UserID.String())
	if err != nil {
		return convertSubscriptionsPostError(ports.ErrInvalidUUID), nil
	}

	// Convert ogen request to domain request
	domainReq := &ports.CreateSubscriptionRequest{
		ServiceName: req.ServiceName,
		Price:       int(req.Price),
		UserID:      userID,
		StartDate:   req.StartDate,
		EndDate:     getStringPtrFromOptNil(req.EndDate),
	}

	// Call domain service
	subscription, err := h.service.CreateSubscription(ctx, domainReq)
	if err != nil {

		log.Error().Err(err).Msg("Failed to create subscription")
		return convertSubscriptionsPostError(err), nil

	}

	// Convert domain response to ogen response
	return convertSubscriptionToOgen(subscription), nil
}

// SubscriptionsIDGet implements api.Handler.
func (h *OgenAdapter) SubscriptionsIDGet(ctx context.Context, params api.SubscriptionsIDGetParams) (api.SubscriptionsIDGetRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	subscription, err := h.service.GetSubscription(ctx, params.ID)
	if err != nil {
		log.Error().Err(err).Str("subscription_id", params.ID.String()).Msg("Failed to get subscription")
		return convertSubscriptionsIDGetError(err), nil
	}

	return convertSubscriptionToOgen(subscription), nil
}

// SubscriptionsIDPut implements api.Handler.
func (h *OgenAdapter) SubscriptionsIDPut(ctx context.Context, req *api.SubscriptionUpdate, params api.SubscriptionsIDPutParams) (api.SubscriptionsIDPutRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	domainReq := &ports.UpdateSubscriptionRequest{
		ServiceName: req.ServiceName,
		Price:       int(req.Price),
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     getStringPtrFromOptNil(req.EndDate),
	}

	subscription, err := h.service.UpdateSubscription(ctx, params.ID, domainReq)
	if err != nil {
		log.Error().Err(err).Str("subscription_id", params.ID.String()).Msg("Failed to update subscription")
		return convertSubscriptionsIDPutError(err), nil
	}

	return convertSubscriptionToOgen(subscription), nil
}

// SubscriptionsIDPatch implements api.Handler.
func (h *OgenAdapter) SubscriptionsIDPatch(ctx context.Context, req *api.SubscriptionPatch, params api.SubscriptionsIDPatchParams) (api.SubscriptionsIDPatchRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	domainReq := &ports.PartialUpdateRequest{
		ServiceName: getStringPtrFromOpt(req.ServiceName),
		Price:       getIntPtrFromOpt(req.Price),
		UserID:      getStringPtrFromUUIDOpt(req.UserID),
		StartDate:   getStringPtrFromOpt(req.StartDate),
		EndDate:     getStringPtrFromOptNil(req.EndDate),
	}

	subscription, err := h.service.PartialUpdateSubscription(ctx, params.ID, domainReq)
	if err != nil {
		log.Error().Err(err).Str("subscription_id", params.ID.String()).Msg("Failed to partially update subscription")
		return convertSubscriptionsIDPatchError(err), nil
	}

	return convertSubscriptionToOgen(subscription), nil
}

// SubscriptionsIDDelete implements api.Handler.
func (h *OgenAdapter) SubscriptionsIDDelete(ctx context.Context, params api.SubscriptionsIDDeleteParams) (api.SubscriptionsIDDeleteRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	err := h.service.DeleteSubscription(ctx, params.ID)
	if err != nil {
		log.Error().Err(err).Str("subscription_id", params.ID.String()).Msg("Failed to delete subscription")
		return convertSubscriptionsIDDeleteError(err), nil
	}

	return &api.SubscriptionsIDDeleteNoContent{}, nil
}

// SubscriptionsSummaryTotalCostGet implements api.Handler.
func (h *OgenAdapter) SubscriptionsSummaryTotalCostGet(ctx context.Context, params api.SubscriptionsSummaryTotalCostGetParams) (api.SubscriptionsSummaryTotalCostGetRes, error) {
	log := logger.WithRequestID(getRequestID(ctx))

	domainReq := &ports.TotalCostRequest{
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		UserIDs:      params.UserIds,
		ServiceNames: params.ServiceNames,
	}

	result, err := h.service.GetTotalCost(ctx, domainReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate total cost")
		return convertSubscriptionsSummaryTotalCostGetError(err), nil
	}

	period := api.SubscriptionsSummaryTotalCostGetOKPeriod{}
	period.SetStartDate(api.NewOptString(result.Period.StartDate))
	period.SetEndDate(api.NewOptString(result.Period.EndDate))

	optPeriod := api.OptSubscriptionsSummaryTotalCostGetOKPeriod{}
	optPeriod.SetTo(period)

	filter := api.SubscriptionsSummaryTotalCostGetOKFilterCriteria{}

	stringIDs := make([]string, len(params.UserIds))
	for i := range params.UserIds {
		stringIDs[i] = params.UserIds[i].String()
	}

	filter.SetUserIds(stringIDs)
	filter.SetServiceNames(result.FilterCriteria.ServiceNames)

	optFilter := api.OptSubscriptionsSummaryTotalCostGetOKFilterCriteria{}
	optFilter.SetTo(filter)

	response := &api.SubscriptionsSummaryTotalCostGetOK{
		TotalCost:      api.NewOptInt(result.TotalCost),
		Period:         optPeriod,
		FilterCriteria: optFilter,
	}

	return response, nil
}

// Helper functions

func convertFilterParams(params api.SubscriptionsGetParams) ports.SubscriptionFilter {
	return ports.SubscriptionFilter{
		UserIDs:       params.UserIds,
		ServiceNames:  params.ServiceNames,
		StartDateFrom: getStringPtrFromOpt(params.StartDateFrom),
		StartDateTo:   getStringPtrFromOpt(params.StartDateTo),
	}
}

func convertSubscriptionToOgen(sub *domain.Subscription) *api.Subscription {
	if sub == nil {
		return nil
	}

	return &api.Subscription{
		ID:          api.NewOptUUID(sub.ID),
		ServiceName: api.NewOptString(sub.ServiceName),
		Price:       api.NewOptInt32(int32(sub.Price)),
		UserID:      api.NewOptUUID(sub.UserID),
		StartDate:   api.NewOptString(sub.StartDate),
		EndDate:     api.NewOptNilString(*sub.EndDate),
		CreatedAt:   api.NewOptDateTime(sub.CreatedAt),
		UpdatedAt:   api.NewOptDateTime(sub.UpdatedAt),
	}
}

func convertSubscriptionsToOgen(subscriptions []*domain.Subscription) []api.Subscription {
	result := make([]api.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = api.Subscription{
			ID:          api.NewOptUUID(sub.ID),
			ServiceName: api.NewOptString(sub.ServiceName),
			Price:       api.NewOptInt32(int32(sub.Price)),
			UserID:      api.NewOptUUID(sub.UserID),
			StartDate:   api.NewOptString(sub.StartDate),
			EndDate:     api.NewOptNilString(*sub.EndDate),
			CreatedAt:   api.NewOptDateTime(sub.CreatedAt),
			UpdatedAt:   api.NewOptDateTime(sub.UpdatedAt),
		}
	}
	return result
}

func convertPaginationToOgen(meta *ports.PaginationMetadata) api.OptPagination {
	if meta == nil {
		return api.OptPagination{}
	}
	return api.OptPagination{
		Value: api.Pagination{
			Page:  api.NewOptInt(meta.Page),
			Limit: api.NewOptInt(meta.Limit),
			Total: api.NewOptInt(meta.Total),
			Pages: api.NewOptInt(meta.TotalPages),
		},
	}
}

func convertError(err error) *api.ErrorStatusCode {
	// Проверка стандартных ошибок
	switch {
	case errors.Is(err, ports.ErrSubscriptionNotFound):
		return &api.ErrorStatusCode{
			StatusCode: 404,
			Response: api.Error{
				Error:   "not_found",
				Message: err.Error(),
			},
		}
	case errors.Is(err, ports.ErrInvalidDateformat),
		errors.Is(err, ports.ErrInvalidUUID),
		errors.Is(err, ports.ErrInvalidPrice),
		errors.Is(err, ports.ErrStartDateAfterEndDate),
		errors.Is(err, ports.ErrInvalidDateRange):
		return &api.ErrorStatusCode{
			StatusCode: 400,
			Response: api.Error{
				Error:   "validation_error",
				Message: err.Error(),
			},
		}
	case errors.Is(err, ports.ErrDuplicateSubscription):
		return &api.ErrorStatusCode{
			StatusCode: 409,
			Response: api.Error{
				Error:   "conflict",
				Message: err.Error(),
			},
		}
	default:
		// Для кастомных DomainError
		var domainErr *ports.DomainError
		if errors.As(err, &domainErr) {
			return &api.ErrorStatusCode{
				StatusCode: getStatusCodeFromDomainError(domainErr),
				Response: api.Error{
					Error:   domainErr.Code,
					Message: domainErr.Message,
				},
			}
		}

		// Любая другая ошибка
		return &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Error:   "internal_error",
				Message: "Internal server error",
			},
		}
	}
}

// Error conversion functions for each operation

func convertSubscriptionsPostError(err error) *api.SubscriptionsPostBadRequest {
	errorResponse := createErrorResponse(err)
	// Просто приводим тип, поскольку SubscriptionsPostBadRequest это алиас к Error
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
	case errors.Is(err, ports.ErrSubscriptionNotFound):
		return 404
	case errors.Is(err, ports.ErrInvalidDateformat),
		errors.Is(err, ports.ErrInvalidUUID),
		errors.Is(err, ports.ErrInvalidPrice),
		errors.Is(err, ports.ErrStartDateAfterEndDate),
		errors.Is(err, ports.ErrInvalidDateRange):
		return 400
	case errors.Is(err, ports.ErrDuplicateSubscription):
		return 409
	default:
		return 500
	}
}

func getErrorCode(err error) string {
	switch {
	case errors.Is(err, ports.ErrSubscriptionNotFound):
		return "not_found"
	case errors.Is(err, ports.ErrInvalidDateformat):
		return "invalid_date_format"
	case errors.Is(err, ports.ErrInvalidUUID):
		return "invalid_uuid"
	case errors.Is(err, ports.ErrInvalidPrice):
		return "invalid_price"
	case errors.Is(err, ports.ErrStartDateAfterEndDate):
		return "invalid_date_range"
	case errors.Is(err, ports.ErrInvalidDateRange):
		return "invalid_date_range"
	case errors.Is(err, ports.ErrDuplicateSubscription):
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
func getStatusCodeFromDomainError(err *ports.DomainError) int {
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

// Utility functions for working with ogen optional types

func getStringPtrFromOptNil(opt api.OptNilString) *string {
	if !opt.Set {
		return nil
	}
	return &opt.Value
}

func getStringPtrFromOpt(opt api.OptString) *string {
	if !opt.Set {
		return nil
	}
	return &opt.Value
}

func getIntPtrFromOpt(opt api.OptInt32) *int {
	if !opt.Set {
		return nil
	}
	value := int(opt.Value)
	return &value
}

func getStringPtrFromUUIDOpt(opt api.OptUUID) *uuid.UUID {
	if !opt.Set {
		return nil
	}
	value := opt.Value
	return &value
}

func getIntOrDefault(get func() (int, bool), defaultValue int) int {
	if value, set := get(); set {
		return value
	}
	return defaultValue
}

func getRequestID(ctx context.Context) string {
	// Implement request ID extraction from context
	return "unknown"
}
