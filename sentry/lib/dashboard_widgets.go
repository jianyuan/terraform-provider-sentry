package sentry

import (
	"context"
	"fmt"
)

// DashboardWidget represents a Dashboard Widget.
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/rest_framework/dashboard.py#L230-L243
type DashboardWidget struct {
	ID          *string                 `json:"id,omitempty"`
	Title       *string                 `json:"title,omitempty"`
	DisplayType *string                 `json:"displayType,omitempty"`
	Interval    *string                 `json:"interval,omitempty"`
	Queries     []*DashboardWidgetQuery `json:"queries,omitempty"`
	WidgetType  *string                 `json:"widgetType,omitempty"`
	Limit       *int                    `json:"limit,omitempty"`
	Layout      *DashboardWidgetLayout  `json:"layout,omitempty"`
}

type DashboardWidgetLayout struct {
	X    *int `json:"x,omitempty"`
	Y    *int `json:"y,omitempty"`
	W    *int `json:"w,omitempty"`
	H    *int `json:"h,omitempty"`
	MinH *int `json:"minH,omitempty"`
}

type DashboardWidgetQuery struct {
	ID           *string  `json:"id,omitempty"`
	Fields       []string `json:"fields,omitempty"`
	Aggregates   []string `json:"aggregates,omitempty"`
	Columns      []string `json:"columns,omitempty"`
	FieldAliases []string `json:"fieldAliases,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Conditions   *string  `json:"conditions,omitempty"`
	OrderBy      *string  `json:"orderby,omitempty"`
}

// DashboardWidgetsService provides methods for accessing Sentry dashboard widget API endpoints.
type DashboardWidgetsService service

type DashboardWidgetErrors map[string][]string

// Validate a dashboard widget configuration.
func (s *DashboardWidgetsService) Validate(ctx context.Context, organizationSlug string, widget *DashboardWidget) (DashboardWidgetErrors, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/dashboards/widgets/", organizationSlug)

	req, err := s.client.NewRequest("POST", u, widget)
	if err != nil {
		return nil, nil, err
	}

	widgetErrors := make(DashboardWidgetErrors)
	resp, err := s.client.Do(ctx, req, &widgetErrors)
	if err != nil {
		return nil, resp, err
	}
	if len(widgetErrors) == 0 {
		return nil, resp, err
	}
	return widgetErrors, resp, err
}
