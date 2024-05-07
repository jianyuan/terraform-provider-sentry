package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccProjectDataSource_MigrateFromPluginSDK(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "data.sentry_project.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`))),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("is_public"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.12.3",
					},
				},
				Config:            testAccProjectDataSourceConfig(teamName, projectName),
				ConfigStateChecks: checks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccProjectDataSourceConfig(teamName, projectName),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("date_created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("color"), knownvalue.NotNull()),
				),
			},
		},
	})
}

func testAccProjectDataSourceConfig(teamName, projectName string) string {
	return testAccProjectResourceConfig(teamName, projectName) + `
data "sentry_project" "test" {
	organization = sentry_project.test.organization
	slug         = sentry_project.test.id
}
`
}
