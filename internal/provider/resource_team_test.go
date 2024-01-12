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
	resource.AddTestSweepers("sentry_team", &resource.Sweeper{
		Name: "sentry_team",
		F: func(r string) error {
			ctx := context.Background()

			listParams := &sentry.ListCursorParams{}

			for {
				teams, resp, err := acctest.SharedClient.Teams.List(ctx, acctest.TestOrganization, listParams)
				if err != nil {
					return err
				}

				for _, team := range teams {
					if !strings.HasPrefix(sentry.StringValue(team.Slug), "tf-team") {
						continue
					}

					log.Printf("[INFO] Destroying Team: %s", sentry.StringValue(team.Slug))

					_, err := acctest.SharedClient.Teams.Delete(ctx, acctest.TestOrganization, sentry.StringValue(team.Slug))
					if err != nil {
						log.Printf("[ERROR] Failed to destroy Team %q: %s", sentry.StringValue(team.Slug), err)
						continue
					}

					log.Printf("[INFO] Team %q has been destroyed.", sentry.StringValue(team.Slug))
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
