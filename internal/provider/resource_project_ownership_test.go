package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccProjectOwnershipResource(t *testing.T) {
	rn := "sentry_project_ownership.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	fallThrough := false
	codeownersAutoSync := false
	autoAssignment := "Auto Assign to Issue Owner"
	raw := fmt.Sprintf("path:src/views/* #%s", team)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOwnershipConfig(team, project, fallThrough, codeownersAutoSync, autoAssignment, raw),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "fallthrough", strconv.FormatBool(fallThrough)),
					resource.TestCheckResourceAttr(rn, "codeowners_auto_sync", strconv.FormatBool(codeownersAutoSync)),
					resource.TestCheckResourceAttr(rn, "auto_assignment", autoAssignment),
					resource.TestCheckResourceAttr(rn, "raw", raw),
				),
			},
		},
	})
}

func TestAccProjectOwnershipResource_IllegalAutoAssignment(t *testing.T) {
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	fallThrough := false
	codeownersAutoSync := false
	autoAssignment := "This auto-assignment mode is not supported"
	raw := fmt.Sprintf("path:src/views/* #%s", team)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectOwnershipConfig(team, project, fallThrough, codeownersAutoSync, autoAssignment, raw),
				ExpectError: regexp.MustCompile(`Attribute auto_assignment value must be one of`),
			},
		},
	})
}

func testAccProjectOwnershipConfig(teamName string, projectName string, fallThrough bool, codeownersAutoSync bool, autoAssignment string, raw string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_project_ownership" "test" {
	organization         = sentry_project.test.organization
	project              = sentry_project.test.id
	fallthrough          = %[3]t
	codeowners_auto_sync = %[4]t
	auto_assignment      = "%[5]s"
	raw                  = "%[6]s"
}
`, teamName, projectName, fallThrough, codeownersAutoSync, autoAssignment, raw)
}
