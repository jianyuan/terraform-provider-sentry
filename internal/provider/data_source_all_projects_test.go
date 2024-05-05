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

func TestAccAllProjectsDataSource(t *testing.T) {
	dn := "data.sentry_all_projects.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllProjectsDataSourceConfig(team, project),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dn, tfjsonpath.New("projects"), knownvalue.ListPartial(map[int]knownvalue.Check{
						0: knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":           knownvalue.NotNull(),
							"slug":         knownvalue.NotNull(),
							"name":         knownvalue.NotNull(),
							"platform":     knownvalue.NotNull(),
							"date_created": knownvalue.NotNull(),
							"features":     knownvalue.NotNull(),
							"color":        knownvalue.NotNull(),
							"status":       knownvalue.NotNull(),
							"organization": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":   knownvalue.NotNull(),
								"slug": knownvalue.NotNull(),
								"name": knownvalue.NotNull(),
							}),
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
}
`, teamName, projectName)
}
