package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccProjectInboundDataFilterResource(t *testing.T) {
	rn := "sentry_project_inbound_data_filter.test"
	project := acctest.RandomWithPrefix("tf-project")
	filterId := "browser-extensions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(project, filterId, "active = true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckResourceAttr(rn, "active", "true"),
					resource.TestCheckNoResourceAttr(rn, "subfilters"),
				),
			},
			{
				Config: testAccProjectInboundDataFilterConfig(project, filterId, "active = false"),
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

func TestAccProjectInboundDataFilterResource_Conflict(t *testing.T) {
	project := acctest.RandomWithPrefix("tf-project")
	filterId := "browser-extensions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(project, filterId, `
					active = true
					subfilters = ["android_pre_4", "ie_pre_9"]
				`),
				ExpectError: regexp.MustCompile(`Attribute "active" cannot be specified when "subfilters" is specified`),
			},
		},
	})
}

func TestAccProjectInboundDataFilterResource_LegacyBrowser(t *testing.T) {
	rn := "sentry_project_inbound_data_filter.test"
	project := acctest.RandomWithPrefix("tf-project")
	filterId := "legacy-browsers"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectInboundDataFilterConfig(project, filterId, "subfilters = [\"android_pre_4\", \"ie_pre_9\"]"),
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
				Config: testAccProjectInboundDataFilterConfig(project, filterId, "subfilters = [\"safari_pre_6\", \"android_pre_4\"]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "filter_id", filterId),
					resource.TestCheckNoResourceAttr(rn, "active"),
					resource.TestCheckResourceAttr(rn, "subfilters.#", "2"),
					resource.TestCheckResourceAttr(rn, "subfilters.0", "android_pre_4"),
					resource.TestCheckResourceAttr(rn, "subfilters.1", "safari_pre_6"),
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

func testAccProjectInboundDataFilterConfig(projectName, filterId, body string) string {
	return fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = "%[1]s"
	teams        = ["%[2]s"]
	name         = "%[3]s"
	platform     = "go"
}

resource "sentry_project_inbound_data_filter" "test" {
	organization = "%[1]s"
	project      = sentry_project.test.id
	filter_id    = "%[4]s"
	%[5]s
}
`, acctest.TestOrganization, acctest.TestTeam.Slug, projectName, filterId, body)
}
