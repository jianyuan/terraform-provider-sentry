package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccAllClientKeysDataSource(t *testing.T) {
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	rn := "data.sentry_all_keys.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllClientKeysDataSourceConfig(team, project),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("keys"), knownvalue.ListExact([]knownvalue.Check{
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

func testAccAllClientKeysDataSourceConfig(teamName, projectName string) string {
	return testAccProjectResourceConfig(teamName, projectName) + `
data "sentry_all_keys" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
}
`
}
