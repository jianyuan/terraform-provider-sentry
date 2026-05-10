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
	project := acctest.RandomWithPrefix("tf-project")
	fallThrough := false
	codeownersAutoSync := false
	autoAssignment := "Auto Assign to Issue Owner"
	raw := fmt.Sprintf("path:src/views/* #%s", acctest.TestTeam.Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOwnershipConfig(project, fallThrough, codeownersAutoSync, autoAssignment, raw),
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
	project := acctest.RandomWithPrefix("tf-project")
	fallThrough := false
	codeownersAutoSync := false
	autoAssignment := "This auto-assignment mode is not supported"
	raw := fmt.Sprintf("path:src/views/* #%s", acctest.TestTeam.Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectOwnershipConfig(project, fallThrough, codeownersAutoSync, autoAssignment, raw),
				ExpectError: regexp.MustCompile(`Attribute auto_assignment value must be one of`),
			},
		},
	})
}

func testAccProjectOwnershipConfig(projectName string, fallThrough bool, codeownersAutoSync bool, autoAssignment string, raw string) string {
	return fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = "%[1]s"
	teams        = ["%[2]s"]
	name         = "%[3]s"
	platform     = "go"
}

resource "sentry_project_ownership" "test" {
	organization         = sentry_project.test.organization
	project              = sentry_project.test.id
	fallthrough          = %[4]t
	codeowners_auto_sync = %[5]t
	auto_assignment      = "%[6]s"
	raw                  = "%[7]s"
}
`, acctest.TestOrganization, acctest.TestTeam.Slug, projectName, fallThrough, codeownersAutoSync, autoAssignment, raw)
}
