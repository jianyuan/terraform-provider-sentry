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

func TestAccClientKeysDataSource(t *testing.T) {
	dn := "data.sentry_keys.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeysDataSourceConfig(team, project),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(dn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
					statecheck.ExpectKnownValue(dn, tfjsonpath.New("keys"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.MapExact(map[string]knownvalue.Check{
							"id":                knownvalue.NotNull(),
							"organization":      knownvalue.StringExact(acctest.TestOrganization),
							"project":           knownvalue.StringExact(project),
							"project_id":        knownvalue.NotNull(),
							"name":              knownvalue.StringExact("Default"),
							"public":            knownvalue.NotNull(),
							"secret":            knownvalue.NotNull(),
							"rate_limit_window": knownvalue.Null(),
							"rate_limit_count":  knownvalue.Null(),
							"dsn_public":        knownvalue.NotNull(),
							"dsn_secret":        knownvalue.NotNull(),
							"dsn_csp":           knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

func testAccClientKeysDataSourceConfig(teamName string, projectName string) string {
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

data "sentry_keys" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
}
`, teamName, projectName)
}
