package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryMetricAlertDataSource_basic(t *testing.T) {
	var alert sentry.MetricAlert

	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-metric-alert")
	rn := "sentry_metric_alert.test"
	dn := "data.sentry_metric_alert.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryMetricAlertDataSourceConfig(teamSlug, projectName, alertName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryMetricAlertExists(rn, &alert),
					resource.TestCheckResourceAttr(dn, "organization", testOrganization),
					resource.TestCheckResourceAttr(dn, "project", projectName),
					resource.TestCheckResourceAttrPair(dn, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(dn, "project", rn, "project"),
					resource.TestCheckResourceAttrPair(dn, "internal_id", rn, "internal_id"),
					// TODO: Other fields
				),
			},
		},
	})
}

func testAccSentryMetricAlertDataSourceConfig(teamSlug, projectName, alertName string) string {
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

	triggers {
		actions           = []
		alert_threshold   = 1000
		label             = "critical"
		resolve_threshold = 100.0
		threshold_type    = 0
	}

	triggers {
		actions           = []
		alert_threshold   = 500
		label             = "warning"
		resolve_threshold = 100.0
		threshold_type    = 0
	}
}

data "sentry_metric_alert" "test" {
	organization = sentry_metric_alert.test.organization
	project      = sentry_metric_alert.test.project
	internal_id  = sentry_metric_alert.test.internal_id
}

resource "sentry_metric_alert" "test_copy" {
	organization      = data.sentry_metric_alert.test.organization
	project           = data.sentry_metric_alert.test.project
	name              = "${data.sentry_metric_alert.test.name}-copy"
	dataset           = data.sentry_metric_alert.test.dataset
	query             = data.sentry_metric_alert.test.query
	aggregate         = data.sentry_metric_alert.test.aggregate
	time_window       = data.sentry_metric_alert.test.time_window
	threshold_type    = data.sentry_metric_alert.test.threshold_type
	resolve_threshold = data.sentry_metric_alert.test.resolve_threshold

	dynamic "triggers" {
		for_each = data.sentry_metric_alert.test.triggers
		content {
			actions           = triggers.value.actions
			alert_threshold   = triggers.value.alert_threshold
			label             = triggers.value.label
			resolve_threshold = triggers.value.resolve_threshold
			threshold_type    = triggers.value.threshold_type
		}
	}
}
	`, testOrganization, teamSlug, projectName, alertName)
}
