package sentryclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

// GetProjectIdToSlugMap returns a map of project ID to project slug for all
// projects in the organization.
//
// It uses the organization-scoped project list endpoint rather than the
// global /0/projects/ endpoint, which only returns the projects the token's
// user is a member of and rejects organization-scoped tokens such as internal
// integration tokens.
func GetProjectIdToSlugMap(ctx context.Context, apiClient *apiclient.ClientWithResponses, organization string) (map[string]string, error) {
	projectMap := make(map[string]string)

	params := &apiclient.ListOrganizationProjectsParams{}

	for {
		httpResp, err := apiClient.ListOrganizationProjectsWithResponse(ctx, organization, params)
		if err != nil {
			return nil, err
		} else if httpResp.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
		} else if httpResp.JSON200 == nil {
			return nil, errors.New("got empty response body")
		}

		for _, project := range *httpResp.JSON200 {
			projectMap[project.Id] = project.Slug
		}

		params.Cursor = ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	return projectMap, nil
}

// GetProjectSlug returns the slug of the project identified by projectIdOrSlug.
func GetProjectSlug(ctx context.Context, apiClient *apiclient.ClientWithResponses, organization string, projectIdOrSlug string) (string, error) {
	httpResp, err := apiClient.GetOrganizationProjectWithResponse(ctx, organization, projectIdOrSlug)
	if err != nil {
		return "", err
	} else if httpResp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
	} else if httpResp.JSON200 == nil {
		return "", errors.New("got empty response body")
	}
	return httpResp.JSON200.Slug, nil
}
