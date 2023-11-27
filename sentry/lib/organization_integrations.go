package sentry

import (
	"context"
	"fmt"
	"time"
)

// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/serializers/models/integration.py#L22
type OrganizationIntegrationProvider struct {
	Key        string   `json:"key"`
	Slug       string   `json:"slug"`
	Name       string   `json:"name"`
	CanAdd     bool     `json:"canAdd"`
	CanDisable bool     `json:"canDisable"`
	Features   []string `json:"features"`
}

// IntegrationConfigData for defining integration-specific configuration data.
type IntegrationConfigData map[string]interface{}

// OrganizationIntegration represents an integration added for the organization.
// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/serializers/models/integration.py#L93
type OrganizationIntegration struct {
	// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/serializers/models/integration.py#L35
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	Icon        *string                         `json:"icon"`
	DomainName  string                          `json:"domainName"`
	AccountType *string                         `json:"accountType"`
	Scopes      []string                        `json:"scopes"`
	Status      string                          `json:"status"`
	Provider    OrganizationIntegrationProvider `json:"provider"`

	// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/serializers/models/integration.py#L138
	ConfigData                    *IntegrationConfigData `json:"configData"`
	ExternalId                    string                 `json:"externalId"`
	OrganizationId                int                    `json:"organizationId"`
	OrganizationIntegrationStatus string                 `json:"organizationIntegrationStatus"`
	GracePeriodEnd                *time.Time             `json:"gracePeriodEnd"`
}

// OrganizationIntegrationsService provides methods for accessing Sentry organization integrations API endpoints.
// Paths: https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/urls.py#L1236-L1245
// Endpoints: https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/endpoints/integrations/organization_integrations/index.py
// Endpoints: https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/endpoints/integrations/organization_integrations/details.py
type OrganizationIntegrationsService service

type ListOrganizationIntegrationsParams struct {
	ListCursorParams
	ProviderKey string `url:"provider_key,omitempty"`
}

// List organization integrations.
func (s *OrganizationIntegrationsService) List(ctx context.Context, organizationSlug string, params *ListOrganizationIntegrationsParams) ([]*OrganizationIntegration, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/integrations/", organizationSlug)
	u, err := addQuery(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	integrations := []*OrganizationIntegration{}
	resp, err := s.client.Do(ctx, req, &integrations)
	if err != nil {
		return nil, resp, err
	}
	return integrations, resp, nil
}

// Get organization integration details.
func (s *OrganizationIntegrationsService) Get(ctx context.Context, organizationSlug string, integrationID string) (*OrganizationIntegration, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/integrations/%v/", organizationSlug, integrationID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	integration := new(OrganizationIntegration)
	resp, err := s.client.Do(ctx, req, integration)
	if err != nil {
		return nil, resp, err
	}
	return integration, resp, nil
}

type UpdateConfigOrganizationIntegrationsParams = IntegrationConfigData

// UpdateConfig - update configData for organization integration.
// https://github.com/getsentry/sentry/blob/22.7.0/src/sentry/api/endpoints/integrations/organization_integrations/details.py#L94-L102
func (s *OrganizationIntegrationsService) UpdateConfig(ctx context.Context, organizationSlug string, integrationID string, params *UpdateConfigOrganizationIntegrationsParams) (*Response, error) {
	u := fmt.Sprintf("0/organizations/%v/integrations/%v/", organizationSlug, integrationID)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
