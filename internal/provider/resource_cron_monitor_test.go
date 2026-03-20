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
	resource.AddTestSweepers("sentry_cron_monitor", &resource.Sweeper{
		Name: "sentry_cron_monitor",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListOrganizationMonitorsParams{
				Query: ptr.Ptr("!type:issue_stream type:monitor_check_in_failure"),
			}

			for {
				listHttpResp, err := acctest.SharedApiClient.ListOrganizationMonitorsWithResponse(ctx, acctest.TestOrganization, params)
				if err != nil {
					return err
				} else if listHttpResp.StatusCode() != http.StatusOK || listHttpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Failed to list organization monitors: %s", listHttpResp.Status())
				}

				for _, monitor := range *listHttpResp.JSON200 {
					if !strings.HasPrefix(monitor.Name, "tf-cron-monitor") {
						continue
					}

					deleteHttpResp, err := acctest.SharedApiClient.DeleteProjectMonitorWithResponse(ctx, acctest.TestOrganization, monitor.Id)
					if err != nil {
						log.Printf("[ERROR] Failed to delete cron monitor: %s", err)
					} else if deleteHttpResp.StatusCode() != http.StatusNoContent {
						log.Printf("[ERROR] Failed to delete cron monitor: %s", deleteHttpResp.Status())
					} else {
						log.Printf("[INFO] Deleted cron monitor: %s (ID: %s)", monitor.Name, monitor.Id)
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

func TestAccCronMonitorResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`The argument "schedule" is required, but no definition was found.`),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					"At least one attribute out of [schedule.crontab.<.interval_value,schedule.crontab] must be specified",
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					"At least one attribute out of [schedule.crontab.<.interval_unit,schedule.crontab] must be specified",
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
							interval_value = 1
							interval_unit = "day"
						}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "schedule.interval_value" cannot be specified when "schedule.crontab" is specified`,
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
							interval_value = 1
							interval_unit = "day"
						}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "schedule.interval_unit" cannot be specified when "schedule.crontab" is specified`,
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {
							interval_value = 1
						}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "schedule.interval_unit" must be specified when "schedule.interval_value" is specified`,
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
						}

						owner = {
							user_id = "3"
							team_id = "4"
						}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "owner.team_id" cannot be specified when "owner.user_id" is specified`,
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
						}

						owner = {
							user_id = "3"
							team_id = "4"
						}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "owner.user_id" cannot be specified when "owner.team_id" is specified`,
				),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime_minutes = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
						}

						owner = {}
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					"No attribute specified when one (and only one) of [owner.user_id.<.team_id] is required",
				),
			},
		},
	})
}

func TestAccCronMonitorResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-cron-monitor")
	rn := "sentry_cron_monitor.test"

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
				Config: testAccCronMonitorResourceConfig(teamName, projectName, monitorName, `
					description = "cron monitor description"
					checkin_margin = 1
					failure_issue_threshold = 2
					max_runtime_minutes = 3
					recovery_threshold = 4
					schedule = {
						interval_value = 1
						interval_unit = "day"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.StringExact("cron monitor description")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("failure_issue_threshold"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_runtime_minutes"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("recovery_threshold"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("schedule"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"crontab":        knownvalue.Null(),
						"interval_value": knownvalue.Int64Exact(1),
						"interval_unit":  knownvalue.StringExact("day"),
					})),
				),
			},
			{
				Config: testAccCronMonitorResourceConfig(teamName, projectName, monitorName+"-updated", `
					enabled = true
					checkin_margin = 10
					failure_issue_threshold = 20
					max_runtime_minutes = 30
					recovery_threshold = 40
					schedule = {
						crontab = "0 0 * * *"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin"), knownvalue.Int64Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("failure_issue_threshold"), knownvalue.Int64Exact(20)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_runtime_minutes"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("recovery_threshold"), knownvalue.Int64Exact(40)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("schedule"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"crontab":        knownvalue.StringExact("0 0 * * *"),
						"interval_value": knownvalue.Null(),
						"interval_unit":  knownvalue.Null(),
					})),
				),
			},
			{
				Config: testAccCronMonitorResourceConfig(teamName, projectName, monitorName+"-updated", `
					enabled = false
					checkin_margin = 10
					failure_issue_threshold = 20
					max_runtime_minutes = 30
					recovery_threshold = 40
					schedule = {
						crontab = "0 0 * * *"
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin"), knownvalue.Int64Exact(10)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("failure_issue_threshold"), knownvalue.Int64Exact(20)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_runtime_minutes"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("recovery_threshold"), knownvalue.Int64Exact(40)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("schedule"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"crontab":        knownvalue.StringExact("0 0 * * *"),
						"interval_value": knownvalue.Null(),
						"interval_unit":  knownvalue.Null(),
					})),
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

func testAccCronMonitorResourceConfig(teamName, projectName, name, extras string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		resource "sentry_cron_monitor" "test" {
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
