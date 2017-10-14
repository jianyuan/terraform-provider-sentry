package sentry

import (
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Project represents a Sentry project.
// Based on https://github.com/getsentry/sentry/blob/cc81fff31d4f2c9cede14ce9c479d6f4f78c5e5b/src/sentry/api/serializers/models/project.py#L137.
type Project struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`

	DateCreated time.Time `json:"dateCreated"`
	FirstEvent  time.Time `json:"firstEvent"`

	IsPublic     bool     `json:"isPublic"`
	IsBookmarked bool     `json:"isBookmarked"`
	CallSign     string   `json:"callSign"`
	Color        string   `json:"color"`
	Features     []string `json:"features"`
	Status       string   `json:"status"`

	// TODO: latestRelease
	Options         map[string]interface{} `json:"options"`
	DigestsMinDelay int                    `json:"digestsMinDelay"`
	DigestsMaxDelay int                    `json:"digestsMaxDelay"`
	SubjectPrefix   string                 `json:"subjectPrefix"`
	SubjectTemplate string                 `json:"subjectTemplate"`
	// TODO: plugins
	// TODO: platforms
	ProcessingIssues int `json:"processingIssues"`
	// TODO: defaultEnvironment

	Team         Team         `json:"team"`
	Organization Organization `json:"organization"`
}

// ProjectService provides methods for accessing Sentry project API endpoints.
// https://docs.sentry.io/api/projects/
type ProjectService struct {
	sling *sling.Sling
}

func newProjectService(sling *sling.Sling) *ProjectService {
	return &ProjectService{
		sling: sling,
	}
}

// List projects available.
// https://docs.sentry.io/api/projects/get-project-index/
func (s *ProjectService) List() ([]Project, *http.Response, error) {
	projects := new([]Project)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("projects/").Receive(projects, apiError)
	return *projects, resp, relevantError(err, *apiError)
}

// Get details on an individual project.
// https://docs.sentry.io/api/projects/get-project-details/
func (s *ProjectService) Get(organizationSlug string, slug string) (*Project, *http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("projects/"+organizationSlug+"/"+slug+"/").Receive(project, apiError)
	return project, resp, relevantError(err, *apiError)
}

// CreateProjectParams are the parameters for ProjectService.Create.
type CreateProjectParams struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// Create a new project bound to a team.
// https://docs.sentry.io/api/teams/post-team-project-index/
func (s *ProjectService) Create(organizationSlug string, teamSlug string, params *CreateProjectParams) (*Project, *http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("teams/"+organizationSlug+"/"+teamSlug+"/projects/").BodyJSON(params).Receive(project, apiError)
	return project, resp, relevantError(err, *apiError)
}

// UpdateProjectParams are the parameters for ProjectService.Update.
type UpdateProjectParams struct {
	Name            string                 `json:"name,omitempty"`
	Slug            string                 `json:"slug,omitempty"`
	IsBookmarked    *bool                  `json:"isBookmarked,omitempty"`
	DigestsMinDelay *int                   `json:"digestsMinDelay,omitempty"`
	DigestsMaxDelay *int                   `json:"digestsMaxDelay,omitempty"`
	Options         map[string]interface{} `json:"options,omitempty"`
}

// Update various attributes and configurable settings for a given project.
// https://docs.sentry.io/api/projects/put-project-details/
func (s *ProjectService) Update(organizationSlug string, slug string, params *UpdateProjectParams) (*Project, *http.Response, error) {
	project := new(Project)
	apiError := new(APIError)
	resp, err := s.sling.New().Put("projects/"+organizationSlug+"/"+slug+"/").BodyJSON(params).Receive(project, apiError)
	return project, resp, relevantError(err, *apiError)
}

// Delete a project.
// https://docs.sentry.io/api/projects/delete-project-details/
func (s *ProjectService) Delete(organizationSlug string, slug string) (*http.Response, error) {
	apiError := new(APIError)
	resp, err := s.sling.New().Delete("projects/"+organizationSlug+"/"+slug+"/").Receive(nil, apiError)
	return resp, relevantError(err, *apiError)
}
