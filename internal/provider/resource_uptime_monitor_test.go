package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccUptimeMonitorResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-uptime-monitor")
	rn := "sentry_uptime_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_assignee"), knownvalue.ObjectExact(map[string]knownvalue.Check{
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
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("url"), knownvalue.StringExact("https://sentry.io")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("method"), knownvalue.StringExact("GET")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("body"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("headers"), knownvalue.MapSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval_seconds"), knownvalue.Int64Exact(60)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("timeout_ms"), knownvalue.Int64Exact(5000)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.StringExact("production")),
				),
			},
			{
				Config: testAccUptimeMonitorResourceConfig(teamName, projectName, monitorName+"-updated", `
					url = "https://us.sentry.io"
					method = "POST"
					body = "{}"
					headers = {
						"X-Header-Key" = "X-Header-Value"
					}
					interval_seconds = 300
					timeout_ms = 10000
					
					environment = "production"
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("url"), knownvalue.StringExact("https://us.sentry.io")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("method"), knownvalue.StringExact("POST")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("body"), knownvalue.StringExact("{}")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("headers"), knownvalue.MapExact(map[string]knownvalue.Check{
						"X-Header-Key": knownvalue.StringExact("X-Header-Value"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("interval_seconds"), knownvalue.Int64Exact(300)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("timeout_ms"), knownvalue.Int64Exact(10000)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.StringExact("production")),
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

			default_assignee = {
				team_id = sentry_team.test.internal_id
			}
		}
	`, name, extras)
}
