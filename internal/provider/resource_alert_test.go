package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccAlertResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-monitor")
	alertName := acctest.RandomWithPrefix("tf-alert")
	rn := "sentry_alert.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertResourceConfig(teamName, projectName, monitorName, alertName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alertName)),
				),
			},
		},
	})
}

func testAccAlertResourceConfig(teamName, projectName, monitorName, name string) string {
	return testAccMetricMonitorResourceConfig(teamName, projectName, monitorName) + fmt.Sprintf(`
		resource "sentry_alert" "test" {
			organization = data.sentry_organization.test.slug
			name         = "%[1]s"

			frequency_minutes = 1440
			environment       = "production"
			monitor_ids       = [sentry_metric_monitor.test.id]

			trigger_conditions = [
				"first_seen_event",
				"issue_resolved_trigger",
				"reappeared_event",
				"regression_event",
			]
		}
	`, name)
}
