package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryTeamDataSource_basic(t *testing.T) {
	var team sentry.Team

	teamSlug := acctest.RandomWithPrefix("tf-team")
	rn := "sentry_team.test"
	dn := "data.sentry_team.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryTeamDataSourceConfig(teamSlug),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryTeamExists(rn, &team),
					resource.TestCheckResourceAttrPair(dn, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(dn, "slug", rn, "slug"),
					resource.TestCheckResourceAttrPair(dn, "internal_id", rn, "internal_id"),
					resource.TestCheckResourceAttrPair(dn, "name", rn, "name"),
					resource.TestCheckResourceAttrPair(dn, "has_access", rn, "has_access"),
					resource.TestCheckResourceAttrPair(dn, "is_pending", rn, "is_pending"),
					resource.TestCheckResourceAttrPair(dn, "is_member", rn, "is_member"),
				),
			},
		},
	})
}

func testAccSentryTeamDataSourceConfig(teamSlug string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

data "sentry_team" "test" {
	organization = sentry_team.test.organization
	slug         = sentry_team.test.id
}
	`, teamSlug)
}
