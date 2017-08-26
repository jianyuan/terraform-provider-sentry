package sentry

import (
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Team represents a Sentry team that is bound to an organization.
type Team struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"dateCreated"`
	HasAccess   bool      `json:"hasAccess"`
	IsPending   bool      `json:"isPending"`
	IsMember    bool      `json:"isMember"`
}

// TeamService provides methods for accessing Sentry team API endpoints.
// https://docs.sentry.io/api/teams/
type TeamService struct {
	sling *sling.Sling
}

func newTeamService(sling *sling.Sling) *TeamService {
	return &TeamService{
		sling: sling,
	}
}

// List returns a list of teams bound to an organization.
// https://docs.sentry.io/api/teams/get-organization-teams/
func (s *TeamService) List(organizationSlug string) ([]Team, *http.Response, error) {
	teams := new([]Team)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("organizations/"+organizationSlug+"/teams/").Receive(teams, apiError)
	return *teams, resp, relevantError(err, *apiError)
}

// Get details on an individual team of an organization.
// https://docs.sentry.io/api/teams/get-team-details/
func (s *TeamService) Get(organizationSlug string, slug string) (*Team, *http.Response, error) {
	team := new(Team)
	apiError := new(APIError)
	resp, err := s.sling.New().Get("teams/"+organizationSlug+"/"+slug+"/").Receive(team, apiError)
	return team, resp, relevantError(err, *apiError)
}

// CreateTeamParams are the parameters for TeamService.Create.
type CreateTeamParams struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// Create a new Sentry team bound to an organization.
// https://docs.sentry.io/api/teams/post-organization-teams/
func (s *TeamService) Create(organizationSlug string, params *CreateTeamParams) (*Team, *http.Response, error) {
	team := new(Team)
	apiError := new(APIError)
	resp, err := s.sling.New().Post("organizations/"+organizationSlug+"/teams/").BodyJSON(params).Receive(team, apiError)
	return team, resp, relevantError(err, *apiError)
}

// UpdateTeamParams are the parameters for TeamService.Update.
type UpdateTeamParams struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// Update settings for a given team.
// https://docs.sentry.io/api/teams/put-team-details/
func (s *TeamService) Update(organizationSlug string, slug string, params *UpdateTeamParams) (*Team, *http.Response, error) {
	team := new(Team)
	apiError := new(APIError)
	resp, err := s.sling.New().Put("teams/"+organizationSlug+"/"+slug+"/").BodyJSON(params).Receive(team, apiError)
	return team, resp, relevantError(err, *apiError)
}

// Delete a team.
// https://docs.sentry.io/api/teams/delete-team-details/
func (s *TeamService) Delete(organizationSlug string, slug string) (*http.Response, error) {
	apiError := new(APIError)
	resp, err := s.sling.New().Delete("teams/"+organizationSlug+"/"+slug+"/").Receive(nil, apiError)
	return resp, relevantError(err, *apiError)
}
