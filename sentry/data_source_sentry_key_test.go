package sentry

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSentryKeyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyDataSourceID("data.sentry_key.test_key"),
					resource.TestCheckResourceAttrSet("data.sentry_key.test_key", "name"),
					resource.TestMatchResourceAttr("data.sentry_key.test_key", "public", regexp.MustCompile(`^[0-9a-f]+$`)),
					resource.TestMatchResourceAttr("data.sentry_key.test_key", "secret", regexp.MustCompile(`^[0-9a-f]+$`)),
					resource.TestMatchResourceAttr("data.sentry_key.test_key", "project_id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrSet("data.sentry_key.test_key", "is_active"),
					resource.TestMatchResourceAttr("data.sentry_key.test_key", "dsn_secret", regexp.MustCompile(`^https://`)),
					resource.TestMatchResourceAttr("data.sentry_key.test_key", "dsn_public", regexp.MustCompile(`^https://`)),
					resource.TestMatchResourceAttr("data.sentry_key.test_key", "dsn_csp", regexp.MustCompile(`^https://`)),
				),
			},
		},
	})
}

func TestAccSentryKeyDataSource_first(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyDataSourceFirstConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyDataSourceID("data.sentry_key.test_key"),
				),
			},
		},
	})
}

func TestAccSentryKeyDataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyDataSourceNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyDataSourceID("data.sentry_key.default_key"),
				),
			},
		},
	})
}

func testAccCheckSentryKeyDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Sentry key: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Sentry key data source ID not set")
		}

		return nil
	}
}

var testAccSentryKeyDataSourceConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
  organization = "%s"
  name = "Test team"
}

resource "sentry_project" "test_project" {
  organization = "%s"
  team = "${sentry_team.test_team.id}"
  name = "Test project"
}

data "sentry_key" "test_key" {
  organization = "%s"
  project = "${sentry_project.test_project.id}"
}
`, testOrganization, testOrganization, testOrganization)

// Testing first parameter
var testAccSentryKeyDataSourceFirstConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
  organization = "%s"
  name = "Test team"
}

resource "sentry_project" "test_project" {
  organization = "%s"
  team = "${sentry_team.test_team.id}"
  name = "Test project"
}

resource "sentry_key" "test_key_2" {
  organization = "%s"
  project = "${sentry_project.test_project.id}"
  name = "Test key 2"
}

data "sentry_key" "test_key" {
  organization = "%s"
  project = "${sentry_project.test_project.id}"
  first = true
}
`, testOrganization, testOrganization, testOrganization, testOrganization)

// Testing name parameter
// A key named "Default" is always created along with the project
var testAccSentryKeyDataSourceNameConfig = fmt.Sprintf(`
resource "sentry_team" "test_team" {
  organization = "%s"
  name = "Test team"
}

resource "sentry_project" "test_project" {
  organization = "%s"
  team = "${sentry_team.test_team.id}"
  name = "Test project"
}

data "sentry_key" "default_key" {
  organization = "%s"
  project = "${sentry_project.test_project.id}"
  name = "Default"
}
`, testOrganization, testOrganization, testOrganization)
