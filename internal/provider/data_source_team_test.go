package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccTeamDataSource(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	rn := "sentry_team.test"
	dsn := "data.sentry_team.test"

	configStateChecks := []statecheck.StateCheck{
		statecheck.CompareValuePairs(dsn, tfjsonpath.New("organization"), rn, tfjsonpath.New("organization"), compare.ValuesSame()),
		statecheck.CompareValuePairs(dsn, tfjsonpath.New("slug"), rn, tfjsonpath.New("slug"), compare.ValuesSame()),
		statecheck.CompareValuePairs(dsn, tfjsonpath.New("internal_id"), rn, tfjsonpath.New("internal_id"), compare.ValuesSame()),
		statecheck.CompareValuePairs(dsn, tfjsonpath.New("name"), rn, tfjsonpath.New("name"), compare.ValuesSame()),
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("id"), knownvalue.NotNull()),         // Deprecated
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("has_access"), knownvalue.NotNull()), // Deprecated
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("is_pending"), knownvalue.NotNull()), // Deprecated
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("is_member"), knownvalue.NotNull()),  // Deprecated
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccTeamDataSourceConfig(teamName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: configStateChecks,
			},
		},
	})
}

func testAccTeamDataSourceConfig(teamName string) string {
	return testAccTeamResourceConfig(teamName) + `
data "sentry_team" "test" {
	organization = sentry_team.test.organization
	slug         = sentry_team.test.slug
}
`
}
