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
	return fmt.Sprintf("%v", m)
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
	Name string `url:"name"`
	Slug string `url:"slug,omitempty"`
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
	resp, err := c.sling.New().Post("0/organizations/").BodyForm(params).Receive(&org, &apiErr)
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
