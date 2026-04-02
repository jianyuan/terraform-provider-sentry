package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccCronMonitorDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-cron-monitor")
	rn := "data.sentry_cron_monitor.test"

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
				Config: testAccCronMonitorDataSourceConfig(teamName, projectName, monitorName, `
					description = "cron monitor description"
					checkin_margin_minutes = 1
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin_minutes"), knownvalue.Int64Exact(1)),
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
				Config: testAccCronMonitorDataSourceConfig(teamName, projectName, monitorName+"-updated", `
					enabled = true
					checkin_margin_minutes = 10
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin_minutes"), knownvalue.Int64Exact(10)),
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
				Config: testAccCronMonitorDataSourceConfig(teamName, projectName, monitorName+"-updated", `
					enabled = false
					checkin_margin_minutes = 10
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin_minutes"), knownvalue.Int64Exact(10)),
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
		},
	})
}

func testAccCronMonitorDataSourceConfig(teamName, projectName, name, extras string) string {
	return testAccCronMonitorResourceConfig(teamName, projectName, name, extras) + `
		data "sentry_cron_monitor" "test" {
			organization = sentry_cron_monitor.test.organization
			id           = sentry_cron_monitor.test.id
		}
	`
}
