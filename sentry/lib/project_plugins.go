package sentry

import (
	"context"
	"encoding/json"
	"fmt"
)

// ProjectPluginAsset represents an asset of a plugin.
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

// ProjectPluginsService provides methods for accessing Sentry project
// plugin API endpoints.
type ProjectPluginsService service

// List plugins bound to a project.
func (s *ProjectPluginsService) List(ctx context.Context, organizationSlug string, projectSlug string) ([]*ProjectPlugin, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/plugins/", organizationSlug, projectSlug)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectPlugins := []*ProjectPlugin{}
	resp, err := s.client.Do(ctx, req, &projectPlugins)
	if err != nil {
		return nil, resp, err
	}
	return projectPlugins, resp, nil
}

// Get details of a project plugin.
func (s *ProjectPluginsService) Get(ctx context.Context, organizationSlug string, projectSlug string, id string) (*ProjectPlugin, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/plugins/%v/", organizationSlug, projectSlug, id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectPlugin := new(ProjectPlugin)
	resp, err := s.client.Do(ctx, req, projectPlugin)
	if err != nil {
		return nil, resp, err
	}
	return projectPlugin, resp, nil
}

// UpdateProjectPluginParams are the parameters for TeamService.Update.
type UpdateProjectPluginParams map[string]interface{}

// Update settings for a given team.
// https://docs.sentry.io/api/teams/put-team-details/
func (s *ProjectPluginsService) Update(ctx context.Context, organizationSlug string, projectSlug string, id string, params UpdateProjectPluginParams) (*ProjectPlugin, *Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/plugins/%v/", organizationSlug, projectSlug, id)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	projectPlugin := new(ProjectPlugin)
	resp, err := s.client.Do(ctx, req, projectPlugin)
	if err != nil {
		return nil, resp, err
	}
	return projectPlugin, resp, nil
}

// Enable a project plugin.
func (s *ProjectPluginsService) Enable(ctx context.Context, organizationSlug string, projectSlug string, id string) (*Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/plugins/%v/", organizationSlug, projectSlug, id)
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Disable a project plugin.
func (s *ProjectPluginsService) Disable(ctx context.Context, organizationSlug string, projectSlug string, id string) (*Response, error) {
	u := fmt.Sprintf("0/projects/%v/%v/plugins/%v/", organizationSlug, projectSlug, id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
