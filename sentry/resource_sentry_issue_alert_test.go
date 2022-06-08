package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryIssueAlert_basic(t *testing.T) {
	var rule sentry.IssueAlert

	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-issue-alert")
	rn := "sentry_issue_alert.test_issue_alert"

	check := func(alertName string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryIssueAlertExists(rn, &rule),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "action_match", "any"),
			resource.TestCheckResourceAttr(rn, "filter_match", "any"),
			resource.TestCheckResourceAttr(rn, "name", alertName),
			resource.TestCheckResourceAttr(rn, "environment", ""),
			resource.TestCheckResourceAttr(rn, "project", projectName),
			// Conditions
			resource.TestCheckResourceAttr(rn, "conditions.#", "5"),
			resource.TestCheckResourceAttr(rn, "conditions.0.id", "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.1.id", "sentry.rules.conditions.regression_event.RegressionEventCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.2.id", "sentry.rules.conditions.event_frequency.EventFrequencyCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.3.id", "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.4.id", "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":   "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
				"name": "A new issue is created",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":   "sentry.rules.conditions.regression_event.RegressionEventCondition",
				"name": "The issue changes state from resolved to unresolved",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":             "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
				"name":           "The issue is seen more than 100 times in 1h",
				"value":          "100",
				"comparisonType": "count",
				"interval":       "1h",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":             "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
				"name":           "The issue is seen by more than 100 users in 1h",
				"value":          "100",
				"comparisonType": "count",
				"interval":       "1h",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":             "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
				"name":           "The issue affects more than 50.0 percent of sessions in 1h",
				"value":          "50.0",
				"comparisonType": "count",
				"interval":       "1h",
			}),
			// Filters
			resource.TestCheckResourceAttr(rn, "filters.#", "7"),
			resource.TestCheckResourceAttr(rn, "filters.0.id", "sentry.rules.filters.age_comparison.AgeComparisonFilter"),
			resource.TestCheckResourceAttr(rn, "filters.1.id", "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter"),
			resource.TestCheckResourceAttr(rn, "filters.2.id", "sentry.rules.filters.assigned_to.AssignedToFilter"),
			resource.TestCheckResourceAttr(rn, "filters.3.id", "sentry.rules.filters.latest_release.LatestReleaseFilter"),
			resource.TestCheckResourceAttr(rn, "filters.4.id", "sentry.rules.filters.event_attribute.EventAttributeFilter"),
			resource.TestCheckResourceAttr(rn, "filters.5.id", "sentry.rules.filters.tagged_event.TaggedEventFilter"),
			resource.TestCheckResourceAttr(rn, "filters.6.id", "sentry.rules.filters.level.LevelFilter"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":              "sentry.rules.filters.age_comparison.AgeComparisonFilter",
				"name":            "The issue is older than 10 minute",
				"value":           "10",
				"time":            "minute",
				"comparison_type": "older",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":    "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
				"name":  "The issue has happened at least 10 times",
				"value": "10",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":         "sentry.rules.filters.assigned_to.AssignedToFilter",
				"name":       "The issue is assigned to Team",
				"targetType": "Team",
			}),
			resource.TestCheckResourceAttrPair(rn, "filters.2.targetIdentifier", "sentry_team.test_team", "team_id"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":   "sentry.rules.filters.latest_release.LatestReleaseFilter",
				"name": "The event is from the latest release",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":        "sentry.rules.filters.event_attribute.EventAttributeFilter",
				"name":      "The event's message value contains test",
				"attribute": "message",
				"match":     "co",
				"value":     "test",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":    "sentry.rules.filters.tagged_event.TaggedEventFilter",
				"name":  "The event's tags match test contains test",
				"key":   "test",
				"match": "co",
				"value": "test",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":    "sentry.rules.filters.level.LevelFilter",
				"name":  "The event's level is equal to fatal",
				"match": "eq",
				"level": "50",
			}),
			// Actions
			resource.TestCheckResourceAttr(rn, "actions.#", "4"),
			resource.TestCheckResourceAttr(rn, "actions.0.id", "sentry.mail.actions.NotifyEmailAction"),
			resource.TestCheckResourceAttr(rn, "actions.1.id", "sentry.mail.actions.NotifyEmailAction"),
			resource.TestCheckResourceAttr(rn, "actions.2.id", "sentry.mail.actions.NotifyEmailAction"),
			resource.TestCheckResourceAttr(rn, "actions.3.id", "sentry.rules.actions.notify_event.NotifyEventAction"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"name":             "Send a notification to IssueOwners",
				"targetType":       "IssueOwners",
				"targetIdentifier": "",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"name":             "Send a notification to Team",
				"targetType":       "Team",
				"targetIdentifier": "",
			}),
			resource.TestCheckResourceAttrPair(rn, "actions.1.targetIdentifier", "sentry_team.test_team", "team_id"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"name":             "Send a notification to Member",
				"targetType":       "Member",
				"targetIdentifier": "94401", // TODO
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id":   "sentry.rules.actions.notify_event.NotifyEventAction",
				"name": "Send a notification (for all legacy integrations)",
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryIssueAlertConfig(teamName, projectName, alertName),
				Check:  check(alertName),
			},
			{
				Config: testAccSentryIssueAlertConfig(teamName, projectName, alertName+"-renamed"),
				Check:  check(alertName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryIssueAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_issue_alert" {
			continue
		}

		org, project, id, err := splitThreePartID(rs.Primary.ID, "organization-slug", "project-slug", "id")
		if err != nil {
			return err
		}

		ctx := context.Background()
		rule, resp, err := client.IssueAlerts.Get(ctx, org, project, id)
		if err == nil {
			if rule != nil {
				return errors.New("issue alert still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckSentryIssueAlertExists(n string, rule *sentry.IssueAlert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		org, project, id, err := splitThreePartID(rs.Primary.ID, "organization-slug", "project-slug", "id")
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotRule, _, err := client.IssueAlerts.Get(ctx, org, project, id)
		if err != nil {
			return err
		}
		*rule = *gotRule
		return nil
	}
}

func testAccSentryIssueAlertConfig(teamName, projectName, ruleName string) string {
	return fmt.Sprintf(`
data "sentry_organization" "test_organization" {
	slug = "%[1]s"
}

resource "sentry_team" "test_team" {
	organization = data.sentry_organization.test_organization.id
	name         = "%[2]s"
}

resource "sentry_project" "test_project" {
	organization = sentry_team.test_team.organization
	team         = sentry_team.test_team.id
	name         = "%[3]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test_issue_alert" {
	organization = sentry_project.test_project.organization
	project      = sentry_project.test_project.id
	name         = "%[4]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = [
		{
			id   = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
			name = "A new issue is created"
		},
		{
			id   = "sentry.rules.conditions.regression_event.RegressionEventCondition"
			name = "The issue changes state from resolved to unresolved"
		},
		{
			id             = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
			name           = "The issue is seen more than 100 times in 1h"
			value          = 100
			comparisonType = "count"
			interval       = "1h"
		},
		{
			id             = "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
			name           = "The issue is seen by more than 100 users in 1h"
			value          = 100
			comparisonType = "count"
			interval       = "1h"
		},
		{
			id             = "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
			name           = "The issue affects more than 50.0 percent of sessions in 1h"
			value          = 50.0
			comparisonType = "count"
			interval       = "1h"
		},
	]

	filters = [
		{
			id              = "sentry.rules.filters.age_comparison.AgeComparisonFilter"
			name            = "The issue is older than 10 minute"
			value           = 10
			time            = "minute"
			comparison_type = "older"
		},
		{
			id    = "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter"
			name  = "The issue has happened at least 10 times"
			value = 10
		},
		{
			id               = "sentry.rules.filters.assigned_to.AssignedToFilter"
			name             = "The issue is assigned to Team"
			targetType       = "Team"
			targetIdentifier = sentry_team.test_team.team_id
		},
		{
			id   = "sentry.rules.filters.latest_release.LatestReleaseFilter"
			name = "The event is from the latest release"
		},
		{
			id        = "sentry.rules.filters.event_attribute.EventAttributeFilter"
			name      = "The event's message value contains test"
			attribute = "message"
			match     = "co"
			value     = "test"
		},
		{
			id    = "sentry.rules.filters.tagged_event.TaggedEventFilter"
			name  = "The event's tags match test contains test"
			key   = "test"
			match = "co"
			value = "test"
		},
		{
			id    = "sentry.rules.filters.level.LevelFilter"
			name  = "The event's level is equal to fatal"
			match = "eq"
			level = "50"
		}
	]

	actions = [
		{
			id               = "sentry.mail.actions.NotifyEmailAction"
			name             = "Send a notification to IssueOwners"
			targetType       = "IssueOwners"
			targetIdentifier = ""
		},
		{
			id               = "sentry.mail.actions.NotifyEmailAction"
			name             = "Send a notification to Team"
			targetType       = "Team"
			targetIdentifier = sentry_team.test_team.team_id
		},
		{
			id               = "sentry.mail.actions.NotifyEmailAction"
			name             = "Send a notification to Member"
			targetType       = "Member"
			targetIdentifier = 94401
		},
		{
			id   = "sentry.rules.actions.notify_event.NotifyEventAction"
			name = "Send a notification (for all legacy integrations)"
		}
	]
}
	`, testOrganization, teamName, projectName, ruleName)
}
