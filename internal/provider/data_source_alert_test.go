package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccAlertDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-monitor")
	alertName := acctest.RandomWithPrefix("tf-alert")
	opsgenieTeamName := acctest.RandomWithPrefix("tf-opsgenie")
	rn := "data.sentry_alert.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("monitor_ids"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency_minutes"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("triggers_json"), acctest.StringJson()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_filters_json"), acctest.StringJson()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertDataSourceConfig(teamName, projectName, monitorName, alertName, opsgenieTeamName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alertName)),
				),
			},
			{
				Config: testAccAlertDataSourceConfig(teamName, projectName, monitorName, alertName+"-updated", opsgenieTeamName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alertName+"-updated")),
				),
			},
		},
	})
}

func testAccAlertDataSourceConfig(teamName, projectName, monitorName, name, opsgenieTeamName string) string {
	return testAccAlertResourceConfig(teamName, projectName, monitorName, name, opsgenieTeamName) + `
		data "sentry_alert" "test" {
			organization = sentry_alert.test.organization
			id           = sentry_alert.test.id
		}
	`
}
