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

func TestAccMonitorDataSource_validation(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: testAccMonitorDataSourceConfig(teamName, projectName, `
					id = "1"
					first = true
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute "first" cannot be specified when "id" is specified`),
			},
			{
				PlanOnly: true,
				Config: testAccMonitorDataSourceConfig(teamName, projectName, `
					id = "1"
					project = "2"
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute "project" cannot be specified when "id" is specified`),
			},
			{
				Config:      testAccMonitorDataSourceConfig(teamName, projectName, ``),
				ExpectError: acctest.ExpectLiteralError("Multiple monitors found, please narrow down the search by setting the `type` attribute, and/or set the `first` attribute to `true`."),
			},
		},
	})
}

func TestAccMonitorDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "data.sentry_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorDataSourceConfig(teamName, projectName, `
					project = sentry_project.test.id
					type = "issue_stream"
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("Issue Stream")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
				),
			},
		},
	})
}

func testAccMonitorDataSourceConfig(teamName, projectName, extras string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	}) + fmt.Sprintf(`
		data "sentry_monitor" "test" {
			organization = sentry_project.test.organization
			%s
		}
	`, extras)
}
