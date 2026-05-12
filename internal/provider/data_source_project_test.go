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

func TestAccProjectDataSource(t *testing.T) {
	rn := "data.sentry_project.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact(acctest.TestProjectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.StringExact(acctest.TestProjectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`))),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("is_public"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceConfig(),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(acctest.TestProjectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.StringExact("go")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("date_created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("subject_template"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("color"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.NotNull()),
				),
			},
		},
	})
}

func testAccProjectDataSourceConfig() string {
	return fmt.Sprintf(`
data "sentry_project" "test" {
	organization = "%s"
	slug         = "%s"
}
`, acctest.TestOrganization, acctest.TestProject.Slug)
}
