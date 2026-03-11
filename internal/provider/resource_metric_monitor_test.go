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

func TestAccMetricMonitorResource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-metric-monitor")
	rn := "sentry_metric_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricMonitorResourceConfig(teamName, projectName, monitorName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(monitorName)),
				),
			},
		},
	})
}

func testAccMetricMonitorResourceConfig(teamName, projectName, name string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		resource "sentry_metric_monitor" "test" {
			organization = data.sentry_organization.test.slug
			project      = sentry_project.test.internal_id
			name         = "%[1]s"

			condition_group = {}

			default_assignee = {
				team_id = sentry_team.test.internal_id
			}
		}
	`, name)
}
