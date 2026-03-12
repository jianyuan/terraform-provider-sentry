package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccCronMonitorResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4
					}
				`,
				ExpectError: regexp.MustCompile(`The argument "schedule" is required, but no definition was found.`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {}
					}
				`,
				ExpectError: regexp.MustCompile(
					`At least one attribute out of\n\[schedule.crontab.<.interval_value,schedule.crontab\] must be specified`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {}
					}
				`,
				ExpectError: regexp.MustCompile(
					`At least one attribute out of\n\[schedule.crontab.<.interval_unit,schedule.crontab\] must be specified`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
							interval_value = 1
							interval_unit = "day"
						}
					}
				`,
				ExpectError: regexp.MustCompile(
					`Attribute "schedule.interval_value" cannot be specified when\n"schedule.crontab" is specified`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
							interval_value = 1
							interval_unit = "day"
						}
					}
				`,
				ExpectError: regexp.MustCompile(
					`Attribute "schedule.interval_unit" cannot be specified when\n"schedule.crontab" is specified`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {
							interval_value = 1
						}
					}
				`,
				ExpectError: regexp.MustCompile(
					`Attribute "schedule.interval_unit" must be specified when\n"schedule.interval_value" is specified`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
						}

						default_assignee = {
							user_id = "3"
							team_id = "4"
						}
					}
				`,
				ExpectError: regexp.MustCompile(
					`Attribute "default_assignee.team_id" cannot be specified when\n"default_assignee.user_id" is specified`),
			},
			{
				Config: `
					resource "sentry_cron_monitor" "test" {
						organization = "1"
						project      = "2"
						name         = "cron monitor name"

						checkin_margin = 1
						failure_issue_threshold = 2
						max_runtime = 3
						recovery_threshold = 4

						schedule = {
							crontab = "0 0 * * *"
						}

						default_assignee = {
							user_id = "3"
							team_id = "4"
						}
					}
				`,
				ExpectError: regexp.MustCompile(
					`Attribute "default_assignee.user_id" cannot be specified when\n"default_assignee.team_id" is specified`),
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
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCronMonitorResourceConfig(teamName, projectName, monitorName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
				),
			},
		},
	})
}

func testAccCronMonitorResourceConfig(teamName, projectName, name string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		resource "sentry_cron_monitor" "test" {
			organization = data.sentry_organization.test.slug
			project      = sentry_project.test.slug
			name         = "%[1]s"

			checkin_margin = 1
			failure_issue_threshold = 2
			max_runtime = 3
			recovery_threshold = 4

			schedule = {
				// crontab = "0 0 * * *"
				interval_value = 1
				interval_unit = "day"
			}

			default_assignee = {
				team_id = sentry_team.test.internal_id
			}
		}
	`, name)
}
