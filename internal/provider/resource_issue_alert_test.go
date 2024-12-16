package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccIssueAlertResource_validation(t *testing.T) {
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccIssueAlertConfig(team, project, alert, ``),
				ExpectError: acctest.ExpectLiteralError(`At least one of these attributes must be configured: [actions,actions_v2]`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions = "[]"
				`),
				ExpectError: acctest.ExpectLiteralError(`You must add an action for this alert to fire`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions_v2 = []
				`),
				ExpectError: acctest.ExpectLiteralError(`You must add an action for this alert to fire`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					conditions    = "[]"
					conditions_v2 = []
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute "conditions" cannot be specified when "conditions_v2" is specified`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					filters    = "[]"
					filters_v2 = []
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute "filters" cannot be specified when "filters_v2" is specified`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions    = "[]"
					actions_v2 = []
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute "actions" cannot be specified when "actions_v2" is specified`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions = "[]"
				`),
				ExpectError: acctest.ExpectLiteralError(`You must add an action for this alert to fire`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions_v2 = []
				`),
				ExpectError: acctest.ExpectLiteralError(`You must add an action for this alert to fire`),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions_v2 = [{ }]
				`),
				ExpectError: acctest.ExpectLiteralError(
					`Failed to convert action: [{Exactly one action must be set Exactly one action`,
					`must be set}]`,
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions_v2 = [{ notify_event = { } }]

					filters_v2 = [{ }]
				`),
				ExpectError: acctest.ExpectLiteralError(
					`Failed to convert filter: [{Exactly one filter must be set Exactly one filter`,
					`must be set}]`,
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions_v2 = [{ notify_event = { } }]

					conditions_v2 = [{ }]
				`),
				ExpectError: acctest.ExpectLiteralError(
					`Failed to convert condition: [{Exactly one condition must be set Exactly one`,
					`condition must be set}]`,
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions_v2 = [{ notify_event = { } }]

					conditions_v2 = [
						{ first_seen_event = {}, regression_event = {} },
					]
				`),
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "conditions_v2[0].first_seen_event" cannot be specified when`,
					`"conditions_v2[0].regression_event" is specified`,
				),
			},
		},
	})
}

