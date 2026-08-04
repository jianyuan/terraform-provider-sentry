package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

func init() {
	resource.AddTestSweepers("sentry_team", &resource.Sweeper{
		Name: "sentry_team",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListOrganizationTeamsParams{}
			for {
				listHttpResp, err := acctest.SharedApiClient.ListOrganizationTeamsWithResponse(ctx, acctest.TestOrganization, params)
				if err != nil {
					return err
				} else if listHttpResp.StatusCode() != http.StatusOK || listHttpResp.JSON200 == nil {
					return fmt.Errorf("failed to list organization teams: %s", listHttpResp.Status())
				}

				for _, team := range *listHttpResp.JSON200 {
					if !strings.HasPrefix(team.Slug, "tf-team") {
						continue
					}

					log.Printf("[INFO] Destroying Team: %s", team.Slug)

					_, err := acctest.SharedApiClient.DeleteOrganizationTeamWithResponse(ctx, acctest.TestOrganization, team.Id)
					if err != nil {
						log.Printf("[ERROR] Failed to destroy Team %q: %s", team.Slug, err)
						continue
					}

					log.Printf("[INFO] Team %q has been destroyed.", team.Slug)
				}

				params.Cursor = sentryclient.ParseNextPaginationCursor(listHttpResp.HTTPResponse)
				if params.Cursor == nil {
					break
				}
			}

			return nil
		},
	})
}

func testAccTeamResourceConfig(teamName string) string {
	return fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = "%[1]s"
	name         = "%[2]s"
	slug         = "%[2]s"
}
`, acctest.TestOrganization, teamName)
}
