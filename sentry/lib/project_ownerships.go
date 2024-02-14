package sentry

import (
	"context"
	"fmt"
	"time"
)

// https://github.com/getsentry/sentry/blob/master/src/sentry/api/serializers/models/projectownership.py
type ProjectOwnership struct {
	Raw                string    `json:"raw"`
	FallThrough        bool      `json:"fallthrough"`
	DateCreated        time.Time `json:"dateCreated"`
	LastUpdated        time.Time `json:"lastUpdated"`
	IsActive           bool      `json:"isActive"`
	AutoAssignment     bool      `json:"autoAssignment"`
	CodeownersAutoSync *bool     `json:"codeownersAutoSync,omitempty"`
}

// ProjectOwnershipsService provides methods for accessing Sentry project
// ownership API endpoints.
type ProjectOwnershipsService service

// Get details on a project's ownership configuration.
func (s *ProjectOwnershipsService) Get(ctx context.Context, organizationSlug string, projectSlug string) (*ProjectOwnership, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/ownership/", organizationSlug, projectSlug)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	owner := new(ProjectOwnership)
	resp, err := s.client.Do(ctx, req, owner)
	if err != nil {
		return nil, resp, err
	}
	return owner, resp, nil
}

// CreateProjectParams are the parameters for ProjectOwnershipService.Update.
type UpdateProjectOwnershipParams struct {
	Raw                string `json:"raw,omitempty"`
	FallThrough        *bool  `json:"fallthrough,omitempty"`
	AutoAssignment     *bool  `json:"autoAssignment,omitempty"`
	CodeownersAutoSync *bool  `json:"codeownersAutoSync,omitempty"`
}

// Update a Project's Ownership configuration
func (s *ProjectOwnershipsService) Update(ctx context.Context, organizationSlug string, projectSlug string, params *UpdateProjectOwnershipParams) (*ProjectOwnership, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/ownership/", organizationSlug, projectSlug)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	owner := new(ProjectOwnership)
	resp, err := s.client.Do(ctx, req, owner)
	if err != nil {
		return nil, resp, err
	}
	return owner, resp, nil
}
