package sentry

import (
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// OrganizationAvatar represents a Sentry organization's avatar.
type OrganizationAvatar struct {
	UUID *string `json:"avatarUuid"`
	Type string  `json:"avatarType"`
}

// OrganizationQuota represents a Sentry organization's quota.
type OrganizationQuota struct {
	MaxRate         int `json:"maxRate"`
	MaxRateInterval int `json:"maxRateInterval"`
	AccountLimit    int `json:"accountLimit"`
	ProjectLimit    int `json:"projectLimit"`
}

// OrganizationAvailableRole represents a Sentry organization's available role.
type OrganizationAvailableRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Organization represents a Sentry organization.
// Based on https://github.com/getsentry/sentry/blob/cc81fff31d4f2c9cede14ce9c479d6f4f78c5e5b/src/sentry/api/endpoints/organization_details.py#L66.
type Organization struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`

	DateCreated time.Time `json:"dateCreated"`

	Quota OrganizationQuota `json:"quota"`

	Access                []string `json:"access"`
	Features              []string `json:"features"`
	PendingAccessRequests int      `json:"pendingAccessRequests"`

	IsDefault        bool                        `json:"isDefault"`
	DefaultRole      string                      `json:"defaultRole"`
	AvailableRoles   []OrganizationAvailableRole `json:"availableRoles"`
	AccountRateLimit int                         `json:"accountRateLimit"`
	ProjectRateLimit int                         `json:"projectRateLimit"`

	Avatar OrganizationAvatar `json:"avatar"`

	OpenMembership       bool     `json:"openMembership"`
	AllowSharedIssues    bool     `json:"allowSharedIssues"`
	EnhancedPrivacy      bool     `json:"enhancedPrivacy"`
	DataScrubber         bool     `json:"dataScrubber"`
	DataScrubberDefaults bool     `json:"dataScrubberDefaults"`
	SensitiveFields      []string `json:"sensitiveFields"`
	SafeFields           []string `json:"safeFields"`
	ScrubIPAddresses     bool     `json:"scrubIPAddresses"`
	IsEarlyAdopter       bool     `json:"isEarlyAdopter"`

	// TODO: teams
	// TODO: onboardingTasks
}

// OrganizationService provides methods for accessing Sentry organization API endpoints.
// https://docs.sentry.io/api/organizations/
type OrganizationService struct {
	sling *sling.Sling
}

func newOrganizationService(sling *sling.Sling) *OrganizationService {
	return &OrganizationService{
		sling: sling.Path("organizations/"),
	}
}

// ListOrganizationParams are the parameters for OrganizationService.List.
type ListOrganizationParams struct {
	Cursor string `url:"cursor,omitempty"`
}

// List organizations available to the authenticated session.
// https://docs.sentry.io/api/organizations/get-organization-index/
func (s *OrganizationService) List(params *ListOrganizationParams) ([]Organization, *http.Response, error) {
	organizations := new([]Organization)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("").QueryStruct(params).Receive(organizations, apiError)
	return *organizations, resp, relevantError(err, *apiError)
}

// CreateOrganizationParams are the parameters for OrganizationService.Create.
type CreateOrganizationParams struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// Get a Sentry organization.
// https://docs.sentry.io/api/organizations/get-organization-details/
func (s *OrganizationService) Get(slug string) (*Organization, *http.Response, error) {
	org := new(Organization)
	apiError := new(APIError)
	resp, err := s.sling.New().Get(slug+"/").Receive(org, apiError)
	return org, resp, relevantError(err, *apiError)
}

// Create a new Sentry organization.
// https://docs.sentry.io/api/organizations/post-organization-index/
func (s *OrganizationService) Create(params *CreateOrganizationParams) (*Organization, *http.Response, error) {
	org := new(Organization)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("").BodyJSON(params).Receive(org, apiError)
	return org, resp, relevantError(err, *apiError)
}

// UpdateOrganizationParams are the parameters for OrganizationService.Update.
type UpdateOrganizationParams struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// Update a Sentry organization.
// https://docs.sentry.io/api/organizations/put-organization-details/
func (s *OrganizationService) Update(slug string, params *UpdateOrganizationParams) (*Organization, *http.Response, error) {
	org := new(Organization)
	apiError := new(APIError)
	resp, err := s.sling.New().Put(slug+"/").BodyJSON(params).Receive(org, apiError)
	return org, resp, relevantError(err, *apiError)
}

// Delete a Sentry organization.
func (s *OrganizationService) Delete(slug string) (*http.Response, error) {
	apiError := new(APIError)
	resp, err := s.sling.New().Delete(slug+"/").Receive(nil, apiError)
	return resp, relevantError(err, *apiError)
}
