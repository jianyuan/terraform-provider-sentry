package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccOrganizationIntegrationDataSource(t *testing.T) {
	dn := "data.sentry_organization_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationIntegrationDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dn, "id"),
					resource.TestCheckResourceAttrPair(dn, "internal_id", dn, "id"),
					resource.TestCheckResourceAttr(dn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(dn, "provider_key", "github"),
					resource.TestCheckResourceAttr(dn, "name", "jianyuan"),
				),
			},
		},
	})
}

var testAccOrganizationIntegrationDataSourceConfig = testAccOrganizationDataSourceConfig + `
data "sentry_organization_integration" "test" {
	organization = data.sentry_organization.test.slug
	provider_key = "github"
	name         = "jianyuan"
}
`