func TestAccIssueAlertResource_basic(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		// TODO: Precheck acctest.TestOpsgenieIntegrationKey
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					conditions_v2 = [
						{ first_seen_event = {} },
						{ regression_event = {} },
						{ reappeared_event = {} },
						{ new_high_priority_issue = {} },
						{ existing_high_priority_issue = {} },
						{ event_frequency = { comparison_type = "count", value = 100, interval = "1h" } },
						{ event_frequency = { comparison_type = "percent", comparison_interval = "1w", value = 100, interval = "1h" } },
						{ event_unique_user_frequency = { comparison_type = "count", value = 100, interval = "1h" } },
						{ event_unique_user_frequency = { comparison_type = "percent", comparison_interval = "1w", value = 100, interval = "1h" } },
						{ event_frequency_percent = { comparison_type = "count", value = 100, interval = "1h" } },
						{ event_frequency_percent = { comparison_type = "percent", comparison_interval = "1w", value = 100, interval = "1h" } },
					]

					filters_v2 = [
						{ age_comparison = { comparison_type = "older", value = 10, time = "minute" } },
						{ issue_occurrences = { value = 10 } },
						{ assigned_to = { target_type = "Unassigned" } },
						{ assigned_to = { target_type = "Team", target_identifier = sentry_team.test.internal_id } },
						{ latest_adopted_release = { oldest_or_newest = "oldest", older_or_newer = "older", environment = "test" } },
						{ latest_release = {} },
						{ issue_category = { value = "Error" } },
						{ event_attribute = { attribute = "message", match = "CONTAINS", value = "test" } },
						{ event_attribute = { attribute = "message", match = "IS_SET" } },
						{ tagged_event = { key = "key", match = "CONTAINS", value = "value" } },
						{ tagged_event = { key = "key", match = "NOT_SET" } },
						{ level = { match = "EQUAL", level = "error" } },
					]

					actions_v2 = [
						{ notify_email = { target_type = "IssueOwners", fallthrough_type = "ActiveMembers" } },
						{ notify_email = { target_type = "Team", target_identifier = sentry_team.test.internal_id } },
						{ notify_event = { } },
						{
							notify_event_service = {
								service = "terraform-provider-sentry-ea4fdd"
							}
						},
						{
							notify_event_sentry_app = {
								sentry_app_installation_uuid = "d384d654-0e4c-447d-999c-a298fad579a7"

								settings = {
									teamId     = "5538c20b-37cf-4efd-b0aa-83c7f2e691f8"
									assigneeId = "b7afdd84-58b9-48ab-a682-9bb121d9dfbd"
									labelId    = "9f918fa3-9641-4522-950e-84dfb5c21099"
									projectId  = ""
									stateId    = "23e412bc-5abc-4812-916c-f91b4e21a060"
									priority   = "0"
								}
							}
						},
						{
							opsgenie_notify_team = {
								account  = sentry_integration_opsgenie.opsgenie.integration_id
								team     = sentry_integration_opsgenie.opsgenie.id
								priority = "P1"
							}
						},
						{
							pagerduty_notify_service = {
								account  = sentry_integration_pagerduty.pagerduty.integration_id
								service  = sentry_integration_pagerduty.pagerduty.id
								severity = "default"
							}
						},
						{
							slack_notify_service = {
								workspace = data.sentry_organization_integration.slack.id
								channel   = "#general"
								notes     = "Please <http://example.com|click here> for triage information"
							}
						},
						{
							slack_notify_service = {
								workspace = data.sentry_organization_integration.slack.id
								channel   = "#general"
								tags      = ["environment", "level"]
								notes     = "Please <http://example.com|click here> for triage information"
							}
						},
						{
							discord_notify_service = {
								server     = data.sentry_organization_integration.discord.id
								channel_id = "714123428994482189"
							}
						},
						{
							discord_notify_service = {
								server     = data.sentry_organization_integration.discord.id
								channel_id = "714123428994482189"
								tags       = ["environment", "level"]
							}
						},
						{
							github_create_ticket = {
								integration = data.sentry_organization_integration.github.id
								repo        = "terraform-provider-sentry"
								assignee    = "jianyuan"
								labels      = ["bug", "enhancement"]
							}
						},
						{
							azure_devops_create_ticket = {
								integration    = data.sentry_organization_integration.vsts.id
								project        = "123"
								work_item_type = "Microsoft.VSTS.WorkItemTypes.Task"
							}
						}
					]
				`) + fmt.Sprintf(`
					# Opsgenie
					data "sentry_organization_integration" "opsgenie" {
						organization = sentry_project.test.organization
						provider_key = "opsgenie"
						name         = "terraform-provider-sentry"
					}

					resource "sentry_integration_opsgenie" "opsgenie" {
						organization    = data.sentry_organization_integration.opsgenie.organization
						integration_id  = data.sentry_organization_integration.opsgenie.id
						team            = "issue-alert-team"
						integration_key = "%[1]s"
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

					# Slack
					data "sentry_organization_integration" "slack" {
						organization = sentry_project.test.organization
						provider_key = "slack"
						name         = "A2 Marketing"  # TODO: Use a real integration name
					}

					# Discord
					data "sentry_organization_integration" "discord" {
						organization = sentry_project.test.organization
						provider_key = "discord"
						name         = "jy's server"
					}

					# GitHub
					data "sentry_organization_integration" "github" {
						organization = sentry_project.test.organization
						provider_key = "github"
						name         = "jianyuan"
					}

					# Azure DevOps
					data "sentry_organization_integration" "vsts" {
						organization = sentry_project.test.organization
						provider_key = "vsts"
						name         = "jianyuanlee"
					}
				`, acctest.TestOpsgenieIntegrationKey),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"first_seen_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"regression_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"reappeared_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"new_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"existing_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.NotNull(),
								"comparison_type":     knownvalue.StringExact("count"),
								"comparison_interval": knownvalue.Null(),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.NotNull(),
								"comparison_type":     knownvalue.StringExact("percent"),
								"comparison_interval": knownvalue.StringExact("1w"),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_unique_user_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.NotNull(),
								"comparison_type":     knownvalue.StringExact("count"),
								"comparison_interval": knownvalue.Null(),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_unique_user_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.NotNull(),
								"comparison_type":     knownvalue.StringExact("percent"),
								"comparison_interval": knownvalue.StringExact("1w"),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency_percent": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.NotNull(),
								"comparison_type":     knownvalue.StringExact("count"),
								"comparison_interval": knownvalue.Null(),
								"value":               knownvalue.Float64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency_percent": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.NotNull(),
								"comparison_type":     knownvalue.StringExact("percent"),
								"comparison_interval": knownvalue.StringExact("1w"),
								"value":               knownvalue.Float64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"age_comparison": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":            knownvalue.NotNull(),
								"comparison_type": knownvalue.StringExact("older"),
								"value":           knownvalue.Int64Exact(10),
								"time":            knownvalue.StringExact("minute"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"issue_occurrences": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"value": knownvalue.Int64Exact(10),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"assigned_to": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":              knownvalue.NotNull(),
								"target_type":       knownvalue.StringExact("Unassigned"),
								"target_identifier": knownvalue.Null(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"assigned_to": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":              knownvalue.NotNull(),
								"target_type":       knownvalue.StringExact("Team"),
								"target_identifier": knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`)),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"latest_adopted_release": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":             knownvalue.NotNull(),
								"oldest_or_newest": knownvalue.StringExact("oldest"),
								"older_or_newer":   knownvalue.StringExact("older"),
								"environment":      knownvalue.StringExact("test"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"latest_release": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"issue_category": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"value": knownvalue.StringExact("Error"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_attribute": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":      knownvalue.NotNull(),
								"attribute": knownvalue.StringExact("message"),
								"match":     knownvalue.StringExact("CONTAINS"),
								"value":     knownvalue.StringExact("test"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_attribute": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":      knownvalue.NotNull(),
								"attribute": knownvalue.StringExact("message"),
								"match":     knownvalue.StringExact("IS_SET"),
								"value":     knownvalue.Null(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"tagged_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"key":   knownvalue.StringExact("key"),
								"match": knownvalue.StringExact("CONTAINS"),
								"value": knownvalue.StringExact("value"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"tagged_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"key":   knownvalue.StringExact("key"),
								"match": knownvalue.StringExact("NOT_SET"),
								"value": knownvalue.Null(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"level": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"match": knownvalue.StringExact("EQUAL"),
								"level": knownvalue.StringExact("error"),
							}),
						}),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"notify_email": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":              knownvalue.NotNull(),
								"target_type":       knownvalue.StringExact("IssueOwners"),
								"target_identifier": knownvalue.Null(),
								"fallthrough_type":  knownvalue.StringExact("ActiveMembers"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"notify_email": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":              knownvalue.NotNull(),
								"target_type":       knownvalue.StringExact("Team"),
								"target_identifier": knownvalue.NotNull(),
								"fallthrough_type":  knownvalue.Null(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"notify_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"notify_event_service": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.NotNull(),
								"service": knownvalue.StringExact("terraform-provider-sentry-ea4fdd"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"notify_event_sentry_app": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                         knownvalue.NotNull(),
								"sentry_app_installation_uuid": knownvalue.StringExact("d384d654-0e4c-447d-999c-a298fad579a7"),
								"settings": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"teamId":     knownvalue.StringExact("5538c20b-37cf-4efd-b0aa-83c7f2e691f8"),
									"assigneeId": knownvalue.StringExact("b7afdd84-58b9-48ab-a682-9bb121d9dfbd"),
									"labelId":    knownvalue.StringExact("9f918fa3-9641-4522-950e-84dfb5c21099"),
									"projectId":  knownvalue.StringExact(""),
									"stateId":    knownvalue.StringExact("23e412bc-5abc-4812-916c-f91b4e21a060"),
									"priority":   knownvalue.StringExact("0"),
								}),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"opsgenie_notify_team": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":     knownvalue.NotNull(),
								"account":  knownvalue.NotNull(),
								"team":     knownvalue.NotNull(),
								"priority": knownvalue.StringExact("P1"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"pagerduty_notify_service": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":     knownvalue.NotNull(),
								"account":  knownvalue.NotNull(),
								"service":  knownvalue.NotNull(),
								"severity": knownvalue.StringExact("default"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"slack_notify_service": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.NotNull(),
								"workspace":  knownvalue.NotNull(),
								"channel":    knownvalue.StringExact("#general"),
								"channel_id": knownvalue.NotNull(),
								"tags":       knownvalue.Null(),
								"notes":      knownvalue.StringExact("Please <http://example.com|click here> for triage information"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"slack_notify_service": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.NotNull(),
								"workspace":  knownvalue.NotNull(),
								"channel":    knownvalue.StringExact("#general"),
								"channel_id": knownvalue.NotNull(),
								"tags": knownvalue.SetExact([]knownvalue.Check{
									knownvalue.StringExact("environment"),
									knownvalue.StringExact("level"),
								}),
								"notes": knownvalue.StringExact("Please <http://example.com|click here> for triage information"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"discord_notify_service": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.NotNull(),
								"server":     knownvalue.NotNull(),
								"channel_id": knownvalue.NotNull(),
								"tags":       knownvalue.Null(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"discord_notify_service": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.NotNull(),
								"server":     knownvalue.NotNull(),
								"channel_id": knownvalue.NotNull(),
								"tags": knownvalue.SetExact([]knownvalue.Check{
									knownvalue.StringExact("environment"),
									knownvalue.StringExact("level"),
								}),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"github_create_ticket": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":        knownvalue.NotNull(),
								"integration": knownvalue.NotNull(),
								"repo":        knownvalue.StringExact("terraform-provider-sentry"),
								"assignee":    knownvalue.StringExact("jianyuan"),
								"labels": knownvalue.SetExact([]knownvalue.Check{
									knownvalue.StringExact("bug"),
									knownvalue.StringExact("enhancement"),
								}),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"azure_devops_create_ticket": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":           knownvalue.NotNull(),
								"integration":    knownvalue.NotNull(),
								"project":        knownvalue.StringExact("123"),
								"work_item_type": knownvalue.StringExact("Microsoft.VSTS.WorkItemTypes.Task"),
							}),
						}),
					})),
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert+"-updated", `
					conditions_v2 = [
						{ reappeared_event = {} },
						{ new_high_priority_issue = {} },
						{ existing_high_priority_issue = {} },
					]
					filters_v2 = []
					actions_v2 = [
						{ notify_email = { target_type = "IssueOwners", fallthrough_type = "NoOne" } },
					]
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"reappeared_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"new_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"existing_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.NotNull(),
							}),
						}),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters_v2"), knownvalue.ListExact([]knownvalue.Check{})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"notify_email": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":              knownvalue.NotNull(),
								"target_type":       knownvalue.StringExact("IssueOwners"),
								"target_identifier": knownvalue.Null(),
								"fallthrough_type":  knownvalue.StringExact("NoOne"),
							}),
						}),
					})),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
			},
		},
	})
}

func TestAccIssueAlertResource_emptyArray(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

	check := func(alert string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckIssueAlertExists(rn, &alertId),
			resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
				if alertId != value {
					return fmt.Errorf("expected %s, got %s", alertId, value)
				}
				return nil
			}),
		)
	}

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions_v2"), knownvalue.ListSizeExact(0)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters_v2"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					conditions_v2 = []

					actions = <<EOT
					[
						{
							"id": "sentry.mail.actions.NotifyEmailAction",
							"targetType": "IssueOwners"
						}
					]
					EOT
				`),
				Check: check(alert),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert+`updated`, `
					conditions_v2 = []

					actions = <<EOT
					[
						{
							"id": "sentry.mail.actions.NotifyEmailAction",
							"targetType": "IssueOwners"
						}
					]
					EOT
				`),
				Check: check(alert + "-updated"),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"updated")),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
			},
		},
	})
}

