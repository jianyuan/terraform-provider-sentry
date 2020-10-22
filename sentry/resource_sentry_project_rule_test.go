package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryProjectRule_basic(t *testing.T) {
	var rule sentry.Rule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryProjectRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryProjectRuleExists("sentry_rule.test_rule", &rule),
					resource.TestCheckResourceAttr("sentry_rule.test_rule", "name", "Test rule"),
					resource.TestCheckResourceAttr("sentry_rule.test_rule", "environment", ""),
					resource.TestCheckResourceAttr("sentry_rule.test_rule", "frequency", "30"),
				),
			},
		},
	})
}

func testAccCheckSentryProjectRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_rule" {
			continue
		}

		rules, resp, err := client.Rules.List(testOrganization, rs.Primary.Attributes["project"])
		if err == nil {
			for _, rule := range rules {
				if rule.ID == rs.Primary.ID {
					return errors.New("Project rule still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckSentryProjectRuleExists(n string, rule *sentry.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		sentryRules, _, err := client.Rules.List(
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
		)
		if err != nil {
			return err
		}
		for _, sentryRule := range sentryRules {
			if sentryRule.ID == rs.Primary.ID {
				*rule = sentryRule
				break
			}
		}
		return nil
	}
}

var testAccSentryProjectRuleConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
	organization = "%s"
	name         = "Test team"
}

resource "sentry_project" "test_project" {
	organization = "%s"
	team         = "${sentry_team.test_team.id}"
	name         = "Test project"
	platform     = "go"
}

resource "sentry_rule" "test_rule" {
	organization = "%s"
	project      = "${sentry_project.test_project.id}"
	name         = "Test rule"

	conditions = [
		{
			id   = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
			name = "A new issue is created"
		}
	]

	filters = [
		{
			id         = "sentry.rules.filters.assigned_to.AssignedToFilter"
			targetType = "Unassigned"
		}
	]

	actions = [
		{
			id               = "sentry.mail.actions.NotifyEmailAction"
			name             = "Send an email to IssueOwners"
			targetIdentifier = ""
			targetType       = "IssueOwners"
		}
	]
}
`, testOrganization, testOrganization, testOrganization)
