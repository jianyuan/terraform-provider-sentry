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

func TestAccSentryProject_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	check := func(projectName string) resource.TestCheckFunc {
		var project sentry.Project

		return resource.ComposeTestCheckFunc(
			testAccCheckSentryProjectExists(rn, &project),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "team", teamName),
			resource.TestCheckResourceAttr(rn, "name", projectName),
			resource.TestCheckResourceAttrSet(rn, "slug"),
			resource.TestCheckResourceAttr(rn, "platform", "go"),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrWith(rn, "internal_id", func(v string) error {
				want := project.ID
				if v != want {
					return fmt.Errorf("got project ID %s; want %s", v, want)
				}
				return nil
			}),
			resource.TestCheckResourceAttrPair(rn, "project_id", rn, "internal_id"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectConfig(teamName, projectName),
				Check:  check(projectName),
			},
			{
				Config: testAccSentryProjectConfig(teamName, projectName+"-renamed"),
				Check:  check(projectName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccSentryProjectImportStateIdFunc(rn),
				ImportStateVerify: true,
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

		ctx := context.Background()
		proj, resp, err := client.Projects.Get(ctx, testOrganization, rs.Primary.ID)
		if err == nil {
			if proj != nil {
				return errors.New("project still exists")
			}
		}
		if resp.StatusCode != 403 && resp.StatusCode != 404 {
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
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotProj, _, err := client.Projects.Get(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		}
		*proj = *gotProj
		return nil
	}
}

func testAccSentryProjectImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		org := rs.Primary.Attributes["organization"]
		projectSlug := rs.Primary.ID
		return buildTwoPartID(org, projectSlug), nil
	}
}

func testAccSentryProjectConfig(teamName, projectName string) string {
	return fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%[1]s"
}

resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[2]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	team         = sentry_team.test.id
	name         = "%[3]s"
	platform     = "go"
}
	`, testOrganization, teamName, projectName)
}
