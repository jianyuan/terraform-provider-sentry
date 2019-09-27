package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryTeam_basic(t *testing.T) {
	var team sentry.Team

	random := acctest.RandInt()
	newTeamSlug := fmt.Sprintf("test-team-changed-%d", random)

	testAccSentryTeamUpdateConfig := fmt.Sprintf(`
    resource "sentry_team" "test_team" {
      organization = "%s"
      name = "Test team changed"
      slug = "%s"
    }
	`, testOrganization, newTeamSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryTeamConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryTeamExists("sentry_team.test_team", &team),
					testAccCheckSentryTeamAttributes(&team, &testAccSentryTeamExpectedAttributes{
						Name:        "Test team",
						SlugPresent: true,
					}),
				),
			},
			{
				Config: testAccSentryTeamUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryTeamExists("sentry_team.test_team", &team),
					testAccCheckSentryTeamAttributes(&team, &testAccSentryTeamExpectedAttributes{
						Name: "Test team changed",
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
	Name string

	SlugPresent bool
	Slug        string
}

func testAccCheckSentryTeamAttributes(team *sentry.Team, want *testAccSentryTeamExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if team.Name != want.Name {
			return fmt.Errorf("got team %q; want %q", team.Name, want.Name)
		}

		if want.SlugPresent && team.Slug == "" {
			return errors.New("got empty slug; want non-empty slug")
		}

		if want.Slug != "" && team.Slug != want.Slug {
			return fmt.Errorf("got slug %q; want %q", team.Slug, want.Slug)
		}

		return nil
	}
}

var testAccSentryTeamConfig = fmt.Sprintf(`
  resource "sentry_team" "test_team" {
    organization = "%s"
    name = "Test team"
  }
`, testOrganization)
