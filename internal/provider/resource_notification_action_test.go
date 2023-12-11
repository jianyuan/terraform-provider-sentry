package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccNotificationActionResource(t *testing.T) {
	rn := "sentry_notification_action.test"
	teamSlug := acctest.RandomWithPrefix("tf-team")
	project1Slug := acctest.RandomWithPrefix("tf-project")
	project2Slug := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationActionConfig(teamSlug, project1Slug, project2Slug, "[sentry_project.test_1.slug]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "trigger_type", "spike-protection"),
					resource.TestCheckResourceAttr(rn, "service_type", "sentry_notification"),
					resource.TestCheckResourceAttr(rn, "target_identifier", "default"),
					resource.TestCheckResourceAttr(rn, "target_display", "default"),
					resource.TestCheckResourceAttr(rn, "project_slugs.#", "1"),
					resource.TestCheckResourceAttr(rn, "project_slugs.0", project1Slug),
				),
			},
			{
				Config: testAccNotificationActionConfig(teamSlug, project1Slug, project2Slug, "[sentry_project.test_2.slug]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "trigger_type", "spike-protection"),
					resource.TestCheckResourceAttr(rn, "service_type", "sentry_notification"),
					resource.TestCheckResourceAttr(rn, "target_identifier", "default"),
					resource.TestCheckResourceAttr(rn, "target_display", "default"),
					resource.TestCheckResourceAttr(rn, "project_slugs.#", "1"),
					resource.TestCheckResourceAttr(rn, "project_slugs.0", project2Slug),
				),
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					org := rs.Primary.Attributes["organization"]
					actionId := rs.Primary.ID
					return buildTwoPartID(org, actionId), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotificationActionConfig(teamName string, project1Name string, project2Name string, projectSlugs string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
  organization = data.sentry_organization.test.id
  name         = "%[1]s"
  slug         = "%[1]s"
}

resource "sentry_project" "test_1" {
  organization = sentry_team.test.organization
  teams        = [sentry_team.test.id]
  name         = "%[2]s"
  platform     = "go"
}

resource "sentry_project" "test_2" {
  organization = sentry_team.test.organization
  teams        = [sentry_team.test.id]
  name         = "%[3]s"
  platform     = "go"
}

resource "sentry_notification_action" "test" {
  organization      = sentry_team.test.organization
  trigger_type      = "spike-protection"
  service_type      = "sentry_notification"
  target_identifier = "default"
  target_display    = "default"
  project_slugs     = %[4]s
}
`, teamName, project1Name, project2Name, projectSlugs)
}
