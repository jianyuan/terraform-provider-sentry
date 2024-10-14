package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccAllProjectsDataSource(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "data.sentry_all_projects.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllProjectsDataSourceConfig(teamName, projectName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_slugs"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.StringExact(projectName),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"internal_id":  knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`)),
							"slug":         knownvalue.StringExact(projectName),
							"name":         knownvalue.StringExact(projectName),
							"platform":     knownvalue.StringExact("go"),
							"date_created": knownvalue.NotNull(),
							"features":     knownvalue.NotNull(),
							"color":        knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

func testAccAllProjectsDataSourceConfig(teamName string, projectName string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

data "sentry_all_projects" "test" {
	organization = data.sentry_organization.test.slug

	depends_on = [sentry_project.test]
}
`, teamName, projectName)
}
