package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryProjectAPMRule_basic(t *testing.T) {
	var apmRule sentry.APMRule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryProjectAPMRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectAPMRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryProjectAPMRuleExists("sentry_apm_rule.test_apm_rule", &apmRule),
					resource.TestCheckResourceAttr("sentry_apm_rule.test_apm_rule", "name", "Test apm rule"),
					resource.TestCheckResourceAttr("sentry_apm_rule", "environment", ""),
					resource.TestCheckResourceAttr("sentry_apm_rule", "dataset", "transactions"),
				),
			},
		},
	})
}

func testAccCheckSentryProjectAPMRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_apm_rule" {
			continue
		}

		apmRules, resp, err := client.APMRules.List(testOrganization, rs.Primary.Attributes["project"])
		if err == nil {
			for _, apmRule := range apmRules {
				if apmRule.ID == rs.Primary.ID {
					return errors.New("Project apm rule still exists")
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

func testAccCheckSentryProjectAPMRuleExists(n string, apmRule *sentry.APMRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		sentryApmRules, _, err := client.APMRules.List(
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
		)
		if err != nil {
			return err
		}
		for _, sentryApmRule := range sentryApmRules {
			if sentryApmRule.ID == rs.Primary.ID {
				*apmRule = sentryApmRule
				break
			}
		}
		return nil
	}
}

var testAccSentryProjectAPMRuleConfig = fmt.Sprintf(`
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

resource "sentry_apm_rule" "test_apm_rule" {
	organization      = "%s"
	project           = sentry_project.test_project.id
	name              = "Test apm rule"
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
  
	projects = ["Test project"]
}

`, testOrganization, testOrganization, testOrganization)
