package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryTeam_basic(t *testing.T) {
	var team sentry.Team

	teamName := acctest.RandomWithPrefix("tf-team")
	rn := "sentry_team.test_team"

	check := func(teamName string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryTeamExists(rn, &team),
			resource.TestCheckResourceAttrPair(rn, "organization", "data.sentry_organization.test_organization", "id"),
			resource.TestCheckResourceAttr(rn, "name", teamName),
			resource.TestCheckResourceAttr(rn, "slug", teamName),
			resource.TestCheckResourceAttrWith(rn, "internal_id", func(v string) error {
				want := sentry.StringValue(team.ID)
				if v != want {
					return fmt.Errorf("got team ID %s; want %s", v, want)
				}
				return nil
			}),
			resource.TestCheckResourceAttrPair(rn, "internal_id", rn, "team_id"),
			resource.TestCheckResourceAttrSet(rn, "has_access"),
			resource.TestCheckResourceAttrSet(rn, "is_pending"),
			resource.TestCheckResourceAttrSet(rn, "is_member"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryTeamConfig(teamName),
				Check:  check(teamName),
			},
			{
				Config: testAccSentryTeamConfig(teamName + "-renamed"),
				Check:  check(teamName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccSentryTeamImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_team" {
			continue
		}

		ctx := context.Background()
		team, resp, err := client.Teams.Get(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err == nil {
			if team != nil {
				return errors.New("Team still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryTeamExists(n string, team *sentry.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		org := rs.Primary.Attributes["organization"]
		teamSlug := rs.Primary.ID
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotTeam, _, err := client.Teams.Get(ctx, org, teamSlug)
		if err != nil {
			return err
		}
		*team = *gotTeam
		return nil
	}
}

func testAccSentryTeamImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		org := rs.Primary.Attributes["organization"]
		teamSlug := rs.Primary.ID
		return buildTwoPartID(org, teamSlug), nil
	}
}

func testAccSentryTeamConfig(teamName string) string {
	return fmt.Sprintf(`
data "sentry_organization" "test_organization" {
	slug = "%[1]s"
}

resource "sentry_team" "test_team" {
	organization = data.sentry_organization.test_organization.id
	name         = "%[2]s"
	slug         = "%[2]s"
}
	`, testOrganization, teamName)
}
