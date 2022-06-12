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

func TestAccSentryMetricAlert_basic(t *testing.T) {
	var alert sentry.MetricAlert

	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-issue-alert")
	rn := "sentry_metric_alert.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryMetricAlertConfig(teamSlug, projectName, alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryMetricAlertExists(rn, &alert),
					resource.TestCheckResourceAttr(rn, "organization", testOrganization),
					resource.TestCheckResourceAttr(rn, "project", projectName),
					resource.TestCheckResourceAttr(rn, "projects.#", "1"),
					resource.TestCheckResourceAttr(rn, "projects.0", projectName),
					resource.TestCheckResourceAttr(rn, "name", alertName),
					resource.TestCheckResourceAttr(rn, "environment", ""),
					resource.TestCheckResourceAttr(rn, "dataset", "transactions"),
					resource.TestCheckResourceAttr(rn, "query", "http.url:http://testservice.com/stats"),
					resource.TestCheckResourceAttr(rn, "aggregate", "p50(transaction.duration)"),
					resource.TestCheckResourceAttr(rn, "time_window", "50"),
					resource.TestCheckResourceAttr(rn, "threshold_type", "0"),
					resource.TestCheckResourceAttr(rn, "resolve_threshold", "100"),
					resource.TestCheckResourceAttr(rn, "projects.#", "1"),
					resource.TestCheckResourceAttrPair(rn, "projects.0", "sentry_project.test", "id"),
					resource.TestCheckResourceAttrSet(rn, "internal_id"),
				),
			},
		},
	})
}

func testAccCheckSentryMetricAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_metric_alert" {
			continue
		}

		org, project, id, err := splitSentryAlertID(rs.Primary.ID)
		if err != nil {
			return err
		}

		ctx := context.Background()
		alert, resp, err := client.MetricAlerts.Get(ctx, org, project, id)
		if err == nil {
			if alert != nil {
				return errors.New("metric alert still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckSentryMetricAlertExists(n string, alert *sentry.MetricAlert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		org, project, alertID, err := splitSentryAlertID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotAlert, _, err := client.MetricAlerts.Get(ctx, org, project, alertID)
		if err != nil {
			return err
		}
		*alert = *gotAlert
		return nil
	}
}

func testAccSentryMetricAlertConfig(teamSlug, projectName, alertName string) string {
	return fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%[1]s"
}

resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[2]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	team         = sentry_team.test.id
	name         = "%[3]s"
	platform     = "go"
}

resource "sentry_metric_alert" "test" {
	organization      = sentry_project.test.organization
	project           = sentry_project.test.id
	name              = "%[4]s"
	dataset           = "transactions"
	query             = "http.url:http://testservice.com/stats"
	aggregate         = "p50(transaction.duration)"
	time_window       = 50.0
	threshold_type    = 0
	resolve_threshold = 100.0

	trigger {
		actions           = []
		alert_threshold   = 1000
		label             = "critical"
		resolve_threshold = 100.0
		threshold_type    = 0
	}

	trigger {
		actions           = []
		alert_threshold   = 500
		label             = "warning"
		resolve_threshold = 100.0
		threshold_type    = 0
	}
}
	`, testOrganization, teamSlug, projectName, alertName)
}
