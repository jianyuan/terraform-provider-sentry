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

func TestAccMetricMonitorResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "sentry_metric_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`The argument "aggregatex" is required, but no definition was found.`),
			},
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
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"user_id": knownvalue.Null(),
			"team_id": knownvalue.NotNull(),
		})),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricMonitorResourceConfig(teamName, projectName, monitorName, `
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
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("aggregate"), knownvalue.StringExact("count()")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("dataset"), knownvalue.StringExact("events")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("event_types"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("default"),
						knownvalue.StringExact("error"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition_group"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"logic_type": knownvalue.StringExact("any"),
						"conditions": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"type":             knownvalue.StringExact("gt"),
								"comparison":       knownvalue.Int64Exact(100),
								"condition_result": knownvalue.Int64Exact(75),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"type":             knownvalue.StringExact("lte"),
								"comparison":       knownvalue.Int64Exact(50),
								"condition_result": knownvalue.Int64Exact(0),
							}),
						}),
					})),
				),
			},
			{
				Config: testAccMetricMonitorResourceConfig(teamName, projectName, monitorName+"-updated", `
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
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
				),
			},
			{
				Config: testAccMetricMonitorResourceConfig(teamName, projectName, monitorName+"-updated", `
					enabled = false

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
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetricMonitorResourceConfig(teamName, projectName, name, extras string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		resource "sentry_metric_monitor" "test" {
			organization = data.sentry_organization.test.slug
			project      = sentry_project.test.slug
			name         = "%[1]s"

			%[2]s

			owner = {
				team_id = sentry_team.test.internal_id
			}
		}
	`, name, extras)
}
