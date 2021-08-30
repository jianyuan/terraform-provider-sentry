package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryTeam_basic(t *testing.T) {
	var team sentry.Team

	random := acctest.RandString(4)
	teamSlug := fmt.Sprintf("test-%s", random)
	newTeamSlug := fmt.Sprintf("test-%s-changed", random)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTestAccSentryTeamConfig(teamSlug),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryTeamExists("sentry_team.test_team", &team),
					testAccCheckSentryTeamAttributes(&team, &testAccSentryTeamExpectedAttributes{
						Slug: teamSlug,
					}),
				),
			},
			{
				Config: buildTestAccSentryTeamConfig(newTeamSlug),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryTeamExists("sentry_team.test_team", &team),
					testAccCheckSentryTeamAttributes(&team, &testAccSentryTeamExpectedAttributes{
						Slug: newTeamSlug,
					}),
				),
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

		team, resp, err := client.Teams.Get(
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No team ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		sentryTeam, _, err := client.Teams.Get(
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		}
		*team = *sentryTeam
		return nil
	}
}

type testAccSentryTeamExpectedAttributes struct {
	Slug string
}

func testAccCheckSentryTeamAttributes(team *sentry.Team, want *testAccSentryTeamExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if want.Slug != "" && team.Slug != want.Slug {
			return fmt.Errorf("got slug %q; want %q", team.Slug, want.Slug)
		}

		return nil
	}
}

func buildTestAccSentryTeamConfig(slug string) string {
	return fmt.Sprintf(`
		resource "sentry_team" "test_team" {
			organization = "%s"
			slug = "%s"
		}
	`, testOrganization, slug)
}
