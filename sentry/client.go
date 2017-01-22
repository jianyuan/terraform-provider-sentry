package sentry

import (
	"fmt"
	"log"
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

func (m APIError) Error() error {
	if !m.HasError() {
		return nil
	}
	if detail, ok := m["detail"].(string); ok {
		// Some endpoints return an empty detail
		if detail == "" {
			return nil
		}
		return fmt.Errorf(detail)
	}
	return fmt.Errorf("API errored: %v", m)
}

func relevantError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func notFoundError(resp *http.Response, resource string) error {
	if resp != nil && resp.StatusCode == 404 {
		return fmt.Errorf("%s not found", resource)
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
	return &org, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization"))
}

func (c *Client) CreateOrganization(params *CreateOrganizationParams) (*Organization, *http.Response, error) {
	var org Organization
	apiErr := make(APIError)
	resp, err := c.sling.New().Post("0/organizations/").BodyJSON(params).Receive(&org, &apiErr)
	return &org, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization"))
}

func (c *Client) UpdateOrganization(slug string, params *UpdateOrganizationParams) (*Organization, *http.Response, error) {
	var org Organization
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/", slug)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&org, &apiErr)
	return &org, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization"))
}

func (c *Client) DeleteOrganization(slug string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/", slug)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization"))
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
	return &team, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or team"))
}

func (c *Client) CreateTeam(organizationSlug string, params *CreateTeamParams) (*Team, *http.Response, error) {
	var team Team
	apiErr := make(APIError)
	path := fmt.Sprintf("0/organizations/%s/teams/", organizationSlug)
	resp, err := c.sling.New().Post(path).BodyJSON(params).Receive(&team, &apiErr)
	return &team, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization"))
}

func (c *Client) UpdateTeam(organizationSlug, slug string, params *UpdateTeamParams) (*Team, *http.Response, error) {
	var team Team
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&team, &apiErr)
	return &team, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or team"))
}

func (c *Client) DeleteTeam(organizationSlug, slug string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or team"))
}

type ProjectOptions struct {
	ResolveAge int `json:"sentry:resolve_age"`
}

type Project struct {
	Organization Organization   `json:"organization"`
	Team         Team           `json:"team"`
	ID           string         `json:"id"`
	Slug         string         `json:"slug"`
	Name         string         `json:"name"`
	Options      ProjectOptions `json:"options"`
}

type CreateProjectParams struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

type UpdateProjectOptionsParams struct {
	ResolveAge int `json:"sentry:resolve_age,omitempty"`
}

type UpdateProjectParams struct {
	Name    string                     `json:"name"`
	Slug    string                     `json:"slug,omitempty"`
	Options UpdateProjectOptionsParams `json:"options,omitempty"`
}

func (c *Client) GetProject(organizationSlug, slug string) (*Project, *http.Response, error) {
	var proj Project
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Get(path).Receive(&proj, &apiErr)
	return &proj, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or project"))
}

func (c *Client) CreateProject(organizationSlug, teamSlug string, params *CreateProjectParams) (*Project, *http.Response, error) {
	var proj Project
	apiErr := make(APIError)
	path := fmt.Sprintf("0/teams/%s/%s/projects/", organizationSlug, teamSlug)
	resp, err := c.sling.New().Post(path).BodyJSON(params).Receive(&proj, &apiErr)
	return &proj, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or team"))
}

func (c *Client) UpdateProject(organizationSlug, slug string, params *UpdateProjectParams) (*Project, *http.Response, error) {
	var proj Project
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&proj, &apiErr)
	return &proj, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or project"))
}

func (c *Client) DeleteProject(organizationSlug, slug string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/", organizationSlug, slug)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or project"))
}

type DSN struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
	CSP    string `json:"csp"`
}

type Key struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Secret string `json:"secret"`
	DSN    DSN    `json:"dsn"`
}

type CreateKeyParams struct {
	Name string `json:"name"`
}

type UpdateKeyParams struct {
	Name string `json:"name"`
}

func (c *Client) GetKey(organizationSlug, projectSlug, keyID string) (*Key, *http.Response, error) {
	var key Key

	keys := new([]Key)
	apiErr := make(APIError)

	path := fmt.Sprintf("0/projects/%s/%s/keys/", organizationSlug, projectSlug)
	log.Printf("[DEBUG] Client.GetKey %s", path)

	resp, err := c.sling.New().Get(path).Receive(keys, &apiErr)

	for _, v := range *keys {
		if v.ID == keyID {
			key = v
		}
	}

	log.Printf("[DEBUG] Client.GetKey response %s", resp)

	return &key, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or project"))
}

func (c *Client) CreateKey(organizationSlug, projectSlug string, params *CreateKeyParams) (*Key, *http.Response, error) {
	var key Key
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/keys/", organizationSlug, projectSlug)
	resp, err := c.sling.New().Post(path).BodyJSON(params).Receive(&key, &apiErr)
	return &key, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization or project"))
}

func (c *Client) UpdateKey(organizationSlug, projectSlug string, keyID string, params *UpdateKeyParams) (*Key, *http.Response, error) {
	var key Key
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/keys/%s/", organizationSlug, projectSlug, keyID)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&key, &apiErr)
	return &key, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization, project or key"))
}

func (c *Client) DeleteKey(organizationSlug, projectSlug, keyID string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/keys/%s/", organizationSlug, projectSlug, keyID)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization, project or key"))
}

type PluginConfigEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Plugin struct {
	ID      string              `json:"id"`
	Enabled bool                `json:"enabled"`
	Config  []PluginConfigEntry `json:"config,omitempty"`
}

func (c *Client) GetPlugin(organizationSlug, projectSlug, id string) (*Plugin, *http.Response, error) {
	var plugin Plugin
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/plugins/%s/", organizationSlug, projectSlug, id)
	resp, err := c.sling.New().Get(path).Receive(&plugin, &apiErr)

	log.Printf("[DEBUG] Client.GetPlugin %s\n%v", path, resp)

	return &plugin, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization, project or plugin"))
}

func (c *Client) UpdatePlugin(organizationSlug, projectSlug, id string, params map[string]interface{}) (*Plugin, *http.Response, error) {
	var plugin Plugin
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/plugins/%s/", organizationSlug, projectSlug, id)
	resp, err := c.sling.New().Put(path).BodyJSON(params).Receive(&plugin, &apiErr)

	log.Printf("[DEBUG] Client.UpdatePlugin %s\n%v", path, resp)

	return &plugin, resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization, project or plugin"))
}

func (c *Client) EnablePlugin(organizationSlug, projectSlug, id string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/plugins/%s/", organizationSlug, projectSlug, id)
	resp, err := c.sling.New().Post(path).Receive(nil, &apiErr)

	log.Printf("[DEBUG] Client.EnablePlugin %s\n%v", path, resp)

	return resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization, project or plugin"))
}

func (c *Client) DisablePlugin(organizationSlug, projectSlug, id string) (*http.Response, error) {
	apiErr := make(APIError)
	path := fmt.Sprintf("0/projects/%s/%s/plugins/%s/", organizationSlug, projectSlug, id)
	resp, err := c.sling.New().Delete(path).Receive(nil, &apiErr)
	return resp, relevantError(err, apiErr.Error(), notFoundError(resp, "organization, project or plugin"))
}
