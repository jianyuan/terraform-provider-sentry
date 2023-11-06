package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccTeamDataSource(t *testing.T) {
	ctx := context.Background()

	teamSlug := acctest.RandomWithPrefix("tf-team")
	rn := "sentry_team.test"
	dsn := "data.sentry_team.test"

	var v sentry.Team

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamDataSourceConfig(teamSlug),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamExists(ctx, rn, &v),
					resource.TestCheckResourceAttrPair(dsn, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(dsn, "slug", rn, "slug"),
					resource.TestCheckResourceAttrPair(dsn, "internal_id", rn, "internal_id"),
					resource.TestCheckResourceAttrPair(dsn, "name", rn, "name"),
					resource.TestCheckResourceAttrPair(dsn, "has_access", rn, "has_access"),
					resource.TestCheckResourceAttrPair(dsn, "is_pending", rn, "is_pending"),
					resource.TestCheckResourceAttrPair(dsn, "is_member", rn, "is_member"),
				),
			},
		},
	})
}

func testAccTeamDataSourceConfig(teamName string) string {
	return testAccTeamConfig(teamName) + `
data "sentry_team" "test" {
  organization = sentry_team.test.organization
  slug         = sentry_team.test.id
}
`
}
