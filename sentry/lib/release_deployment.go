package sentry

import (
	"context"
	"fmt"
	"time"
)

type ReleaseDeploymentsService service

type ReleaseDeployment struct {
	ID           string     `json:"id"`
	Name         *string    `json:"name,omitempty"`
	Environment  string     `json:"environment,omitempty"`
	URL          *string    `json:"url,omitempty"`
	Projects     []string   `json:"projects,omitempty"`
	DateStarted  *time.Time `json:"dateStarted,omitempty"`
	DateFinished *time.Time `json:"dateFinished,omitempty"`
}

// Get a Release Deploy for a project.
func (s *ReleaseDeploymentsService) Get(ctx context.Context, organizationSlug string, version string, deployID string) (*ReleaseDeployment, *Response, error) {

	lastCursor := ""

	// Search for the deployment ID by using the list endpoint. When we have
	// found the first match return immediately
	for {
		params := ListCursorParams{
			Cursor: lastCursor,
		}

		u := fmt.Sprintf("0/organizations/%v/releases/%s/deploys/", organizationSlug, version)
		u, err := addQuery(u, params)
		if err != nil {
			return nil, nil, err
		}

		req, err := s.client.NewRequest("GET", u, nil)
		if err != nil {
			return nil, nil, err
		}

		deployments := new([]ReleaseDeployment)
		resp, err := s.client.Do(ctx, req, deployments)
		if err != nil {
			return nil, resp, err
		}

		for i := range *deployments {
			d := (*deployments)[i]
			if d.ID == deployID {
				return &d, resp, nil
			}
		}

		// No matches in the current page and no further pages to check
		if resp.Cursor == "" {
			return nil, resp, nil
		}
		lastCursor = resp.Cursor
	}
}

// Create a new Release Deploy to a project.
func (s *ReleaseDeploymentsService) Create(ctx context.Context, organizationSlug string, version string, params *ReleaseDeployment) (*ReleaseDeployment, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/releases/%s/deploys/", organizationSlug, version)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	deploy := new(ReleaseDeployment)
	resp, err := s.client.Do(ctx, req, deploy)
	if err != nil {
		return nil, resp, err
	}

	return deploy, resp, nil
}
