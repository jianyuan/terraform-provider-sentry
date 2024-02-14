package sentry

import (
	"context"
	"fmt"
	"time"
)

// Team represents a Sentry team that is bound to an organization.
// https://github.com/getsentry/sentry/blob/22.5.0/src/sentry/api/serializers/models/team.py#L109-L119
type Team struct {
	ID          *string    `json:"id,omitempty"`
	Slug        *string    `json:"slug,omitempty"`
	Name        *string    `json:"name,omitempty"`
	DateCreated *time.Time `json:"dateCreated,omitempty"`
	IsMember    *bool      `json:"isMember,omitempty"`
	TeamRole    *string    `json:"teamRole,omitempty"`
	HasAccess   *bool      `json:"hasAccess,omitempty"`
	IsPending   *bool      `json:"isPending,omitempty"`
	MemberCount *int       `json:"memberCount,omitempty"`
	Avatar      *Avatar    `json:"avatar,omitempty"`
	// TODO: externalTeams
	// TODO: projects
}

// TeamsService provides methods for accessing Sentry team API endpoints.
// https://docs.sentry.io/api/teams/
type TeamsService service

// List returns a list of teams bound to an organization.
// https://docs.sentry.io/api/teams/list-an-organizations-teams/
func (s *TeamsService) List(ctx context.Context, organizationSlug string) ([]*Team, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/teams/", organizationSlug)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	teams := []*Team{}
	resp, err := s.client.Do(ctx, req, &teams)
	if err != nil {
		return nil, resp, err
	}
	return teams, resp, nil
}

// Get details on an individual team of an organization.
// https://docs.sentry.io/api/teams/retrieve-a-team/
func (s *TeamsService) Get(ctx context.Context, organizationSlug string, slug string) (*Team, *Response, error) {
	u := fmt.Sprintf("0/teams/%v/%v/", organizationSlug, slug)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	team := new(Team)
	resp, err := s.client.Do(ctx, req, team)
	if err != nil {
		return nil, resp, err
	}
	return team, resp, nil
}

// CreateTeamParams are the parameters for TeamService.Create.
type CreateTeamParams struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// Create a new Sentry team bound to an organization.
// https://docs.sentry.io/api/teams/create-a-new-team/
func (s *TeamsService) Create(ctx context.Context, organizationSlug string, params *CreateTeamParams) (*Team, *Response, error) {
	u := fmt.Sprintf("0/organizations/%v/teams/", organizationSlug)
	req, err := s.client.NewRequest("POST", u, params)
	if err != nil {
		return nil, nil, err
	}

	team := new(Team)
	resp, err := s.client.Do(ctx, req, team)
	if err != nil {
		return nil, resp, err
	}
	return team, resp, nil
}

// UpdateTeamParams are the parameters for TeamService.Update.
type UpdateTeamParams struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// Update settings for a given team.
// https://docs.sentry.io/api/teams/update-a-team/
func (s *TeamsService) Update(ctx context.Context, organizationSlug string, slug string, params *UpdateTeamParams) (*Team, *Response, error) {
	u := fmt.Sprintf("0/teams/%v/%v/", organizationSlug, slug)
	req, err := s.client.NewRequest("PUT", u, params)
	if err != nil {
		return nil, nil, err
	}

	team := new(Team)
	resp, err := s.client.Do(ctx, req, team)
	if err != nil {
		return nil, resp, err
	}
	return team, resp, nil
}

// Delete a team.
// https://docs.sentry.io/api/teams/update-a-team/
func (s *TeamsService) Delete(ctx context.Context, organizationSlug string, slug string) (*Response, error) {
	u := fmt.Sprintf("0/teams/%v/%v/", organizationSlug, slug)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
