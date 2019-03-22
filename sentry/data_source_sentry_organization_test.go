package sentry

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSentryOrganizationDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryOrganizationDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sentry_organization.test", "name"),
					resource.TestMatchResourceAttr("data.sentry_organization.test", "internal_id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr("data.sentry_organization.test", "slug", testOrganization),
				),
			},
		},
	})
}

var testAccSentryOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%s"
}
`, testOrganization)
