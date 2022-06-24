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
	var alertCopy sentry.MetricAlert

	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-metric-alert")
	rn := "sentry_metric_alert.test"
	dn := "data.sentry_metric_alert.test"
	rnCopy := "sentry_metric_alert.test_copy"

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
					resource.TestCheckResourceAttrPair(dn, "name", rn, "name"),
					resource.TestCheckResourceAttrPair(dn, "environment", rn, "environment"),
					resource.TestCheckResourceAttrPair(dn, "dataset", rn, "dataset"),
					resource.TestCheckResourceAttrPair(dn, "query", rn, "query"),
					resource.TestCheckResourceAttrPair(dn, "aggregate", rn, "aggregate"),
					resource.TestCheckResourceAttrPair(dn, "time_window", rn, "time_window"),
					resource.TestCheckResourceAttrPair(dn, "threshold_type", rn, "threshold_type"),
					resource.TestCheckResourceAttrPair(dn, "resolve_threshold", rn, "resolve_threshold"),
					resource.TestCheckResourceAttr(dn, "projects.#", "1"),
					resource.TestCheckResourceAttrPair(dn, "projects.0", rn, "projects.0"),
					resource.TestCheckResourceAttrPair(dn, "owners", rn, "owners"),
					resource.TestCheckResourceAttr(dn, "trigger.#", "2"),
					resource.TestCheckResourceAttrPair(dn, "triggers.0", rn, "triggers.0"),
					resource.TestCheckResourceAttr(dn, "trigger.0.action.#", "0"),
					resource.TestCheckResourceAttrPair(dn, "triggers.1", rn, "triggers.1"),
					resource.TestCheckResourceAttr(dn, "trigger.1.action.#", "1"),
					testAccCheckSentryMetricAlertExists(rnCopy, &alertCopy),
					resource.TestCheckResourceAttrPair(rnCopy, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(rnCopy, "project", rn, "project"),
					resource.TestCheckResourceAttrWith(rnCopy, "internal_id", func(v string) error {
						want := sentry.StringValue(alertCopy.ID)
						if v != want {
							return fmt.Errorf("got metric alert ID %s; want %s", v, want)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith(rnCopy, "name", func(v string) error {
						want := sentry.StringValue(alertCopy.Name)
						if v != want {
							return fmt.Errorf("got name ID %s; want %s", v, want)
						}
						return nil
					}),
					resource.TestCheckResourceAttrPair(rnCopy, "environment", rn, "environment"),
					resource.TestCheckResourceAttrPair(rnCopy, "dataset", rn, "dataset"),
					resource.TestCheckResourceAttrPair(rnCopy, "query", rn, "query"),
					resource.TestCheckResourceAttrPair(rnCopy, "aggregate", rn, "aggregate"),
					resource.TestCheckResourceAttrPair(rnCopy, "time_window", rn, "time_window"),
					resource.TestCheckResourceAttrPair(rnCopy, "threshold_type", rn, "threshold_type"),
					resource.TestCheckResourceAttrPair(rnCopy, "resolve_threshold", rn, "resolve_threshold"),
					resource.TestCheckResourceAttr(rnCopy, "projects.#", "1"),
					resource.TestCheckResourceAttrPair(rnCopy, "projects.0", rn, "projects.0"),
					resource.TestCheckResourceAttrPair(rnCopy, "owners", rn, "owners"),
					resource.TestCheckResourceAttr(rnCopy, "trigger.#", "2"),
					resource.TestCheckResourceAttrPair(rnCopy, "trigger.0", rn, "trigger.0"),
					resource.TestCheckResourceAttrPair(rnCopy, "trigger.1", rn, "trigger.1"),
				),
			},
		},
	})
}

func testAccSentryMetricAlertDataSourceConfig(teamSlug, projectName, alertName string) string {
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

resource "sentry_metric_alert" "test" {
	organization      = sentry_project.test.organization
	project           = sentry_project.test.id
	name              = "%[3]s"
	dataset           = "transactions"
	query             = "http.url:http://testservice.com/stats"
	aggregate         = "p50(transaction.duration)"
	time_window       = 50.0
	threshold_type    = 0
	resolve_threshold = 100.0

	trigger {
		alert_threshold   = 1000
		label             = "critical"
		resolve_threshold = 100.0
		threshold_type    = 0
	}

	trigger {
		action {
			type              = "email"
			target_type       = "team"
			target_identifier = sentry_team.test.internal_id
		}

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

	dynamic "trigger" {
		for_each = data.sentry_metric_alert.test.trigger
		content {
			dynamic "action" {
				for_each = trigger.value.action
				content {
					type              = action.value.type
					target_type       = action.value.target_type
					target_identifier = action.value.target_identifier
				}
			}

			alert_threshold   = trigger.value.alert_threshold
			label             = trigger.value.label
			resolve_threshold = trigger.value.resolve_threshold
			threshold_type    = trigger.value.threshold_type
		}
	}
}
	`, teamSlug, projectName, alertName)
}
