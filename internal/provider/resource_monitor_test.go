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

func TestAccMonitorResource(t *testing.T) {
	rn := "sentry_monitor.test"
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-monitor")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig(teamName, projectName, monitorName, "0 * * * *", 1, 30),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("schedule_crontab"), knownvalue.StringExact("0 * * * *")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("timezone"), knownvalue.StringExact("UTC")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_runtime"), knownvalue.Int64Exact(30)),
				},
			},
			{
				Config: testAccMonitorResourceConfig(teamName, projectName, monitorName, "*/30 * * * *", 5, 45),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("schedule_crontab"), knownvalue.StringExact("*/30 * * * *")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("timezone"), knownvalue.StringExact("UTC")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("checkin_margin"), knownvalue.Int64Exact(5)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("max_runtime"), knownvalue.Int64Exact(45)),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitorResourceConfig(teamName, projectName, monitorName, schedule string, checkinMargin, maxRuntime int) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.slug
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.slug]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_monitor" "test" {
	organization            = sentry_project.test.organization
	project                 = sentry_project.test.id
	name                    = "%[3]s"
	schedule_crontab        = "%[4]s"
	timezone                = "UTC"
	checkin_margin          = %[5]d
	max_runtime             = %[6]d
	failure_issue_threshold = 1
	recovery_threshold      = 1
}
`, teamName, projectName, monitorName, schedule, checkinMargin, maxRuntime)
}
