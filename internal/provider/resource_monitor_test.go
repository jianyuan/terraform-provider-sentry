package provider

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccMonitorResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "sentry_monitor" "test" {
						organization = "value"
						project      = "value"
						name         = "value"

						config = {}
					}
				`,
				ExpectError: regexp.MustCompile(`At least one of these attributes must be configured`),
			},
			{
				Config: `
					resource "sentry_monitor" "test" {
						organization = "value"
						project      = "value"
						name         = "value"
						slug         = "1234"

						config = {
							schedule_crontab = "* * * * *"
						}
					}
				`,
				ExpectError: regexp.MustCompile(`must contain only lowercase letters, numbers, underscores, or dashes and cannot be all digits`),
			},
			{
				Config: `
					resource "sentry_monitor" "test" {
						organization = "value"
						project      = "value"
						name         = "value"
						slug         = "ABC"

						config = {
							schedule_crontab = "* * * * *"
						}
					}
				`,
				ExpectError: regexp.MustCompile(`must contain only lowercase letters, numbers, underscores, or dashes and cannot be all digits`),
			},
			{
				Config: `
					resource "sentry_monitor" "test" {
						organization = "value"
						project      = "value"
						name         = "value"
						slug         = "` + strings.Repeat("a", 51) + `"

						config = {
							schedule_crontab = "* * * * *"
						}
					}
				`,
				ExpectError: regexp.MustCompile(`Attribute must be between 1 and 50 characters`),
			},
			{
				Config: `
					resource "sentry_monitor" "test" {
						organization = "value"
						project      = "value"
						name         = "` + strings.Repeat("a", 129) + `"

						config = {
							schedule_crontab = "* * * * *"
						}
					}
				`,
				ExpectError: regexp.MustCompile(`Attribute must be between 1 and 128 characters`),
			},
		},
	})
}

func TestAccMonitorResource_basic(t *testing.T) {
	rn := "sentry_monitor.test"
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	monitorName := acctest.RandomWithPrefix("tf-monitor")
	monitorSlug := acctest.RandomWithPrefix("tf-monitor")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig(teamName, projectName, monitorName, monitorSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", projectName),
					resource.TestCheckResourceAttr(rn, "name", monitorName),
					resource.TestCheckResourceAttr(rn, "slug", monitorSlug),
					resource.TestCheckResourceAttr(rn, "config.schedule_crontab", "* * * * *"),
				),
			},
			{
				Config: testAccMonitorResourceConfig(teamName, projectName, monitorName+"-updated", monitorSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", monitorName+"-updated"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.TwoPartImportStateIdFunc(rn, "organization"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitorResourceConfig(teamName, projectName, monitorName, monitorSlug string) string {
	return testAccOrganizationDataSourceConfig + `
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.slug
	name         = "` + teamName + `"
	slug         = "` + teamName + `"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.slug]
	name         = "` + projectName + `"
	platform     = "go"
}

resource "sentry_monitor" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "` + monitorName + `"
	slug         = "` + monitorSlug + `"

	config = {
		schedule_crontab = "* * * * *"
	}
}
`
}
