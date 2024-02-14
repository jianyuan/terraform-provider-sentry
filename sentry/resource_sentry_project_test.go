package sentry

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryProject_basic(t *testing.T) {
	teamName1 := acctest.RandomWithPrefix("tf-team")
	teamName2 := acctest.RandomWithPrefix("tf-team")
	teamName3 := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	check := func(projectName string, teamNames []string) resource.TestCheckFunc {
		var projectID string

		fs := resource.ComposeTestCheckFunc(
			testAccCheckSentryProjectExists(rn, &projectID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "teams.#", strconv.Itoa(len(teamNames))),
			resource.TestCheckResourceAttr(rn, "name", projectName),
			resource.TestCheckResourceAttrSet(rn, "slug"),
			resource.TestCheckResourceAttr(rn, "platform", "go"),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &projectID),
			resource.TestCheckResourceAttrPair(rn, "project_id", rn, "internal_id"),
		)
		for _, teamName := range teamNames {
			fs = resource.ComposeTestCheckFunc(fs, resource.TestCheckTypeSetElemAttr(rn, "teams.*", teamName))
		}
		return fs
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectConfig_teams([]string{teamName1}, projectName),
				Check:  check(projectName, []string{teamName1}),
			},
			{
				Config: testAccSentryProjectConfig_teams([]string{teamName2}, projectName+"-renamed"),
				Check:  check(projectName+"-renamed", []string{teamName2}),
			},
			{
				Config: testAccSentryProjectConfig_teams([]string{teamName2, teamName3}, projectName+"-renamed"),
				Check:  check(projectName+"-renamed", []string{teamName2, teamName3}),
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

func TestAccSentryProject_teamMigration(t *testing.T) {
	teams := []string{
		acctest.RandomWithPrefix("tf-team"),
		acctest.RandomWithPrefix("tf-team"),
		acctest.RandomWithPrefix("tf-team"),
	}
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	check := func(team string, teams []string) resource.TestCheckFunc {
		var projectID string

		fs := resource.ComposeTestCheckFunc(
			testAccCheckSentryProjectExists(rn, &projectID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "name", projectName),
			resource.TestCheckResourceAttrSet(rn, "slug"),
			resource.TestCheckResourceAttr(rn, "platform", "go"),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &projectID),
			resource.TestCheckResourceAttrPair(rn, "project_id", rn, "internal_id"),
		)
		if team != "" {
			fs = resource.ComposeTestCheckFunc(fs, resource.TestCheckResourceAttr(rn, "team", team))
		}
		for _, team := range teams {
			fs = resource.ComposeTestCheckFunc(fs, resource.TestCheckTypeSetElemAttr(rn, "teams.*", team))
		}

		return fs
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectConfig_teams_old(teams, projectName),
				Check:  check(teams[0], nil),
			},
			{
				Config: testAccSentryProjectConfig_teams(teams, projectName),
				Check:  check("", teams),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       testAccSentryProjectImportStateIdFunc(rn),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"team"},
			},
		},
	})
}

func TestAccSentryProject_deprecatedTeam(t *testing.T) {
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
				Config: testAccSentryProjectConfig_team(teamName, projectName),
				Check:  check(projectName),
			},
			{
				Config: testAccSentryProjectConfig_team(teamName, projectName+"-renamed"),
				Check:  check(projectName + "-renamed"),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       testAccSentryProjectImportStateIdFunc(rn),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"team"},
			},
		},
	})
}

func TestAccSentryProject_noTeam(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSentryProjectConfig_noTeam(teamName, projectName),
				ExpectError: regexp.MustCompile("one of team or teams must be configured"),
			},
		},
	})
}

func TestAccSentryProject_teamConflict(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSentryProjectConfig_teamConflict(teamName, projectName),
				ExpectError: regexp.MustCompile("\"team\": conflicts with teams"),
			},
			{
				Config:      testAccSentryProjectConfig_teamConflict(teamName, projectName),
				ExpectError: regexp.MustCompile("\"teams\": conflicts with team"),
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

func testAccSentryProjectConfig_team(teamName, projectName string) string {
	return testAccSentryTeamConfig(teamName) + fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	team         = sentry_team.test.slug
	name         = "%[1]s"
	platform     = "go"
}
	`, projectName)
}

func testAccSentryProjectConfig_noTeam(teamName, projectName string) string {
	return testAccSentryTeamConfig(teamName) + fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	name         = "%[1]s"
	platform     = "go"
}
	`, projectName)
}

func testAccSentryProjectConfig_teamConflict(teamName, projectName string) string {
	return testAccSentryTeamConfig(teamName) + fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	team         = sentry_team.test.slug
	teams        = [sentry_team.test.slug]
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

func testAccSentryProjectConfig_teams_old(teamNames []string, projectName string) string {
	config := testAccSentryOrganizationDataSourceConfig

	teamSlugs := make([]string, 0, len(teamNames))
	for i, teamName := range teamNames {
		config += fmt.Sprintf(`
resource "sentry_team" "test_%[1]d" {
	organization = data.sentry_organization.test.id
	name         = "%[2]s"
	slug         = "%[2]s"
}
		`, i, teamName)
		teamSlugs = append(teamSlugs, fmt.Sprintf("sentry_team.test_%d.slug", i))
	}

	config += fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test_0.organization
	team         = %[2]s
	name         = "%[1]s"
	platform     = "go"
}
	`, projectName, teamSlugs[0])

	return config
}

func testAccSentryProjectConfig_teams(teamNames []string, projectName string) string {
	config := testAccSentryOrganizationDataSourceConfig

	teamSlugs := make([]string, 0, len(teamNames))
	for i, teamName := range teamNames {
		config += fmt.Sprintf(`
resource "sentry_team" "test_%[1]d" {
	organization = data.sentry_organization.test.id
	name         = "%[2]s"
	slug         = "%[2]s"
}
		`, i, teamName)
		teamSlugs = append(teamSlugs, fmt.Sprintf("sentry_team.test_%d.slug", i))
	}

	config += fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test_0.organization
	teams        = [%[2]s]
	name         = "%[1]s"
	platform     = "go"
}
	`, projectName, strings.Join(teamSlugs, ", "))

	return config
}
