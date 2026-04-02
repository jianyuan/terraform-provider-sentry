package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccUptimeMonitorDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-uptime-monitor")
	rn := "data.sentry_uptime_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
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
				Config: testAccUptimeMonitorDataSourceConfig(teamName, projectName, monitorName, `
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
				Config: testAccUptimeMonitorDataSourceConfig(teamName, projectName, monitorName+"-updated", `
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
		},
	})
}

func testAccUptimeMonitorDataSourceConfig(teamName, projectName, name, extras string) string {
	return testAccUptimeMonitorResourceConfig(teamName, projectName, name, extras) + `
		data "sentry_uptime_monitor" "test" {
			organization = sentry_uptime_monitor.test.organization
			id           = sentry_uptime_monitor.test.id
		}
	`
}
