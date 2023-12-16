package provider

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func init() {
	resource.AddTestSweepers("sentry_project", &resource.Sweeper{
		Name: "sentry_project",
		F: func(r string) error {
			ctx := context.Background()

			listParams := &sentry.ListCursorParams{}

			for {
				projects, resp, err := acctest.SharedClient.Projects.List(ctx, listParams)
				if err != nil {
					return err
				}

				for _, project := range projects {
					if !strings.HasPrefix(project.Slug, "tf-project") {
						continue
					}

					log.Printf("[INFO] Destroying Project: %s", project.Slug)

					_, err := acctest.SharedClient.Projects.Delete(ctx, acctest.TestOrganization, project.Slug)
					if err != nil {
						log.Printf("[ERROR] Failed to destroy Project %q: %s", project.Slug, err)
						continue
					}

					log.Printf("[INFO] Project %q has been destroyed.", project.Slug)
				}

				if resp.Cursor == "" {
					break
				}
				listParams.Cursor = resp.Cursor
			}

			return nil
		},
	})
}
