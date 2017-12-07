package sentry

import (
	"encoding/json"
	"net/http"

	"github.com/dghubble/sling"
)

// ProjectPlugin represents an asset of a plugin.
type ProjectPluginAsset struct {
	URL string `json:"url"`
}

// ProjectPluginConfig represents the configuration of a plugin.
// Based on https://github.com/getsentry/sentry/blob/96bc1c63df5ec73fe12c136ada11561bf52f1ec9/src/sentry/api/serializers/models/plugin.py#L62-L94.
type ProjectPluginConfig struct {
	Name         string          `json:"name"`
	Label        string          `json:"label"`
	Type         string          `json:"type"`
	Required     bool            `json:"required"`
	Help         string          `json:"help"`
	Placeholder  string          `json:"placeholder"`
	Choices      json.RawMessage `json:"choices"`
	ReadOnly     bool            `json:"readonly"`
	DefaultValue interface{}     `json:"defaultValue"`
	Value        interface{}     `json:"value"`
}

// ProjectPlugin represents a plugin bound to a project.
// Based on https://github.com/getsentry/sentry/blob/96bc1c63df5ec73fe12c136ada11561bf52f1ec9/src/sentry/api/serializers/models/plugin.py#L11.
type ProjectPlugin struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	CanDisable bool                   `json:"canDisable"`
	IsTestable bool                   `json:"isTestable"`
	Metadata   map[string]interface{} `json:"metadata"`
	Contexts   []string               `json:"contexts"`
	Status     string                 `json:"status"`
	Assets     []ProjectPluginAsset   `json:"assets"`
	Doc        string                 `json:"doc"`
	Config     []ProjectPluginConfig  `json:"config"`
}

// ProjectPluginService provides methods for accessing Sentry project
// plugin API endpoints.
type ProjectPluginService struct {
	sling *sling.Sling
}

func newProjectPluginService(sling *sling.Sling) *ProjectPluginService {
	return &ProjectPluginService{
		sling: sling,
	}
}

// List plugins bound to a project.
func (s *ProjectPluginService) List(organizationSlug string, projectSlug string) ([]ProjectPlugin, *http.Response, error) {
	projectPlugins := new([]ProjectPlugin)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("projects/"+organizationSlug+"/"+projectSlug+"/plugins/").Receive(projectPlugins, apiError)
	return *projectPlugins, resp, relevantError(err, *apiError)
}

// Get details of a project plugin.
func (s *ProjectPluginService) Get(organizationSlug string, projectSlug string, id string) (*ProjectPlugin, *http.Response, error) {
	projectPlugin := new(ProjectPlugin)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("projects/"+organizationSlug+"/"+projectSlug+"/plugins/"+id+"/").Receive(projectPlugin, apiError)
	return projectPlugin, resp, relevantError(err, *apiError)
}

// UpdateTeamParams are the parameters for TeamService.Update.
type UpdateProjectPluginParams map[string]interface{}

// Update settings for a given team.
// https://docs.sentry.io/api/teams/put-team-details/
func (s *ProjectPluginService) Update(organizationSlug string, projectSlug string, id string, params UpdateProjectPluginParams) (*ProjectPlugin, *http.Response, error) {
	projectPlugin := new(ProjectPlugin)
	apiError := new(APIError)
	resp, err := s.sling.New().Put("projects/"+organizationSlug+"/"+projectSlug+"/plugins/"+id+"/").BodyJSON(params).Receive(projectPlugin, apiError)
	return projectPlugin, resp, relevantError(err, *apiError)
}

// Enable a project plugin.
func (s *ProjectPluginService) Enable(organizationSlug string, projectSlug string, id string) (*http.Response, error) {
	apiError := new(APIError)
	resp, err := s.sling.New().Post("projects/"+organizationSlug+"/"+projectSlug+"/plugins/"+id+"/").Receive(nil, apiError)
	return resp, relevantError(err, *apiError)
}

// Disable a project plugin.
func (s *ProjectPluginService) Disable(organizationSlug string, projectSlug string, id string) (*http.Response, error) {
	apiError := new(APIError)
	resp, err := s.sling.New().Delete("projects/"+organizationSlug+"/"+projectSlug+"/plugins/"+id+"/").Receive(nil, apiError)
	return resp, relevantError(err, *apiError)
}
