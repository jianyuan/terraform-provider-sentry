package sentryclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func GetProjectIdToSlugMap(ctx context.Context, client *sentry.Client) (map[string]string, error) {
	projectMap := make(map[string]string)

	listParams := &sentry.ListProjectsParams{}

	for {
		projects, resp, err := client.Projects.List(ctx, listParams)
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			projectMap[project.ID] = project.Slug
		}

		if resp.Cursor == "" {
			break
		}
		listParams.Cursor = resp.Cursor
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
