package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccAllOrganizationMembersDataSource(t *testing.T) {
	rn := "data.sentry_all_organization_members.test"
	email := acctest.RandomWithPrefix("tf-member") + "@example.com"
	role := "member"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllOrganizationMembersConfig(email, role),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("members"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":    knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`)),
							"email": knownvalue.StringExact(email),
							"role":  knownvalue.StringExact(role),
						}),
					})),
				},
			},
		},
	})
}

func testAccAllOrganizationMembersConfig(email string, role string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_organization_member" "test" {
	organization = data.sentry_organization.test.id
	email        = "%[1]s"
	role         = "%[2]s"
}

data "sentry_all_organization_members" "test" {
	organization = data.sentry_organization.test.id

	depends_on = [sentry_organization_member.test]
}
`, email, role)
}
