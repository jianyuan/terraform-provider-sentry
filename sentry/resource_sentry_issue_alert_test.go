package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/deste-org/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryIssueAlert_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-issue-alert")
	rn := "sentry_issue_alert.test"

	var alertID string

	check := func(alertName string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryIssueAlertExists(rn, &alertID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "project", projectName),
			resource.TestCheckResourceAttr(rn, "projects.#", "1"),
			resource.TestCheckResourceAttr(rn, "projects.0", projectName),
			resource.TestCheckResourceAttr(rn, "name", alertName),
			resource.TestCheckResourceAttr(rn, "environment", ""),
			resource.TestCheckResourceAttr(rn, "action_match", "any"),
			resource.TestCheckResourceAttr(rn, "filter_match", "any"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &alertID),
			// Conditions
			resource.TestCheckResourceAttr(rn, "conditions.#", "5"),
			resource.TestCheckResourceAttr(rn, "conditions.0.id", "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.1.id", "sentry.rules.conditions.regression_event.RegressionEventCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.2.id", "sentry.rules.conditions.event_frequency.EventFrequencyCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.3.id", "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"),
			resource.TestCheckResourceAttr(rn, "conditions.4.id", "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id": "sentry.rules.conditions.regression_event.RegressionEventCondition",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":             "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
				"value":          "100",
				"comparisonType": "count",
				"interval":       "1h",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":             "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
				"value":          "100",
				"comparisonType": "count",
				"interval":       "1h",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "conditions.*", map[string]string{
				"id":             "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
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
				"value":           "10",
				"time":            "minute",
				"comparison_type": "older",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":    "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
				"value": "10",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":         "sentry.rules.filters.assigned_to.AssignedToFilter",
				"targetType": "Team",
			}),
			resource.TestCheckResourceAttrPair(rn, "filters.2.targetIdentifier", "sentry_team.test", "team_id"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id": "sentry.rules.filters.latest_release.LatestReleaseFilter",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":        "sentry.rules.filters.event_attribute.EventAttributeFilter",
				"attribute": "message",
				"match":     "co",
				"value":     "test",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":    "sentry.rules.filters.tagged_event.TaggedEventFilter",
				"key":   "test",
				"match": "co",
				"value": "test",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "filters.*", map[string]string{
				"id":    "sentry.rules.filters.level.LevelFilter",
				"match": "eq",
				"level": "50",
			}),
			// Actions
			resource.TestCheckResourceAttr(rn, "actions.#", "3"),
			resource.TestCheckResourceAttr(rn, "actions.0.id", "sentry.mail.actions.NotifyEmailAction"),
			resource.TestCheckResourceAttr(rn, "actions.1.id", "sentry.mail.actions.NotifyEmailAction"),
			resource.TestCheckResourceAttr(rn, "actions.2.id", "sentry.rules.actions.notify_event.NotifyEventAction"),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"targetType":       "IssueOwners",
				"targetIdentifier": "",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id":               "sentry.mail.actions.NotifyEmailAction",
				"targetType":       "Team",
				"targetIdentifier": "",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(rn, "actions.*", map[string]string{
				"id": "sentry.rules.actions.notify_event.NotifyEventAction",
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
				ResourceName: rn,
				ImportState:  true,
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

		org, project, id, err := splitSentryAlertID(rs.Primary.ID)
		if err != nil {
			return err
		}

		ctx := context.Background()
		alert, resp, err := client.IssueAlerts.Get(ctx, org, project, id)
		if err == nil {
			if alert != nil {
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

func testAccCheckSentryIssueAlertExists(n string, alertID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		org, project, id, err := splitSentryAlertID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotAlert, _, err := client.IssueAlerts.Get(ctx, org, project, id)
		if err != nil {
			return err
		}
		*alertID = sentry.StringValue(gotAlert.ID)
		return nil
	}
}

func testAccSentryIssueAlertConfig(teamName, projectName, alertName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + fmt.Sprintf(`
resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[1]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = [
		{
			id = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
		},
		{
			id = "sentry.rules.conditions.regression_event.RegressionEventCondition"
		},
		{
			id             = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
			value          = 100
			comparisonType = "count"
			interval       = "1h"
		},
		{
			id             = "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
			value          = 100
			comparisonType = "count"
			interval       = "1h"
		},
		{
			id             = "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
			value          = "50.0"
			comparisonType = "count"
			interval       = "1h"
		},
	]

	filters = [
		{
			id              = "sentry.rules.filters.age_comparison.AgeComparisonFilter"
			value           = 10
			time            = "minute"
			comparison_type = "older"
		},
		{
			id    = "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter"
			value = 10
		},
		{
			id               = "sentry.rules.filters.assigned_to.AssignedToFilter"
			targetType       = "Team"
			targetIdentifier = sentry_team.test.team_id
		},
		{
			id = "sentry.rules.filters.latest_release.LatestReleaseFilter"
		},
		{
			id        = "sentry.rules.filters.event_attribute.EventAttributeFilter"
			attribute = "message"
			match     = "co"
			value     = "test"
		},
		{
			id    = "sentry.rules.filters.tagged_event.TaggedEventFilter"
			key   = "test"
			match = "co"
			value = "test"
		},
		{
			id    = "sentry.rules.filters.level.LevelFilter"
			match = "eq"
			level = "50"
		}
	]

	actions = [
		{
			id               = "sentry.mail.actions.NotifyEmailAction"
			targetType       = "IssueOwners"
			targetIdentifier = ""
		},
		{
			id               = "sentry.mail.actions.NotifyEmailAction"
			targetType       = "Team"
			targetIdentifier = sentry_team.test.team_id
		},
		{
			id = "sentry.rules.actions.notify_event.NotifyEventAction"
		}
	]
}
	`, alertName)
}
