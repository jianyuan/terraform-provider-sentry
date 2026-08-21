package provider

import (
	"context"
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
					resource.TestCheckResourceAttr(rn, "effective_role", "contributor"),
					resource.TestCheckResourceAttrPair(rn, "member_id", "sentry_organization_member.test_1", "internal_id"),
					resource.TestCheckResourceAttrPair(rn, "team", "sentry_team.test", "slug"),
				),
			},
			{
				Config: testAccTeamMemberConfig(team, member1Email, member2Email, "sentry_organization_member.test_1", "admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "admin"),
					resource.TestCheckResourceAttr(rn, "effective_role", "admin"),
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
				ImportStateVerifyIgnore: []string{
					"role",
				},
			},
		},
	})
}

// TestAccTeamMemberResource_teamDeletedOutOfBand reproduces the case where the
// team backing a membership is removed from Sentry outside of Terraform (e.g. it
// was deleted or re-slugged). On destroy, the DELETE against the stale team slug
// returns 404; the resource must treat that as a successful delete instead of
// erroring. Without the fix, teardown fails with "Unable to delete, got error:
// ... 404 The requested resource does not exist".
func TestAccTeamMemberResource_teamDeletedOutOfBand(t *testing.T) {
	team := acctest.RandomWithPrefix("tf-team")
	memberEmail := acctest.RandomWithPrefix("tf-member") + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMemberConfig_minimumPriority(team, memberEmail, "member", "contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sentry_team_member.test", "organization", acctest.TestOrganization),
				),
			},
			{
				// Delete the team out-of-band so the membership's team slug no
				// longer exists, then destroy. The membership's DELETE now 404s
				// and must be treated as already-deleted rather than an error.
				Config:  testAccTeamMemberConfig_minimumPriority(team, memberEmail, "member", "contributor"),
				Destroy: true,
				PreConfig: func() {
					_, err := acctest.SharedClient.Teams.Delete(context.Background(), acctest.TestOrganization, team)
					if err != nil {
						t.Fatalf("failed to delete team out-of-band: %s", err)
					}
				},
			},
		},
	})
}

func testAccTeamMemberConfig(teamName, member1Email, member2Email, memberResourceName, teamRole string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.slug
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_organization_member" "test_1" {
	organization = data.sentry_organization.test.slug
	email        = "%[2]s"
	role         = "member"
}

resource "sentry_organization_member" "test_2" {
	organization = data.sentry_organization.test.slug
	email        = "%[3]s"
	role         = "member"
}

resource "sentry_team_member" "test" {
	organization = data.sentry_organization.test.slug
	team         = sentry_team.test.slug
	member_id    = %[4]s.internal_id
	role         = "%[5]s"
}
`, teamName, member1Email, member2Email, memberResourceName, teamRole)
}

func TestAccTeamMemberResource_minimumPriority(t *testing.T) {
	rn := "sentry_team_member.test"
	team := acctest.RandomWithPrefix("tf-team")
	memberEmail := acctest.RandomWithPrefix("tf-member") + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamMemberConfig_minimumPriority(team, memberEmail, "member", "contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "contributor"),
					resource.TestCheckResourceAttr(rn, "effective_role", "contributor"),
				),
			},
			{
				Config: testAccTeamMemberConfig_minimumPriority(team, memberEmail, "member", "admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "admin"),
					resource.TestCheckResourceAttr(rn, "effective_role", "admin"),
				),
			},
			{
				Config: testAccTeamMemberConfig_minimumPriority(team, memberEmail, "owner", "contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "contributor"),
					resource.TestCheckResourceAttr(rn, "effective_role", "admin"),
				),
			},
			{
				Config: testAccTeamMemberConfig_minimumPriority(team, memberEmail, "owner", "admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "role", "admin"),
					resource.TestCheckResourceAttr(rn, "effective_role", "admin"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"role",
				},
			},
		},
	})
}

func testAccTeamMemberConfig_minimumPriority(teamName, memberEmail, memberRole, teamRole string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.slug
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_organization_member" "test" {
	organization = data.sentry_organization.test.slug
	email        = "%[2]s"
	role         = "%[3]s"
}

resource "sentry_team_member" "test" {
	organization = data.sentry_organization.test.slug
	team         = sentry_team.test.slug
	member_id    = sentry_organization_member.test.internal_id
	role         = "%[4]s"
}
`, teamName, memberEmail, memberRole, teamRole)
}
