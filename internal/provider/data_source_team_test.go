package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mzglinski/terraform-provider-sentry/internal/acctest"
)

func TestAccTeamDataSource(t *testing.T) {
	dsn := "data.sentry_team.test"

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("slug"), knownvalue.StringExact(acctest.TestTeam.Slug)),
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("internal_id"), knownvalue.StringExact(acctest.TestTeam.Id)),
		statecheck.ExpectKnownValue(dsn, tfjsonpath.New("name"), knownvalue.StringExact(acctest.TestTeam.Name)),
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
				Config:            testAccTeamDataSourceConfig(),
				ConfigStateChecks: configStateChecks,
			},
		},
	})
}

func testAccTeamDataSourceConfig() string {
	return fmt.Sprintf(`
data "sentry_team" "test" {
	organization = "%s"
	slug         = "%s"
}
`, acctest.TestOrganization, acctest.TestTeam.Slug)
}
