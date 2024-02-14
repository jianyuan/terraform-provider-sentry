package sentry

import (
	"context"
	"fmt"
	"time"
)

// Dashboard represents a Dashboard.
type Dashboard struct {
	ID          *string            `json:"id,omitempty"`
	Title       *string            `json:"title,omitempty"`
	DateCreated *time.Time         `json:"dateCreated,omitempty"`
	Widgets     []*DashboardWidget `json:"widgets,omitempty"`
}

// DashboardsService provides methods for accessing Sentry dashboard API endpoints.
type DashboardsService service

// List dashboards in an organization.
func (s *DashboardsService) List(ctx context.Context, organizationSlug string, params *ListCursorParams) ([]*Dashboard, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/dashboards/", organizationSlug)
	u, err := addQuery(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var dashboards []*Dashboard
	resp, err := s.client.Do(ctx, req, &dashboards)
	if err != nil {
		return nil, resp, err
	}
	return dashboards, resp, nil
}

// Get details on a dashboard.
func (s *DashboardsService) Get(ctx context.Context, organizationSlug string, id string) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/dashboards/%v/", organizationSlug, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	dashboard := new(Dashboard)
	resp, err := s.client.Do(ctx, req, dashboard)
	if err != nil {
		return nil, resp, err
	}
	return dashboard, resp, nil
}

// Create a dashboard.
func (s *DashboardsService) Create(ctx context.Context, organizationSlug string, params *Dashboard) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/dashboards/", organizationSlug)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	dashboard := new(Dashboard)
	resp, err := s.client.Do(ctx, req, dashboard)
	if err != nil {
		return nil, resp, err
	}
	return dashboard, resp, nil
}

// Update a dashboard.
func (s *DashboardsService) Update(ctx context.Context, organizationSlug string, id string, params *Dashboard) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/dashboards/%v/", organizationSlug, id)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	dashboard := new(Dashboard)
	resp, err := s.client.Do(ctx, req, dashboard)
	if err != nil {
		return nil, resp, err
	}
	return dashboard, resp, nil
}

// Delete a dashboard.
func (s *DashboardsService) Delete(ctx context.Context, organizationSlug string, id string) (*Response, error) {
	u := fmt.Sprintf("0/organizations/%v/dashboards/%v/", organizationSlug, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
