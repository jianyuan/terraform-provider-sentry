package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSentryIssueAlertDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-issue-alert")
	rn := "sentry_issue_alert.test"
	dn := "data.sentry_issue_alert.test"
	rnCopy := "sentry_issue_alert.test_copy"

	var alertID string
	var alertCopyID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryIssueAlertDataSourceConfig(teamName, projectName, alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryIssueAlertExists(rn, &alertID),
					resource.TestCheckResourceAttr(dn, "organization", testOrganization),
					resource.TestCheckResourceAttr(dn, "project", projectName),
					resource.TestCheckResourceAttrPair(dn, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(dn, "project", rn, "project"),
					resource.TestCheckResourceAttrPair(dn, "internal_id", rn, "internal_id"),
					resource.TestCheckResourceAttrPair(dn, "conditions", rn, "conditions"),
					resource.TestCheckResourceAttrPair(dn, "filters", rn, "filters"),
					resource.TestCheckResourceAttrPair(dn, "actions", rn, "actions"),
					resource.TestCheckResourceAttrPair(dn, "action_match", rn, "action_match"),
					resource.TestCheckResourceAttrPair(dn, "filter_match", rn, "filter_match"),
					resource.TestCheckResourceAttrPair(dn, "frequency", rn, "frequency"),
					resource.TestCheckResourceAttrPair(dn, "name", rn, "name"),
					resource.TestCheckResourceAttrPair(dn, "environment", rn, "environment"),
					resource.TestCheckResourceAttrPair(dn, "actions", rn, "actions"),
					testAccCheckSentryIssueAlertExists(rnCopy, &alertCopyID),
					resource.TestCheckResourceAttrPair(rnCopy, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(rnCopy, "project", rn, "project"),
					resource.TestCheckResourceAttrPtr(rnCopy, "internal_id", &alertCopyID),
					resource.TestCheckResourceAttrPair(rnCopy, "conditions", rn, "conditions"),
					resource.TestCheckResourceAttrPair(rnCopy, "filters", rn, "filters"),
					resource.TestCheckResourceAttrPair(rnCopy, "actions", rn, "actions"),
					resource.TestCheckResourceAttrPair(rnCopy, "action_match", rn, "action_match"),
					resource.TestCheckResourceAttrPair(rnCopy, "filter_match", rn, "filter_match"),
					resource.TestCheckResourceAttrPair(rnCopy, "frequency", rn, "frequency"),
					resource.TestCheckResourceAttr(rnCopy, "name", alertName+"-copy"),
					resource.TestCheckResourceAttrPair(rnCopy, "environment", rn, "environment"),
					resource.TestCheckResourceAttrPair(rnCopy, "actions", rn, "actions"),
				),
			},
		},
	})
}

func testAccSentryIssueAlertDataSourceConfig(teamName, projectName, alertName string) string {
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
			id               = "sentry.mail.actions.NotifyEmailAction"
			targetType       = "Member"
			targetIdentifier = 94401
		},
		{
			id = "sentry.rules.actions.notify_event.NotifyEventAction"
		}
	]
}

data "sentry_issue_alert" "test" {
	organization = sentry_issue_alert.test.organization
	project      = sentry_issue_alert.test.project
	internal_id  = sentry_issue_alert.test.internal_id
}

resource "sentry_issue_alert" "test_copy" {
	organization = data.sentry_issue_alert.test.organization
	project      = data.sentry_issue_alert.test.project
	name         = "${data.sentry_issue_alert.test.name}-copy"

	action_match = data.sentry_issue_alert.test.action_match
	filter_match = data.sentry_issue_alert.test.filter_match
	frequency    = data.sentry_issue_alert.test.frequency

	conditions = data.sentry_issue_alert.test.conditions
	filters    = data.sentry_issue_alert.test.filters
	actions    = data.sentry_issue_alert.test.actions
}
	`, alertName)
}
