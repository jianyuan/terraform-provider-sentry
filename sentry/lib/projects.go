package sentry

import (
	"context"
	"fmt"
	"time"
)

// Project represents a Sentry project.
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/models/project.py
type Project struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`

	IsPublic     bool   `json:"isPublic"`
	IsBookmarked bool   `json:"isBookmarked"`
	Color        string `json:"color"`

	DateCreated time.Time `json:"dateCreated"`
	FirstEvent  time.Time `json:"firstEvent"`

	Features []string `json:"features"`
	Status   string   `json:"status"`
	Platform string   `json:"platform"`

	IsInternal bool `json:"isInternal"`
	IsMember   bool `json:"isMember"`
	HasAccess  bool `json:"hasAccess"`

	Avatar Avatar `json:"avatar"`

	// TODO: latestRelease
	Options map[string]interface{} `json:"options"`

	DigestsMinDelay      int      `json:"digestsMinDelay"`
	DigestsMaxDelay      int      `json:"digestsMaxDelay"`
	SubjectPrefix        string   `json:"subjectPrefix"`
	AllowedDomains       []string `json:"allowedDomains"`
	ResolveAge           int      `json:"resolveAge"`
	DataScrubber         bool     `json:"dataScrubber"`
	DataScrubberDefaults bool     `json:"dataScrubberDefaults"`
	FingerprintingRules  string   `json:"fingerprintingRules"`
	GroupingEnhancements string   `json:"groupingEnhancements"`
	SafeFields           []string `json:"safeFields"`
	SensitiveFields      []string `json:"sensitiveFields"`
	SubjectTemplate      string   `json:"subjectTemplate"`
	SecurityToken        string   `json:"securityToken"`
	SecurityTokenHeader  *string  `json:"securityTokenHeader"`
	VerifySSL            bool     `json:"verifySSL"`
	ScrubIPAddresses     bool     `json:"scrubIPAddresses"`
	ScrapeJavaScript     bool     `json:"scrapeJavaScript"`

	Organization Organization `json:"organization"`
	// TODO: plugins
	// TODO: platforms
	ProcessingIssues int `json:"processingIssues"`
	// TODO: defaultEnvironment

	Team  Team   `json:"team"`
	Teams []Team `json:"teams"`
}

// ProjectSummary represents the summary of a Sentry project.
type ProjectSummary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	IsBookmarked bool   `json:"isBookmarked"`
	IsMember     bool   `json:"isMember"`
	HasAccess    bool   `json:"hasAccess"`

	DateCreated time.Time `json:"dateCreated"`
	FirstEvent  time.Time `json:"firstEvent"`

	Platform  *string  `json:"platform"`
	Platforms []string `json:"platforms"`

	Team  *ProjectSummaryTeam  `json:"team"`
	Teams []ProjectSummaryTeam `json:"teams"`
	// TODO: deploys
}

// ProjectSummaryTeam represents a team in a ProjectSummary.
type ProjectSummaryTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProjectsService provides methods for accessing Sentry project API endpoints.
// https://docs.sentry.io/api/projects/
type ProjectsService service

// List projects available.
// https://docs.sentry.io/api/projects/list-your-projects/
func (s *ProjectsService) List(ctx context.Context) ([]*Project, *Response, error) {
	u := "0/projects/"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projects := []*Project{}
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}

// Get details on an individual project.
// https://docs.sentry.io/api/projects/retrieve-a-project/
func (s *ProjectsService) Get(ctx context.Context, organizationSlug string, slug string) (*Project, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/", organizationSlug, slug)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// CreateProjectParams are the parameters for ProjectService.Create.
type CreateProjectParams struct {
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Platform string `json:"platform,omitempty"`
}

// Create a new project bound to a team.
func (s *ProjectsService) Create(ctx context.Context, organizationSlug string, teamSlug string, params *CreateProjectParams) (*Project, *Response, error) {
	u := fmt.Sprintf("0/teams/%v/%v/projects/", organizationSlug, teamSlug)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// UpdateProjectParams are the parameters for ProjectService.Update.
type UpdateProjectParams struct {
	Name                 string                 `json:"name,omitempty"`
	Slug                 string                 `json:"slug,omitempty"`
	Platform             string                 `json:"platform,omitempty"`
	IsBookmarked         *bool                  `json:"isBookmarked,omitempty"`
	DigestsMinDelay      *int                   `json:"digestsMinDelay,omitempty"`
	DigestsMaxDelay      *int                   `json:"digestsMaxDelay,omitempty"`
	ResolveAge           *int                   `json:"resolveAge,omitempty"`
	Options              map[string]interface{} `json:"options,omitempty"`
	AllowedDomains       []string               `json:"allowedDomains,omitempty"`
	FingerprintingRules  string                 `json:"fingerprintingRules,omitempty"`
	GroupingEnhancements string                 `json:"groupingEnhancements,omitempty"`
}

// Update various attributes and configurable settings for a given project.
// https://docs.sentry.io/api/projects/update-a-project/
func (s *ProjectsService) Update(ctx context.Context, organizationSlug string, slug string, params *UpdateProjectParams) (*Project, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/", organizationSlug, slug)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// Delete a project.
// https://docs.sentry.io/api/projects/delete-a-project/
func (s *ProjectsService) Delete(ctx context.Context, organizationSlug string, slug string) (*Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/", organizationSlug, slug)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddTeam add a team to a project.
func (s *ProjectsService) AddTeam(ctx context.Context, organizationSlug string, slug string, teamSlug string) (*Project, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/teams/%v/", organizationSlug, slug, teamSlug)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, nil
}

// RemoveTeam remove a team from a project.
func (s *ProjectsService) RemoveTeam(ctx context.Context, organizationSlug string, slug string, teamSlug string) (*Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/teams/%v/", organizationSlug, slug, teamSlug)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
