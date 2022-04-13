package sentry

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/canva/go-sentry/sentry"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryProjectFilter_basic(t *testing.T) {
	var filterConfig sentry.FilterConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryProjectFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectFilterConfig,
				Check:  testFilterConfig("sentry_filter.test_filter", &filterConfig, true, []string{"ie_pre_9", "ie10"}),
			},
		},
	})
}

func testAccCheckSentryProjectFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		filterConfig, _, err := client.ProjectFilter.GetFilterConfig(testOrganization, rs.Primary.Attributes["project"])
		if err != nil {
			return err
		}
		if filterConfig.BrowserExtension != false {
			return fmt.Errorf("got browser_extension %t; want false", filterConfig.BrowserExtension)
		}
		if len(filterConfig.LegacyBrowsers) != 0 {
			return fmt.Errorf("got legacy_browser %v; want []", filterConfig.LegacyBrowsers)
		}
		return nil
	}

	return nil
}

func testFilterConfig(n string, filterConfig *sentry.FilterConfig, browserExtension bool, legacyBrowsers []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		filterConfig, _, err := client.ProjectFilter.GetFilterConfig(testOrganization, rs.Primary.Attributes["project"])
		if err != nil {
			return err
		}
		if filterConfig.BrowserExtension != browserExtension {
			return fmt.Errorf("got browser_extension %t; want %t", filterConfig.BrowserExtension, browserExtension)
		}

		if !cmp.Equal(filterConfig.LegacyBrowsers, legacyBrowsers, cmp.Transformer("sort", func(in []string) []string {
			sort.Strings(in)
			return in
		})) {
			return fmt.Errorf("got legacy_browser %v; want %v", filterConfig.LegacyBrowsers, legacyBrowsers)
		}

		return nil
	}
}

var testAccSentryProjectFilterConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
	organization = "%s"
	name         = "Test team"
}

resource "sentry_project" "test_project" {
	organization = "%s"
	team         = sentry_team.test_team.id
	name         = "Test project"
	platform     = "go"
}

resource "sentry_filter" "test_filter" {
	organization = "%s"
	project = sentry_project.test_project.id
	browser_extension = true
	legacy_browsers = ["ie_pre_9","ie10"]
}
`, testOrganization, testOrganization, testOrganization)