func TestAccIssueAlertResource_jsonValues(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

	check := func(alert string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckIssueAlertExists(rn, &alertId),
			resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
				if alertId != value {
					return fmt.Errorf("expected %s, got %s", alertId, value)
				}
				return nil
			}),
		)
	}

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig_jsonValues(team, project, alert),
				Check:  check(alert),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertConfig_jsonValues(team, project, alert+"-updated"),
				Check:  check(alert + "-updated"),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
			},
		},
	})
}

func TestAccIssueAlertResource_jsonValues_emptyArray(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

	check := func(alert string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckIssueAlertExists(rn, &alertId),
			resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
				if alertId != value {
					return fmt.Errorf("expected %s, got %s", alertId, value)
				}
				return nil
			}),
		)
	}

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.StringExact(`[]`)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions_v2"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters_v2"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					conditions = "[]"

					actions = <<EOT
					[
						{
							"id": "sentry.mail.actions.NotifyEmailAction",
							"targetType": "IssueOwners"
						}
					]
					EOT
				`),
				Check: check(alert),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert+`updated`, `
					conditions = "[]"

					actions = <<EOT
					[
						{
							"id": "sentry.mail.actions.NotifyEmailAction",
							"targetType": "IssueOwners"
						}
					]
					EOT
				`),
				Check: check(alert + "-updated"),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"updated")),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
			},
		},
	})
}

func TestAccIssueAlertResource_upgradeFromVersion(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.11.2",
					},
				},
				Config: testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = [
		{
			id = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
		},
		{
			id = "sentry.rules.conditions.regression_event.RegressionEventCondition"
		}
	]

	actions = [
		{
			id = "sentry.rules.actions.notify_event.NotifyEventAction"
		}
	]
}
`, team, project, alert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "id"),
					resource.TestCheckResourceAttrSet(rn, "internal_id"),
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "name", alert),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config: testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = jsonencode(
		[
			{
				id = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
			},
			{
				id = "sentry.rules.conditions.regression_event.RegressionEventCondition"
			}
		]
	)

	actions = jsonencode(
		[
			{
				"id": "sentry.rules.actions.notify_event.NotifyEventAction"
			}
		]
	)
}
`, team, project, alert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIssueAlertExists(rn, &alertId),
					resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
						if alertId != value {
							return fmt.Errorf("expected %s, got %s", alertId, value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "name", alert),
					resource.TestCheckResourceAttr(rn, "action_match", "any"),
					resource.TestCheckResourceAttr(rn, "filter_match", "any"),
					resource.TestCheckResourceAttr(rn, "frequency", "30"),
					resource.TestCheckResourceAttrSet(rn, "conditions"),
					resource.TestCheckResourceAttrSet(rn, "actions"),
				),
			},
		},
	})
}

func testAccCheckIssueAlertDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_issue_alert" {
			continue
		}

		ctx := context.Background()
		alert, resp, err := acctest.SharedClient.IssueAlerts.Get(ctx, rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], rs.Primary.ID)
		if err == nil {
			if alert != nil {
				return fmt.Errorf("issue alert %q still exists", rs.Primary.ID)
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckIssueAlertExists(n string, alertId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		var resolvedAlertId string
		// Support schema v1 and below
		if value, ok := rs.Primary.Attributes["internal_id"]; ok {
			resolvedAlertId = value
		} else {
			resolvedAlertId = rs.Primary.ID
		}

		ctx := context.Background()
		gotAlert, _, err := acctest.SharedClient.IssueAlerts.Get(ctx, rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], resolvedAlertId)
		if err != nil {
			return err
		}
		*alertId = sentry.StringValue(gotAlert.ID)
		return nil
	}
}

func testAccIssueAlertConfig(team string, project string, alert string, extras string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	%[4]s
}
`, team, project, alert, extras)
}

