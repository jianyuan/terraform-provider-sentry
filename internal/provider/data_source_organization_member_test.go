package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccOrganizationMemberDataSource(t *testing.T) {
	rn := "sentry_organization_member.test"
	dn := "data.sentry_organization_member.test"
	email := acctest.RandomWithPrefix("tf-member") + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationMemberConfig(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dn, "id", rn, "internal_id"),
					resource.TestCheckResourceAttr(dn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(dn, "email", email),
					resource.TestCheckResourceAttr(dn, "role", "member"),
				),
			},
		},
	})
}

func testAccOrganizationMemberConfig(email string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_organization_member" "test" {
	organization = data.sentry_organization.test.id
	email        = "%[1]s"
	role         = "member"
}

data "sentry_organization_member" "test" {
	organization = data.sentry_organization.test.id
	email        = sentry_organization_member.test.email
}
`, email)
}
