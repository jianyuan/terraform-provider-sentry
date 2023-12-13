package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccTeamMemberResource(t *testing.T) {
	rn := "sentry_team_member.test"
	team := acctest.RandomWithPrefix("tf-team")
	member1Email := acctest.RandomWithPrefix("tf-member") + "@example.com"
	member2Email := acctest.RandomWithPrefix("tf-member") + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMemberConfig(team, member1Email, member2Email, "sentry_organization_member.test_1", "contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "contributor"),
					resource.TestCheckResourceAttrPair(rn, "member_id", "sentry_organization_member.test_1", "internal_id"),
					resource.TestCheckResourceAttrPair(rn, "team", "sentry_team.test", "slug"),
				),
			},
			{
				Config: testAccTeamMemberConfig(team, member1Email, member2Email, "sentry_organization_member.test_1", "admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "admin"),
					resource.TestCheckResourceAttrPair(rn, "member_id", "sentry_organization_member.test_1", "internal_id"),
					resource.TestCheckResourceAttrPair(rn, "team", "sentry_team.test", "slug"),
				),
			},
			{
				Config: testAccTeamMemberConfig(team, member1Email, member2Email, "sentry_organization_member.test_2", "contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "contributor"),
					resource.TestCheckResourceAttrPair(rn, "member_id", "sentry_organization_member.test_2", "internal_id"),
					resource.TestCheckResourceAttrPair(rn, "team", "sentry_team.test", "slug"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTeamMemberConfig(teamName, member1Email, member2Email, memberResourceName, memberRole string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_organization_member" "test_1" {
	organization = data.sentry_organization.test.id
	email        = "%[2]s"
	role         = "member"
}

resource "sentry_organization_member" "test_2" {
	organization = data.sentry_organization.test.id
	email        = "%[3]s"
	role         = "member"
}

resource "sentry_team_member" "test" {
	organization = data.sentry_organization.test.id
	team         = sentry_team.test.id
	member_id    = %[4]s.internal_id
	role         = "%[5]s"
}
`, teamName, member1Email, member2Email, memberResourceName, memberRole)
}
