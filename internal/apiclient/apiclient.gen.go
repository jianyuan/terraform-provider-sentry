// Package apiclient provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Organization defines model for Organization.
type Organization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Project defines model for Project.
type Project struct {
	Color           string                 `json:"color"`
	DateCreated     time.Time              `json:"dateCreated"`
	DigestsMaxDelay int64                  `json:"digestsMaxDelay"`
	DigestsMinDelay int64                  `json:"digestsMinDelay"`
	Features        []string               `json:"features"`
	Id              string                 `json:"id"`
	IsPublic        bool                   `json:"isPublic"`
	Name            string                 `json:"name"`
	Options         map[string]interface{} `json:"options"`
	Organization    Organization           `json:"organization"`
	Platform        *string                `json:"platform"`
	ResolveAge      int64                  `json:"resolveAge"`
	Slug            string                 `json:"slug"`
	Teams           []Team                 `json:"teams"`
}

// Team defines model for Team.
type Team struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// OrganizationIdOrSlug defines model for organization_id_or_slug.
type OrganizationIdOrSlug = string

// ProjectIdOrSlug defines model for project_id_or_slug.
type ProjectIdOrSlug = string

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetOrganization request
	GetOrganization(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetOrganizationProject request
	GetOrganizationProject(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, projectIdOrSlug ProjectIdOrSlug, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetOrganization(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetOrganizationRequest(c.Server, organizationIdOrSlug)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetOrganizationProject(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, projectIdOrSlug ProjectIdOrSlug, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetOrganizationProjectRequest(c.Server, organizationIdOrSlug, projectIdOrSlug)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetOrganizationRequest generates requests for GetOrganization
func NewGetOrganizationRequest(server string, organizationIdOrSlug OrganizationIdOrSlug) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "organization_id_or_slug", runtime.ParamLocationPath, organizationIdOrSlug)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/0/organizations/%s/", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetOrganizationProjectRequest generates requests for GetOrganizationProject
func NewGetOrganizationProjectRequest(server string, organizationIdOrSlug OrganizationIdOrSlug, projectIdOrSlug ProjectIdOrSlug) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "organization_id_or_slug", runtime.ParamLocationPath, organizationIdOrSlug)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "project_id_or_slug", runtime.ParamLocationPath, projectIdOrSlug)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/0/projects/%s/%s/", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetOrganizationWithResponse request
	GetOrganizationWithResponse(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, reqEditors ...RequestEditorFn) (*GetOrganizationResponse, error)

	// GetOrganizationProjectWithResponse request
	GetOrganizationProjectWithResponse(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, projectIdOrSlug ProjectIdOrSlug, reqEditors ...RequestEditorFn) (*GetOrganizationProjectResponse, error)
}

type GetOrganizationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Organization
}

// Status returns HTTPResponse.Status
func (r GetOrganizationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetOrganizationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetOrganizationProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Project
}

// Status returns HTTPResponse.Status
func (r GetOrganizationProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetOrganizationProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetOrganizationWithResponse request returning *GetOrganizationResponse
func (c *ClientWithResponses) GetOrganizationWithResponse(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, reqEditors ...RequestEditorFn) (*GetOrganizationResponse, error) {
	rsp, err := c.GetOrganization(ctx, organizationIdOrSlug, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetOrganizationResponse(rsp)
}

// GetOrganizationProjectWithResponse request returning *GetOrganizationProjectResponse
func (c *ClientWithResponses) GetOrganizationProjectWithResponse(ctx context.Context, organizationIdOrSlug OrganizationIdOrSlug, projectIdOrSlug ProjectIdOrSlug, reqEditors ...RequestEditorFn) (*GetOrganizationProjectResponse, error) {
	rsp, err := c.GetOrganizationProject(ctx, organizationIdOrSlug, projectIdOrSlug, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetOrganizationProjectResponse(rsp)
}

// ParseGetOrganizationResponse parses an HTTP response from a GetOrganizationWithResponse call
func ParseGetOrganizationResponse(rsp *http.Response) (*GetOrganizationResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetOrganizationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Organization
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetOrganizationProjectResponse parses an HTTP response from a GetOrganizationProjectWithResponse call
func ParseGetOrganizationProjectResponse(rsp *http.Response) (*GetOrganizationProjectResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetOrganizationProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Project
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
