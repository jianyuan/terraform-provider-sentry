package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryIssueAlert_basic(t *testing.T) {
	var rule sentry.IssueAlert

	rn := "sentry_issue_alert.test_issue_alert"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryIssueAlertConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryIssueAlertExists(rn, &rule),
					resource.TestCheckResourceAttr(rn, "organization", testOrganization),
					resource.TestCheckResourceAttr(rn, "conditions.#", "5"),
					resource.TestCheckResourceAttr(rn, "conditions.0.%", "2"),
					resource.TestCheckResourceAttr(rn, "conditions.0.id", "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"),
					resource.TestCheckResourceAttr(rn, "conditions.0.name", "A new issue is created"),
					resource.TestCheckResourceAttr(rn, "conditions.1.%", "2"),
					resource.TestCheckResourceAttr(rn, "conditions.1.id", "sentry.rules.conditions.regression_event.RegressionEventCondition"),
					resource.TestCheckResourceAttr(rn, "conditions.1.name", "The issue changes state from resolved to unresolved"),
					// TODO: conditions
					// TODO: filters
					// TODO: actions
					resource.TestCheckResourceAttr(rn, "action_match", "any"),
					resource.TestCheckResourceAttr(rn, "filter_match", "any"),
					resource.TestCheckResourceAttr(rn, "name", "Test rule"),
					resource.TestCheckResourceAttr(rn, "environment", ""),
					resource.TestCheckResourceAttr(rn, "project", "test-project"),
				),
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

var testAccSentryIssueAlertConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
	organization = "%s"
	name         = "Test team"
}

resource "sentry_project" "test_project" {
	organization = sentry_team.test_team.organization
	team         = sentry_team.test_team.id
	name         = "Test project"
	platform     = "go"
}

resource "sentry_issue_alert" "test_issue_alert" {
	organization = sentry_project.test_project.organization
	project      = sentry_project.test_project.id
	name         = "Test rule"

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
			interval       = "1h"
			id             = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
			comparisonType = "count"
			value          = 100
			name           = "The issue is seen more than 100 times in 1h"
		},
		{
			interval       = "1h"
			id             = "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
			comparisonType = "count"
			value          = 100
			name           = "The issue is seen by more than 100 users in 1h"
		},
		{
			interval       = "1h"
			id             = "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
			comparisonType = "count"
			value          = 100
			name           = "The issue affects more than 100.0 percent of sessions in 1h"
		},
	]

	filters = [
		{
			comparison_type = "older"
			time            = "minute"
			id              = "sentry.rules.filters.age_comparison.AgeComparisonFilter"
			value           = 10
			name            = "The issue is older than 10 minute"
		},
		{
			id    = "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter"
			value = 10
			name  = "The issue has happened at least 10 times"
		},
		{
			targetType       = "Team"
			id               = "sentry.rules.filters.assigned_to.AssignedToFilter"
			targetIdentifier = sentry_team.test_team.team_id
			name             = "The issue is assigned to Team"
		},
		{
			id   = "sentry.rules.filters.latest_release.LatestReleaseFilter"
			name = "The event is from the latest release"
		},
		{
			attribute = "message"
			match     = "co"
			id        = "sentry.rules.filters.event_attribute.EventAttributeFilter"
			value     = "test"
			name      = "The event's message value contains test"
		},
		{
			match = "co"
			id    = "sentry.rules.filters.tagged_event.TaggedEventFilter"
			key   = "test"
			value = "test"
			name  = "The event's tags match test contains test"
		},
		{
			level = "50"
			match = "eq"
			id    = "sentry.rules.filters.level.LevelFilter"
			name  = "The event's level is equal to fatal"
		}
	]

	actions = [
		{
			targetType       = "IssueOwners"
			id               = "sentry.mail.actions.NotifyEmailAction"
			targetIdentifier = ""
			name             = "Send a notification to IssueOwners"
		},
		{
			targetType       = "Team"
			id               = "sentry.mail.actions.NotifyEmailAction"
			targetIdentifier = sentry_team.test_team.team_id
			name             = "Send a notification to Team"
		},
		{
			targetType       = "Member"
			id               = "sentry.mail.actions.NotifyEmailAction"
			targetIdentifier = 94401
			name             = "Send a notification to Member"
		},
		{
			id   = "sentry.rules.actions.notify_event.NotifyEventAction"
			name = "Send a notification (for all legacy integrations)"
		}
	]
}
`, testOrganization)
