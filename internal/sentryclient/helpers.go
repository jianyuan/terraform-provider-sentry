package sentryclient

import (
	"context"

	"github.com/jianyuan/go-sentry/v2/sentry"
)

func GetProjectIdToSlugMap(ctx context.Context, client *sentry.Client) (map[string]string, error) {
	projectMap := make(map[string]string)

	listParams := &sentry.ListCursorParams{}

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
