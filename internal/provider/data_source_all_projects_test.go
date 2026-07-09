package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mzglinski/terraform-provider-sentry/internal/acctest"
)

func TestAccAllProjectsDataSource(t *testing.T) {
	rn := "data.sentry_all_projects.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllProjectsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_slugs"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.StringExact(acctest.TestProject.Slug),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"internal_id":  knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`)),
							"slug":         knownvalue.StringExact(acctest.TestProject.Slug),
							"name":         knownvalue.StringExact(acctest.TestProject.Name),
							"platform":     knownvalue.StringExact("go"),
							"date_created": knownvalue.NotNull(),
							"features":     knownvalue.NotNull(),
							"color":        knownvalue.NotNull(),
							"teams":        knownvalue.NotNull(),
						}),
					})),
				},
			},
		},
	})
}

func testAccAllProjectsDataSourceConfig() string {
	return fmt.Sprintf(`
data "sentry_all_projects" "test" {
	organization = "%s"
}
`, acctest.TestOrganization)
}
