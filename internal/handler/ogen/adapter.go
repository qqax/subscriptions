package ogen

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"subscription/core/domain"
	"subscription/core/ports"
	api "subscription/internal/api/generated" // сгенерированный ogen код
	"subscription/internal/logger"
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
		log.Error().Err(err).Msg("Failed to list server")
		return convertSubscriptionsGetError(err), nil
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
		return convertSubscriptionsPostError(domain.ErrInvalidUUID), nil
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

type contextKey string

const (
	// RequestIDKey стандартный ключ для request ID
	RequestIDKey contextKey = "request_id"
	// XRequestIDKey ключ для заголовка X-Request-ID
	XRequestIDKey contextKey = "x-request-id"
)

func getRequestID(ctx context.Context) string {
	// Пробуем получить request ID из стандартных ключей
	keys := []contextKey{RequestIDKey, XRequestIDKey}
	for _, key := range keys {
		if value, ok := ctx.Value(key).(string); ok && value != "" {
			return value
		}
	}

	// Пробуем получить из HTTP заголовков (если контекст содержит http.Request)
	if req, ok := ctx.Value("http_request").(*http.Request); ok {
		if requestID := req.Header.Get("X-Request-ID"); requestID != "" {
			return requestID
		}
	}

	// Пробуем получить из gRPC метаданных (если применимо)
	//if md, ok := metadata.FromIncomingContext(ctx); ok {
	//	if values := md.Get("x-request-id"); len(values) > 0 {
	//		return values[0]
	//	}
	//}

	return "unknown"
}
