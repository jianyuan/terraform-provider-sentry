package sentry

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/canva/go-sentry/sentry"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryProject_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	check := func(projectName string) resource.TestCheckFunc {
		var projectID string

		return resource.ComposeTestCheckFunc(
			testAccCheckSentryProjectExists(rn, &projectID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "team", teamName),
			resource.TestCheckResourceAttr(rn, "name", projectName),
			resource.TestCheckResourceAttrSet(rn, "slug"),
			resource.TestCheckResourceAttr(rn, "platform", "go"),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &projectID),
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

func TestAccSentryProject_changeTeam(t *testing.T) {
	teamName1 := acctest.RandomWithPrefix("tf-team")
	teamName2 := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	check := func(teamName, projectName string) resource.TestCheckFunc {
		var projectID string

		return resource.ComposeTestCheckFunc(
			testAccCheckSentryProjectExists(rn, &projectID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "team", teamName),
			resource.TestCheckResourceAttr(rn, "name", projectName),
			resource.TestCheckResourceAttrSet(rn, "slug"),
			resource.TestCheckResourceAttr(rn, "platform", "go"),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &projectID),
			resource.TestCheckResourceAttrPair(rn, "project_id", rn, "internal_id"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectConfig_changeTeam(teamName1, teamName2, projectName, "test_1"),
				Check:  check(teamName1, projectName),
			},
			{
				Config: testAccSentryProjectConfig_changeTeam(teamName1, teamName2, projectName, "test_2"),
				Check:  check(teamName2, projectName),
			},
			{
				Config: testAccSentryProjectRemoveKeyConfig,
				Check:  testAccCheckSentryKeyRemoved("sentry_project.test_project_remove_key"),
			},
			{
				Config: testAccSentryProjectRemoveRuleConfig,
				Check:  testAccCheckSentryRuleRemoved("sentry_project.test_project_remove_rule"),
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

func testAccCheckSentryProjectExists(n string, projectID *string) resource.TestCheckFunc {
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
		*projectID = gotProj.ID
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
	return testAccSentryTeamConfig(teamName) + fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	team         = sentry_team.test.slug
	name         = "%[1]s"
	platform     = "go"
}
	`, projectName)
}

func testAccSentryProjectConfig_changeTeam(teamName1, teamName2, projectName, teamResourceName string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test_1" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_team" "test_2" {
	organization = data.sentry_organization.test.id
	name         = "%[2]s"
	slug         = "%[2]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.%[4]s.organization
	team         = sentry_team.%[4]s.slug
	name         = "%[3]s"
	platform     = "go"
}
	`, teamName1, teamName2, projectName, teamResourceName)
}
