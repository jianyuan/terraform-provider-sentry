package sentry

import (
	"net/http"
	"net/url"

	"path"

	"github.com/dghubble/sling"
)

const (
	DefaultBaseURL = "https://sentry.io/api/"
	APIVersion     = "0"
)

type Client struct {
	sling          *sling.Sling
	Organizations  *OrganizationService
	Teams          *TeamService
	Projects       *ProjectService
	ProjectKeys    *ProjectKeyService
	ProjectPlugins *ProjectPluginService
}

// NewClient returns a new Sentry API client.
// If a nil httpClient is given, the http.DefaultClient will be used.
// If a nil baseURL is given, the DefaultBaseURL will be used.
func NewClient(httpClient *http.Client, baseURL *url.URL, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if baseURL == nil {
		baseURL, _ = url.Parse(DefaultBaseURL)
	}
	baseURL.Path = path.Join(baseURL.Path, APIVersion) + "/"

	base := sling.New().Base(baseURL.String()).Client(httpClient)

	if token != "" {
		base.Add("Authorization", "Bearer "+token)
	}

	c := &Client{
		sling:          base,
		Organizations:  newOrganizationService(base.New()),
		Teams:          newTeamService(base.New()),
		Projects:       newProjectService(base.New()),
		ProjectKeys:    newProjectKeyService(base.New()),
		ProjectPlugins: newProjectPluginService(base.New()),
	}
	return c
}
