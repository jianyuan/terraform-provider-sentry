package main

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type Client struct {
	sling *sling.Sling
	token string
}

func NewClient(httpClient *http.Client, token, baseURL string) *Client {
	baseSling := sling.New().Client(httpClient).Base(baseURL)

	if token != "" {
		baseSling = baseSling.Add("Authorization", "Bearer "+token)
	}

	return &Client{
		sling: baseSling,
		token: token,
	}
}

type APIError map[string]interface{}

func (m APIError) HasError() bool {
	return len(m) > 0
}

func (m APIError) Error() string {
	if !m.HasError() {
		return ""
	}
	if detail, ok := m["detail"].(string); ok {
		return detail
	}
	// TODO
	return "field errors"
}

func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError.HasError() {
		return apiError
	}
	return nil
}

type Organization struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CreateOrganizationParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type UpdateOrganizationParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

func (c *Client) GetOrganization(slug string) (*Organization, *http.Response, error) {
	var org Organization
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/", slug)
	resp, err := c.sling.New().Get(path).Receive(&org, &apiErr)
	return &org, resp, relevantError(err, apiErr)
}

func (c *Client) CreateOrganization(params *CreateOrganizationParams) (*Organization, *http.Response, error) {
	var org Organization
	apiErr := make(APIError)
	resp, err := c.sling.New().Post("0/organizations/").BodyJSON(params).Receive(&org, &apiErr)
	return &org, resp, relevantError(err, apiErr)
}

func (c *Client) UpdateOrganization(slug string, params *UpdateOrganizationParams) (*Organization, *http.Response, error) {
	var org Organization
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/", slug)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&org, &apiErr)
	return &org, resp, relevantError(err, apiErr)
}

func (c *Client) DeleteOrganization(slug string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/", slug)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr)
}

type Team struct {
	Organization Organization `json:"organization"`
	ID           string       `json:"id"`
	Slug         string       `json:"slug"`
	Name         string       `json:"name"`
}

type CreateTeamParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type UpdateTeamParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

func (c *Client) GetTeam(organizationSlug, slug string) (*Team, *http.Response, error) {
	var team Team
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Get(path).Receive(&team, &apiErr)
	return &team, resp, relevantError(err, apiErr)
}

func (c *Client) CreateTeam(organizationSlug string, params *CreateTeamParams) (*Team, *http.Response, error) {
	var team Team
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/teams/", organizationSlug)
	resp, err := c.sling.New().Post(path).BodyJSON(params).Receive(&team, &apiErr)
	return &team, resp, relevantError(err, apiErr)
}

func (c *Client) UpdateTeam(organizationSlug, slug string, params *UpdateTeamParams) (*Team, *http.Response, error) {
	var team Team
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&team, &apiErr)
	return &team, resp, relevantError(err, apiErr)
}

func (c *Client) DeleteTeam(organizationSlug, slug string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr)
}

type Project struct {
	Organization Organization `json:"organization"`
	Team         Team         `json:"team"`
	ID           string       `json:"id"`
	Slug         string       `json:"slug"`
	Name         string       `json:"name"`
}

type CreateProjectParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type UpdateProjectParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

func (c *Client) GetProject(organizationSlug, slug string) (*Project, *http.Response, error) {
	var proj Project
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Get(path).Receive(&proj, &apiErr)
	return &proj, resp, relevantError(err, apiErr)
}

func (c *Client) CreateProject(organizationSlug, teamSlug string, params *CreateProjectParams) (*Project, *http.Response, error) {
	var proj Project
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/projects/", organizationSlug, teamSlug)
	resp, err := c.sling.New().Post(path).BodyJSON(params).Receive(&proj, &apiErr)
	return &proj, resp, relevantError(err, apiErr)
}

func (c *Client) UpdateProject(organizationSlug, slug string, params *UpdateProjectParams) (*Project, *http.Response, error) {
	var proj Project
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&proj, &apiErr)
	return &proj, resp, relevantError(err, apiErr)
}

func (c *Client) DeleteProject(organizationSlug, slug string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr)
}
