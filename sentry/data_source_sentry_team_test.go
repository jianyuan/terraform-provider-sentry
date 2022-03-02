package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSentryTeamDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryTeamataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sentry_team.test", "name", "Test team"),
					resource.TestCheckResourceAttr("data.sentry_team.test", "slug", "test-team"),
					resource.TestCheckResourceAttrSet("data.sentry_team.test", "team_id"),
				),
			},
		},
	})
}

// Testing first parameter
var testAccSentryTeamataSourceConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
  organization = "%s"
  name = "Test team"
	slug = "test-team"
}

data "sentry_team" "test_key" {
  organization = "%s"
	slug = "test-team"
}
`, testOrganization, testOrganization)
