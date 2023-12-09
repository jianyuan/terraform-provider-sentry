package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func testAccCheckTeamMemberDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_team_member" {
			continue
		}

		// ctx := context.Background()
		// team, resp, err := acctest.SharedClient.TeamMembers.Get(
		// 	ctx,
		// 	rs.Primary.Attributes["organization"],
		// 	rs.Primary.ID,
		// )
		// if err == nil {
		// 	if team != nil {
		// 		return errors.New("team member still exists")
		// 	}
		// }
		// if resp.StatusCode != 404 {
		// 	return err
		// }
		return nil
	}
	return nil
}

func TestAccTeamMemberResource(t *testing.T) {
	rn := "sentry_team_member.test"
	teamSlug := acctest.RandomWithPrefix("tf-team")
	member1Email := acctest.RandomWithPrefix("tf-member") + "@example.com"
	member2Email := acctest.RandomWithPrefix("tf-member") + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMemberConfig(teamSlug, member1Email, member2Email, "sentry_organization_member.test_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttrPair(rn, "member_id", "sentry_organization_member.test_1", "internal_id"),
					resource.TestCheckResourceAttrPair(rn, "team_slug", "sentry_team.test", "slug"),
				),
			}, {
				Config: testAccTeamMemberConfig(teamSlug, member1Email, member2Email, "sentry_organization_member.test_2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttrPair(rn, "member_id", "sentry_organization_member.test_2", "internal_id"),
					resource.TestCheckResourceAttrPair(rn, "team_slug", "sentry_team.test", "slug"),
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

func testAccTeamMemberConfig(teamName, member1Email, member2Email, memberResourceName string) string {
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
  team_slug    = sentry_team.test.slug
  member_id    = %[4]s.internal_id
}
`, teamName, member1Email, member2Email, memberResourceName)
}
