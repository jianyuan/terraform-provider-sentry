package sentry

import (
	"context"
	"fmt"
	"time"
)

// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/serializers/models/repository.py#L12-L17
type OrganizationRepositoryProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrganizationRepositories represents
// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/serializers/models/repository.py
type OrganizationRepository struct {
	ID            string                         `json:"id"`
	Name          string                         `json:"name"`
	Url           string                         `json:"url"`
	Provider      OrganizationRepositoryProvider `json:"provider"`
	Status        string                         `json:"status"`
	DateCreated   time.Time                      `json:"dateCreated"`
	IntegrationId string                         `json:"integrationId"`
	ExternalSlug  string                         `json:"externalSlug"`
}

// OrganizationRepositoriesService provides methods for accessing Sentry organization repositories API endpoints.
// Paths: https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/urls.py#L1385-L1394
// Endpoints: https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/endpoints/organization_repositories.py
// Endpoints: https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/endpoints/organization_repository_details.py
type OrganizationRepositoriesService service

type ListOrganizationRepositoriesParams struct {
	ListCursorParams
	// omitting status defaults to only active.
	// sending empty string shows everything, which is a more reasonable default.
	Status string `url:"status"`
	Query  string `url:"query,omitempty"`
}

// List organization integrations.
func (s *OrganizationRepositoriesService) List(ctx context.Context, organizationSlug string, params *ListOrganizationRepositoriesParams) ([]*OrganizationRepository, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/repos/", organizationSlug)
	u, err := addQuery(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repos := []*OrganizationRepository{}
	resp, err := s.client.Do(ctx, req, &repos)
	if err != nil {
		return nil, resp, err
	}
	return repos, resp, nil
}

// Fields are different for different providers
type CreateOrganizationRepositoryParams map[string]interface{}

func (s *OrganizationRepositoriesService) Create(ctx context.Context, organizationSlug string, params CreateOrganizationRepositoryParams) (*OrganizationRepository, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/repos/", organizationSlug)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	repo := new(OrganizationRepository)
	resp, err := s.client.Do(ctx, req, repo)
	if err != nil {
		return nil, resp, err
	}
	return repo, resp, nil
}

func (s *OrganizationRepositoriesService) Delete(ctx context.Context, organizationSlug string, repoID string) (*OrganizationRepository, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/repos/%v/", organizationSlug, repoID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repo := new(OrganizationRepository)
	resp, err := s.client.Do(ctx, req, repo)
	if err != nil {
		return nil, resp, err
	}
	return repo, resp, nil
}
