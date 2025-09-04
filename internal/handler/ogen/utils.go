package ogen

import (
	"subscription/core/domain"
	"subscription/core/ports"
	api "subscription/internal/api/generated"
)

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
