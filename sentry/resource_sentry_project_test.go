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

func TestAccSentryProject_basic(t *testing.T) {
	var project sentry.Project

	random := acctest.RandInt()
	newProjectSlug := fmt.Sprintf("test-project-%d", random)

	testAccSentryProjectUpdateConfig := fmt.Sprintf(`
	  resource "sentry_team" "test_team" {
	    organization = "%s"
	    name = "Test team"
	  }

	  resource "sentry_project" "test_project" {
	    organization = "%s"
	    team = "${sentry_team.test_team.id}"
	    name = "Test project changed"
	    slug = "%s"
	  }
	`, testOrganization, testOrganization, newProjectSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryProjectExists("sentry_project.test_project", &project),
					testAccCheckSentryProjectAttributes(&project, &testAccSentryProjectExpectedAttributes{
						Name:         "Test project",
						Organization: testOrganization,
						Team:         "Test team",
						SlugPresent:  true,
					}),
				),
			},
			{
				Config: testAccSentryProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryProjectExists("sentry_project.test_project", &project),
					testAccCheckSentryProjectAttributes(&project, &testAccSentryProjectExpectedAttributes{
						Name:         "Test project changed",
						Organization: testOrganization,
						Team:         "Test team",
						Slug:         newProjectSlug,
					}),
				),
			},
		},
	})
}

func testAccCheckSentryProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_project" {
			continue
		}

		proj, resp, err := client.Projects.Get(testOrganization, rs.Primary.ID)
		if err == nil {
			if proj != nil {
				return errors.New("Project still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryProjectExists(n string, proj *sentry.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		sentryProj, _, err := client.Projects.Get(
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		}
		*proj = *sentryProj
		return nil
	}
}

type testAccSentryProjectExpectedAttributes struct {
	Name         string
	Organization string
	Team         string

	SlugPresent bool
	Slug        string
}

func testAccCheckSentryProjectAttributes(proj *sentry.Project, want *testAccSentryProjectExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if proj.Name != want.Name {
			return fmt.Errorf("got proj %q; want %q", proj.Name, want.Name)
		}

		if proj.Organization.Slug != want.Organization {
			return fmt.Errorf("got organization %q; want %q", proj.Organization.Slug, want.Organization)
		}

		if proj.Team.Name != want.Team {
			return fmt.Errorf("got team %q; want %q", proj.Team.Name, want.Team)
		}

		if want.SlugPresent && proj.Slug == "" {
			return errors.New("got empty slug; want non-empty slug")
		}

		if want.Slug != "" && proj.Slug != want.Slug {
			return fmt.Errorf("got slug %q; want %q", proj.Slug, want.Slug)
		}

		return nil
	}
}

var testAccSentryProjectConfig = fmt.Sprintf(`
  resource "sentry_team" "test_team" {
    organization = "%s"
    name = "Test team"
  }

  resource "sentry_project" "test_project" {
    organization = "%s"
    team = "${sentry_team.test_team.id}"
    name = "Test project"
  }
`, testOrganization, testOrganization)
