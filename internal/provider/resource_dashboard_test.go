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
	resource.AddTestSweepers("sentry_dashboard", &resource.Sweeper{
		Name: "sentry_dashboard",
		F: func(r string) error {
			ctx := context.Background()

			listParams := &sentry.ListCursorParams{}

			for {
				dashboards, resp, err := acctest.SharedClient.Dashboards.List(ctx, acctest.TestOrganization, listParams)
				if err != nil {
					return err
				}

				for _, dashboard := range dashboards {
					if !strings.HasPrefix(sentry.StringValue(dashboard.Title), "tf-dashboard") {
						continue
					}

					log.Printf("[INFO] Destroying Dashboard %q", sentry.StringValue(dashboard.Title))

					_, err := acctest.SharedClient.Dashboards.Delete(ctx, acctest.TestOrganization, sentry.StringValue(dashboard.ID))
					if err != nil {
						log.Printf("[ERROR] Failed to destroy Dashboard %q: %s", sentry.StringValue(dashboard.Title), err)
						continue
					}

					log.Printf("[INFO] Dashboard %q has been destroyed.", sentry.StringValue(dashboard.Title))
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
