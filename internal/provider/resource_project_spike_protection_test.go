package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccProjectSpikeProtectionResource(t *testing.T) {
	rn := "sentry_project_spike_protection.test"
	teamSlug := acctest.RandomWithPrefix("tf-team")
	projectSlug := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSpikeProtectionConfig(teamSlug, projectSlug, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project_slug", projectSlug),
					resource.TestCheckResourceAttr(rn, "enabled", "true"),
				),
			},
			{
				Config: testAccProjectSpikeProtectionConfig(teamSlug, projectSlug, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project_slug", projectSlug),
					resource.TestCheckResourceAttr(rn, "enabled", "false"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectSpikeProtectionConfig(teamName string, projectName string, enabled bool) string {
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

resource "sentry_project_spike_protection" "test" {
	organization = sentry_project.test.organization
	project_slug = sentry_project.test.slug
	enabled      = %[3]t
}
`, teamName, projectName, enabled)
}