func testAccIssueAlertConfig_jsonValues(team string, project string, alert string) string {
	return testAccIssueAlertConfig(team, project, alert, `
	conditions = <<EOT
[
	{
		"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
		"name": "ignored"
	},
	{
		"id": "sentry.rules.conditions.regression_event.RegressionEventCondition"
	},
	{
		"id": "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
		"value": 100,
		"comparisonType": "count",
		"interval": "1h"
	},
	{
		"id": "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
		"value": 100,
		"comparisonType": "count",
		"interval": "1h"
	},
	{
		"id": "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
		"value": 50.0,
		"comparisonType": "count",
		"interval": "1h"
	}
]
EOT

	filters = <<EOT
[
	{
		"id": "sentry.rules.filters.age_comparison.AgeComparisonFilter",
		"value": 10,
		"time": "minute",
		"comparison_type": "older"
	},
	{
		"id": "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
		"value": 10
	},
	{
		"id": "sentry.rules.filters.assigned_to.AssignedToFilter",
		"targetType": "Team",
		"targetIdentifier": ${parseint(sentry_team.test.team_id, 10)}
	},
	{
		"id": "sentry.rules.filters.latest_release.LatestReleaseFilter"
	},
	{
		"id": "sentry.rules.filters.event_attribute.EventAttributeFilter",
		"attribute": "message",
		"match": "co",
		"value": "test"
	},
	{
		"id": "sentry.rules.filters.tagged_event.TaggedEventFilter",
		"key": "test",
		"match": "co",
		"value": "test"
	},
	{
		"id": "sentry.rules.filters.level.LevelFilter",
		"match": "eq",
		"level": "50"
	}
]
EOT

	actions = <<EOT
[
	{
		"id": "sentry.mail.actions.NotifyEmailAction",
		"targetType": "IssueOwners",
		"targetIdentifier": ""
	},
	{
		"id": "sentry.mail.actions.NotifyEmailAction",
		"targetType": "Team",
		"targetIdentifier": ${parseint(sentry_team.test.team_id, 10)}
	},
	{
		"id": "sentry.rules.actions.notify_event.NotifyEventAction"
	}
]
EOT
`)
}
