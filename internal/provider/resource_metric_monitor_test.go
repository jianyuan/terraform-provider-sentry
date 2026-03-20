package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

func init() {
	resource.AddTestSweepers("sentry_metric_monitor", &resource.Sweeper{
		Name: "sentry_metric_monitor",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListOrganizationMonitorsParams{
				Query: ptr.Ptr("!type:issue_stream type:metric_issue"),
			}

			for {
				listHttpResp, err := acctest.SharedApiClient.ListOrganizationMonitorsWithResponse(ctx, acctest.TestOrganization, params)
				if err != nil {
					return err
				} else if listHttpResp.StatusCode() != http.StatusOK || listHttpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Failed to list organization monitors: %s", listHttpResp.Status())
				}

				for _, monitor := range *listHttpResp.JSON200 {
					if !strings.HasPrefix(monitor.Name, "tf-metric-monitor") {
						continue
					}

					deleteHttpResp, err := acctest.SharedApiClient.DeleteProjectMonitorWithResponse(ctx, acctest.TestOrganization, monitor.Id)
					if err != nil {
						log.Printf("[ERROR] Failed to delete metric monitor: %s", err)
					} else if deleteHttpResp.StatusCode() != http.StatusNoContent {
						log.Printf("[ERROR] Failed to delete metric monitor: %s", deleteHttpResp.Status())
					} else {
						log.Printf("[INFO] Deleted metric monitor: %s (ID: %s)", monitor.Name, monitor.Id)
					}
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

func TestAccMetricMonitorResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-metric-monitor")
	rn := "sentry_metric_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricMonitorResourceConfig(teamName, projectName, monitorName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
				),
			},
		},
	})
}

func testAccMetricMonitorResourceConfig(teamName, projectName, name string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		resource "sentry_metric_monitor" "test" {
			organization = data.sentry_organization.test.slug
			project      = sentry_project.test.slug
			name         = "%[1]s"

			aggregate = "count()"
			dataset = "events"
			event_types = ["default", "error"]

			condition_group = {
				conditions = [
					{
						type = "gt"
						comparison = 100
						condition_result = 75
					},
					{
						type = "lte"
						comparison = 50
						condition_result = 0
					},
				]
			}

			issue_detection = {
				type = "static"
			}

			default_assignee = {
				team_id = sentry_team.test.internal_id
			}
		}
	`, name)
}
