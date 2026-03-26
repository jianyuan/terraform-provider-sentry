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
	resource.AddTestSweepers("sentry_uptime_monitor", &resource.Sweeper{
		Name: "sentry_uptime_monitor",
		F: func(r string) error {
			ctx := context.Background()

			params := &apiclient.ListOrganizationMonitorsParams{
				Query: ptr.Ptr("!type:issue_stream type:uptime_domain_failure"),
			}

			for {
				listHttpResp, err := acctest.SharedApiClient.ListOrganizationMonitorsWithResponse(ctx, acctest.TestOrganization, params)
				if err != nil {
					return err
				} else if listHttpResp.StatusCode() != http.StatusOK || listHttpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Failed to list organization monitors: %s", listHttpResp.Status())
				}

				for _, monitor := range *listHttpResp.JSON200 {
					if !strings.HasPrefix(monitor.Name, "tf-uptime-monitor") {
						continue
					}

					deleteHttpResp, err := acctest.SharedApiClient.DeleteProjectMonitorWithResponse(ctx, acctest.TestOrganization, monitor.Id)
					if err != nil {
						log.Printf("[ERROR] Failed to delete uptime monitor: %s", err)
					} else if deleteHttpResp.StatusCode() != http.StatusNoContent {
						log.Printf("[ERROR] Failed to delete uptime monitor: %s", deleteHttpResp.Status())
					} else {
						log.Printf("[INFO] Deleted uptime monitor: %s (ID: %s)", monitor.Name, monitor.Id)
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

func TestAccUptimeMonitorResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "sentry_uptime_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "uptime monitor name"

						url = "https://sentry.io"
						method = "GET"
						body = "with body"
						interval_seconds = 60
						timeout_ms = 5000
						
						environment = "production"
					}
				`,
				ExpectError: acctest.ExpectLiteralError(
					`If method attribute is set and the value is one of "GET", "HEAD", "OPTIONS" this attribute is NULL`,
				),
			},
		},
	})
}

func TestAccUptimeMonitorResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-uptime-monitor")
	rn := "sentry_uptime_monitor.test"

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
				Config: testAccUptimeMonitorResourceConfig(teamName, projectName, monitorName, `
					url = "https://sentry.io"
					method = "GET"
					interval_seconds = 60
					timeout_ms = 5000
					
					environment = "production"

					enabled = true
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("url"), knownvalue.StringExact("https://sentry.io")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("method"), knownvalue.StringExact("GET")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("body"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("headers"), knownvalue.MapSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval_seconds"), knownvalue.Int64Exact(60)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("timeout_ms"), knownvalue.Int64Exact(5000)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.StringExact("production")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("assertion_json"), knownvalue.Null()),
				),
			},
			{
				Config: testAccUptimeMonitorResourceConfig(teamName, projectName, monitorName+"-updated", `
					url = "https://us.sentry.io"
					method = "POST"
					body = <<EOT
						{
							"key": "value"
						}
					EOT
					headers = {
						"X-Header-Key" = "X-Header-Value"
					}
					interval_seconds = 300
					timeout_ms = 10000
					
					environment = "production"

					assertion_json = provider::sentry::assertion(
						provider::sentry::op_and(
							provider::sentry::op_status_code_check("greater_than", 199),
							provider::sentry::op_status_code_check("less_than", 300),
						),
					)
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("url"), knownvalue.StringExact("https://us.sentry.io")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("method"), knownvalue.StringExact("POST")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("body"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("headers"), knownvalue.MapExact(map[string]knownvalue.Check{
						"X-Header-Key": knownvalue.StringExact("X-Header-Value"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval_seconds"), knownvalue.Int64Exact(300)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("timeout_ms"), knownvalue.Int64Exact(10000)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.StringExact("production")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("assertion_json"), knownvalue.NotNull()),
				),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"body", "assertion_json"},
			},
		},
	})
}

func testAccUptimeMonitorResourceConfig(teamName, projectName, name, extras string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		resource "sentry_uptime_monitor" "test" {
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
