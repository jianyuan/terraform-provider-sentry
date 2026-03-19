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

func TestAccAlertResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-monitor")
	alertName := acctest.RandomWithPrefix("tf-alert")
	opsgenieTeamName := acctest.RandomWithPrefix("tf-opsgenie")
	rn := "sentry_alert.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertResourceConfig(teamName, projectName, monitorName, alertName, opsgenieTeamName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alertName)),
				),
			},
			{
				Config: testAccAlertResourceConfig(teamName, projectName, monitorName, alertName+"-updated", opsgenieTeamName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alertName+"-updated")),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.TwoPartImportStateIdFunc(rn, "organization"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlertResourceConfig(teamName, projectName, monitorName, name, opsgenieTeamName string) string {
	return testAccMetricMonitorResourceConfig(teamName, projectName, monitorName) + fmt.Sprintf(`
		resource "sentry_alert" "test" {
			organization = data.sentry_organization.test.slug
			name         = "%[1]s"

			frequency_minutes = 1440
			environment       = "production"
			monitor_ids       = [sentry_metric_monitor.test.id]

			trigger_conditions = [
				"first_seen_event",
				"issue_resolved_trigger",
				"reappeared_event",
				"regression_event",
			]

			action_filters = [
				{
					logic_type = "all"
					conditions = [
						{
							age_comparison = {
								value = 1
								time = "minute"
								comparison_type = "older"
							}
						},
						{
							assigned_to = {
								target_type = "Team"
								target_id = sentry_team.test.internal_id
							}
						},
						{
							issue_category = {
								value = 1
							}
						},
						{
							issue_occurrences = {
								value = 1
							}
						},
						{
							issue_priority_deescalating = {}
						},
						{
							issue_priority_greater_or_equal = {
								comparison = 75
							}
						},
						{
							event_unique_user_frequency_count = {
								value = 1
								interval = "5m"
							}
						},
						{
							event_unique_user_frequency_count = {
								value = 1
								filters = [
									{ "key": "foo", "match": "eq", "value": "bar" }
								]
								interval = "1m"
							}
						},
						{
							event_frequency_count = {
								value = 100
								interval = "1m"
							}
						},
						{
							event_frequency_percent = {
								value = 100
								interval = "1h"
								comparison_interval = "1w"
							}
						},
						{
							percent_sessions_count = {
								value = 10
								interval = "1h"
							}
						},
						{
							percent_sessions_percent = {
								value = 10
								interval = "1h"
								comparison_interval = "1w"
							}
						},
						{
							event_attribute = {
								attribute = "message"
								match = "co"
								value = "bar"
							}
						},
						{
							tagged_event = {
								key = "level"
								match = "eq"
								value = "error"
							}
						},
						{
							tagged_event = {
								key = "level"
								match = "is"
							}
						},
						{
							latest_release = {}
						},
						{
							latest_adopted_release = {
								environment = "test"
								age_comparison = "older"
								release_age_type = "oldest"
							}
						},
						{
							level = {
								match = "eq"
								level = 50
							}
						}
					]
					actions = [
						{
							email = {
								target_type = "issue_owners"
								fallthrough_type = "AllMembers"
							}
						},
						{
							email = {
								target_type = "team"
								target_id = sentry_team.test.internal_id
							}
						},
						{
							plugin = {}
						},
						{
							slack = {
								integration_id = data.sentry_organization_integration.slack.id
								channel_name   = "#general"
								tags           = "one, two,three"
								notes          = "Please <http://example.com|click here> for triage information"
							}
						},
						// FIXME:
						// {
						// 	slack = {
						// 		integration_id = data.sentry_organization_integration.slack.id
						// 		channel_name   = "general"
						// 		notes          = "Please <http://example.com|click here> for triage information"
						// 	}
						// },
						{
							pagerduty = {
								integration_id = sentry_integration_pagerduty.pagerduty.integration_id
								service_name   = sentry_integration_pagerduty.pagerduty.service
								service_id     = sentry_integration_pagerduty.pagerduty.id
								severity       = "default"
							}
						},
						{
							discord = {
								integration_id = data.sentry_organization_integration.discord.id
								channel_id     = "714123428994482189"
							}
						},
						{
							discord = {
								integration_id = data.sentry_organization_integration.discord.id
								channel_id     = "714123428994482189"
								tags           = "environment, level"
							}
						},
						{
							opsgenie = {
								integration_id = sentry_integration_opsgenie.opsgenie.integration_id
								team_id        = sentry_integration_opsgenie.opsgenie.id
								team_name      = sentry_integration_opsgenie.opsgenie.team
								priority       = "P1"
							}
						},
						{
							github = {
								integration_id = data.sentry_organization_integration.github.id
								repo           = "terraform-provider-sentry"
								assignee       = "jianyuan"
								labels         = ["bug"]
							}
						},
					]
				}
			]
		}
	`, name) + fmt.Sprintf(`
		# Slack
		data "sentry_organization_integration" "slack" {
			organization = data.sentry_organization.test.slug
			provider_key = "slack"
			name         = "A2 Marketing"  # TODO: Use a real integration name
		}

		# PagerDuty
		data "sentry_organization_integration" "pagerduty" {
			organization = sentry_project.test.organization
			provider_key = "pagerduty"
			name         = "terraform-provider-sentry"
		}

		resource "sentry_integration_pagerduty" "pagerduty" {
			organization    = data.sentry_organization_integration.pagerduty.organization
			integration_id  = data.sentry_organization_integration.pagerduty.id
			service         = "issue-alert-service"
			integration_key = "issue-alert-integration-key"
		}

		# Discord
		data "sentry_organization_integration" "discord" {
			organization = sentry_project.test.organization
			provider_key = "discord"
			name         = "jy's server"
		}

		# Opsgenie
		data "sentry_organization_integration" "opsgenie" {
			organization = sentry_project.test.organization
			provider_key = "opsgenie"
			name         = "terraform-provider-sentry"
		}

		resource "sentry_integration_opsgenie" "opsgenie" {
			organization    = data.sentry_organization_integration.opsgenie.organization
			integration_id  = data.sentry_organization_integration.opsgenie.id
			integration_key = "%[1]s"
			team            = "%[2]s"
		}

		# GitHub
		data "sentry_organization_integration" "github" {
			organization = sentry_project.test.organization
			provider_key = "github"
			name         = "jianyuan"
		}
	`, acctest.TestOpsgenieIntegrationKey, opsgenieTeamName)
}
