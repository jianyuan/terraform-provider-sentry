package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSentryTeam_basic(t *testing.T) {
	var team Team

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
						Name:         "Test team",
						Organization: testOrganization,
						SlugPresent:  true,
					}),
				),
			},
			{
				Config: testAccSentryTeamUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryTeamExists("sentry_team.test_team", &team),
					testAccCheckSentryTeamAttributes(&team, &testAccSentryTeamExpectedAttributes{
						Name:         "Test team changed",
						Organization: testOrganization,
						Slug:         newTeamSlug,
					}),
				),
			},
		},
	})
}

func testAccCheckSentryTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_team" {
			continue
		}

		team, resp, err := client.GetTeam(testOrganization, rs.Primary.ID)
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

func testAccCheckSentryTeamExists(n string, team *Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No team ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		sentryTeam, _, err := client.GetTeam(testOrganization, rs.Primary.ID)
		if err != nil {
			return err
		}
		*team = *sentryTeam
		return nil
	}
}

type testAccSentryTeamExpectedAttributes struct {
	Name         string
	Organization string

	SlugPresent bool
	Slug        string
}

func (attrs *testAccSentryTeamExpectedAttributes) internalValidate() error {
	if attrs == nil {
		return errors.New("attribute is nil")
	}

	if attrs.SlugPresent && attrs.Slug != "" {
		return errors.New("cannot provide both SlugPresent and Slug")
	}

	return nil
}

func testAccCheckSentryTeamAttributes(team *Team, want *testAccSentryTeamExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := want.internalValidate(); err != nil {
			return err
		}

		if team.Name != want.Name {
			return fmt.Errorf("got team %q; want %q", team.Name, want.Name)
		}

		if team.Organization.Slug != want.Organization {
			return fmt.Errorf("got organization %q; want %q", team.Organization.Slug, want.Organization)
		}

		if want.SlugPresent && team.Slug == "" {
			return errors.New("got empty slug; want non-empty slug")
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
