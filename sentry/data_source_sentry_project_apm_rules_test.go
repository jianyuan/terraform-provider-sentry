package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSentyAPMRulesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryAPMRulesDataSourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

var testAccSentryAPMRulesDataSourceConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
	organization = "%s"
	name = "Test team"
}

resource "sentry_project" "test_project" {
	organization = "%s"
	team = "${sentry_team.test_team.id}"
	name = "Test project"
}

//create single apm rule to test that it does pull back apm rules

data "sentry_apm_rules" "test" {
    organization = "%s"
    project = "${sentry_project.test_project.id}"
}
`, testOrganization, testOrganization, testOrganization)
