package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccProjectInboundDataFilterResource(t *testing.T) {
	rn := "sentry_project_inbound_data_filter.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	filterId := "browser-extensions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(team, project, filterId, "active = true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckResourceAttr(rn, "active", "true"),
					resource.TestCheckNoResourceAttr(rn, "subfilters"),
				),
			},
			{
				Config: testAccProjectInboundDataFilterConfig(team, project, filterId, "active = false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckResourceAttr(rn, "active", "false"),
					resource.TestCheckNoResourceAttr(rn, "subfilters"),
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

func TestAccProjectInboundDataFilterResource_LegacyBrowser(t *testing.T) {
	rn := "sentry_project_inbound_data_filter.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	filterId := "legacy-browsers"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(team, project, filterId, "subfilters = [\"android_pre_4\", \"ie_pre_9\"]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckNoResourceAttr(rn, "active"),
					resource.TestCheckResourceAttr(rn, "subfilters.#", "2"),
					resource.TestCheckResourceAttr(rn, "subfilters.0", "android_pre_4"),
					resource.TestCheckResourceAttr(rn, "subfilters.1", "ie_pre_9"),
				),
			},
			// NOTE: Intentionally not sorting subfilters to show that the order does not matter during import.
			{
				Config: testAccProjectInboundDataFilterConfig(team, project, filterId, "subfilters = [\"safari_pre_6\", \"android_pre_4\"]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckNoResourceAttr(rn, "active"),
					resource.TestCheckResourceAttr(rn, "subfilters.#", "2"),
					resource.TestCheckResourceAttr(rn, "subfilters.0", "safari_pre_6"),
					resource.TestCheckResourceAttr(rn, "subfilters.1", "android_pre_4"),
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

func testAccProjectInboundDataFilterConfig(teamName, projectName, filterId, body string) string {
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

resource "sentry_project_inbound_data_filter" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	filter_id    = "%[3]s"
	%[4]s
}
`, teamName, projectName, filterId, body)
}
