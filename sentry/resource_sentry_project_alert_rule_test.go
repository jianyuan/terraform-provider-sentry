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

func TestAccSentryProjectAlertRule_basic(t *testing.T) {
	var alertRule sentry.MetricAlert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryProjectAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryProjectAlertRuleExists("sentry_alert_rule.test_alert_rule", &alertRule),
					resource.TestCheckResourceAttr("sentry_alert_rule.test_alert_rule", "name", "Test alert rule"),
					resource.TestCheckResourceAttr("sentry_alert_rule.test_alert_rule", "environment", ""),
					resource.TestCheckResourceAttr("sentry_alert_rule.test_alert_rule", "dataset", "transactions"),
				),
			},
		},
	})
}

func testAccCheckSentryProjectAlertRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_alert_rule" {
			continue
		}

		ctx := context.Background()
		alertRules, resp, err := client.MetricAlerts.List(ctx, testOrganization, rs.Primary.Attributes["project"])
		if err == nil {
			for _, alertRule := range alertRules {
				if alertRule.ID == rs.Primary.ID {
					return errors.New("Project alert rule still exists")
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

func testAccCheckSentryProjectAlertRuleExists(n string, alertRule *sentry.MetricAlert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		sentryAlertRules, _, err := client.MetricAlerts.List(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
		)
		if err != nil {
			return err
		}
		for _, sentryAlertRule := range sentryAlertRules {
			if sentryAlertRule.ID == rs.Primary.ID {
				*alertRule = *sentryAlertRule
				break
			}
		}
		return nil
	}
}

var testAccSentryProjectAlertRuleConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
	organization = "%s"
	name		= "Test team"	
}

resource "sentry_project" "test_project" {
	organization = "%s"
	team = sentry_team.test_team.id
	name = "Test project"
	platform = "go"
}

resource "sentry_alert_rule" "test_alert_rule" {
	organization      = "%s"
	project           = sentry_project.test_project.id
	name              = "Test alert rule"
	dataset           = "transactions"
	query             = "http.url:http://testservice.com/stats"
	time_window       = 50.0
	aggregate         = "p50(transaction.duration)"
	threshold_type    = 0
	resolve_threshold = 100.0
  
	triggers {
	  actions           = []
	  alert_threshold   = 1000
	  label             = "critical"
	  resolve_threshold = 100.0
	  threshold_type    = 0
	}
  
	triggers {
	  actions = []
	  alert_threshold   = 500
	  label             = "warning"
	  resolve_threshold = 100.0
	  threshold_type    = 0
	}
  
	projects = [
		sentry_project.test_project.id,
	]
}

`, testOrganization, testOrganization, testOrganization)
