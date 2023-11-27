package sentry

import (
	"context"
	"fmt"
	"time"
)

// OrganizationStatus represents a Sentry organization's status.
type OrganizationStatus struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

// OrganizationQuota represents a Sentry organization's quota.
type OrganizationQuota struct {
	MaxRate         *int `json:"maxRate"`
	MaxRateInterval *int `json:"maxRateInterval"`
	AccountLimit    *int `json:"accountLimit"`
	ProjectLimit    *int `json:"projectLimit"`
}

// OrganizationAvailableRole represents a Sentry organization's available role.
type OrganizationAvailableRole struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

// Organization represents detailed information about a Sentry organization.
// Based on https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/models/organization.py#L263-L288
type Organization struct {
	// Basic
	ID                       *string             `json:"id,omitempty"`
	Slug                     *string             `json:"slug,omitempty"`
	Status                   *OrganizationStatus `json:"status,omitempty"`
	Name                     *string             `json:"name,omitempty"`
	DateCreated              *time.Time          `json:"dateCreated,omitempty"`
	IsEarlyAdopter           *bool               `json:"isEarlyAdopter,omitempty"`
	Require2FA               *bool               `json:"require2FA,omitempty"`
	RequireEmailVerification *bool               `json:"requireEmailVerification,omitempty"`
	Avatar                   *Avatar             `json:"avatar,omitempty"`
	Features                 []string            `json:"features,omitempty"`

	// Detailed
	// TODO: experiments
	Quota                *OrganizationQuota          `json:"quota,omitempty"`
	IsDefault            *bool                       `json:"isDefault,omitempty"`
	DefaultRole          *string                     `json:"defaultRole,omitempty"`
	AvailableRoles       []OrganizationAvailableRole `json:"availableRoles,omitempty"`
	OpenMembership       *bool                       `json:"openMembership,omitempty"`
	AllowSharedIssues    *bool                       `json:"allowSharedIssues,omitempty"`
	EnhancedPrivacy      *bool                       `json:"enhancedPrivacy,omitempty"`
	DataScrubber         *bool                       `json:"dataScrubber,omitempty"`
	DataScrubberDefaults *bool                       `json:"dataScrubberDefaults,omitempty"`
	SensitiveFields      []string                    `json:"sensitiveFields,omitempty"`
	SafeFields           []string                    `json:"safeFields,omitempty"`
	StoreCrashReports    *int                        `json:"storeCrashReports,omitempty"`
	AttachmentsRole      *string                     `json:"attachmentsRole,omitempty"`
	DebugFilesRole       *string                     `json:"debugFilesRole,omitempty"`
	EventsMemberAdmin    *bool                       `json:"eventsMemberAdmin,omitempty"`
	AlertsMemberWrite    *bool                       `json:"alertsMemberWrite,omitempty"`
	ScrubIPAddresses     *bool                       `json:"scrubIPAddresses,omitempty"`
	ScrapeJavaScript     *bool                       `json:"scrapeJavaScript,omitempty"`
	AllowJoinRequests    *bool                       `json:"allowJoinRequests,omitempty"`
	RelayPiiConfig       *string                     `json:"relayPiiConfig,omitempty"`
	// TODO: trustedRelays
	Access                []string `json:"access,omitempty"`
	Role                  *string  `json:"role,omitempty"`
	PendingAccessRequests *int     `json:"pendingAccessRequests,omitempty"`
	// TODO: onboardingTasks
}

// OrganizationsService provides methods for accessing Sentry organization API endpoints.
// https://docs.sentry.io/api/organizations/
type OrganizationsService service

// List organizations available to the authenticated session.
// https://docs.sentry.io/api/organizations/list-your-organizations/
func (s *OrganizationsService) List(ctx context.Context, params *ListCursorParams) ([]*Organization, *Response, error) {
	u, err := addQuery("0/organizations/", params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := []*Organization{}
	resp, err := s.client.Do(ctx, req, &orgs)
	if err != nil {
		return nil, resp, err
	}
	return orgs, resp, nil
}

// Get a Sentry organization.
// https://docs.sentry.io/api/organizations/retrieve-an-organization/
func (s *OrganizationsService) Get(ctx context.Context, slug string) (*Organization, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/", slug)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	org := new(Organization)
	resp, err := s.client.Do(ctx, req, org)
	if err != nil {
		return nil, resp, err
	}
	return org, resp, nil
}

// CreateOrganizationParams are the parameters for OrganizationService.Create.
type CreateOrganizationParams struct {
	Name       *string `json:"name,omitempty"`
	Slug       *string `json:"slug,omitempty"`
	AgreeTerms *bool   `json:"agreeTerms,omitempty"`
}

// Create a new Sentry organization.
func (s *OrganizationsService) Create(ctx context.Context, params *CreateOrganizationParams) (*Organization, *Response, error) {
	u := "0/organizations/"
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	org := new(Organization)
	resp, err := s.client.Do(ctx, req, org)
	if err != nil {
		return nil, resp, err
	}
	return org, resp, nil
}

// UpdateOrganizationParams are the parameters for OrganizationService.Update.
type UpdateOrganizationParams struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// Update a Sentry organization.
// https://docs.sentry.io/api/organizations/update-an-organization/
func (s *OrganizationsService) Update(ctx context.Context, slug string, params *UpdateOrganizationParams) (*Organization, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/", slug)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	org := new(Organization)
	resp, err := s.client.Do(ctx, req, org)
	if err != nil {
		return nil, resp, err
	}
	return org, resp, nil
}

// Delete a Sentry organization.
func (s *OrganizationsService) Delete(ctx context.Context, slug string) (*Response, error) {
	u := fmt.Sprintf("0/organizations/%v/", slug)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
