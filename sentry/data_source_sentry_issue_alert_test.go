package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryIssueAlertDataSource_basic(t *testing.T) {
	var alert sentry.IssueAlert
	var alertCopy sentry.IssueAlert

	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-issue-alert")
	rn := "sentry_issue_alert.test"
	dn := "data.sentry_issue_alert.test"
	rnCopy := "sentry_issue_alert.test_copy"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryIssueAlertDataSourceConfig(teamSlug, projectName, alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryIssueAlertExists(rn, &alert),
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
					testAccCheckSentryIssueAlertExists(rnCopy, &alertCopy),
					resource.TestCheckResourceAttrPair(rnCopy, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(rnCopy, "project", rn, "project"),
					resource.TestCheckResourceAttrWith(rnCopy, "internal_id", func(v string) error {
						want := sentry.StringValue(alertCopy.ID)
						if v != want {
							return fmt.Errorf("got issue alert ID %s; want %s", v, want)
						}
						return nil
					}),
					resource.TestCheckResourceAttrPair(rnCopy, "conditions", rn, "conditions"),
					resource.TestCheckResourceAttrPair(rnCopy, "filters", rn, "filters"),
					resource.TestCheckResourceAttrPair(rnCopy, "actions", rn, "actions"),
					resource.TestCheckResourceAttrPair(rnCopy, "action_match", rn, "action_match"),
					resource.TestCheckResourceAttrPair(rnCopy, "filter_match", rn, "filter_match"),
					resource.TestCheckResourceAttrPair(rnCopy, "frequency", rn, "frequency"),
					resource.TestCheckResourceAttrWith(rnCopy, "name", func(v string) error {
						want := sentry.StringValue(alertCopy.Name)
						if v != want {
							return fmt.Errorf("got name ID %s; want %s", v, want)
						}
						return nil
					}),
					resource.TestCheckResourceAttrPair(rnCopy, "environment", rn, "environment"),
					resource.TestCheckResourceAttrPair(rnCopy, "actions", rn, "actions"),
				),
			},
		},
	})
}

func testAccSentryIssueAlertDataSourceConfig(teamSlug, projectName, alertName string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	team         = sentry_team.test.id
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
			targetIdentifier = sentry_team.test.team_id
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
			targetIdentifier = sentry_team.test.team_id
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
	`, teamSlug, projectName, alertName)
}
