package sentry

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSentyAlertRulesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryAlertRulesDataSourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

var testAccSentryAlertRulesDataSourceConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
	organization = "%s"
	name = "Test team"
}

resource "sentry_project" "test_project" {
	organization = "%s"
	team = "${sentry_team.test_team.id}"
	name = "Test project"
}

data "sentry_alert_rules" "test" {
    organization = "%s"
    project = "${sentry_project.test_project.id}"
}
`, testOrganization, testOrganization, testOrganization)
