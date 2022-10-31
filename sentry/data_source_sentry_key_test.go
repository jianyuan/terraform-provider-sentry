package sentry

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryKeyDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	dn := "data.sentry_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyDataSourceConfig(teamName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyDataSourceID(dn),
					resource.TestCheckResourceAttr(dn, "name", "Default"),
					resource.TestMatchResourceAttr(dn, "public", regexp.MustCompile(`^[0-9a-f]+$`)),
					resource.TestMatchResourceAttr(dn, "secret", regexp.MustCompile(`^[0-9a-f]+$`)),
					resource.TestMatchResourceAttr(dn, "project_id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrSet(dn, "is_active"),
					resource.TestMatchResourceAttr(dn, "dsn_secret", regexp.MustCompile(`^https://`)),
					resource.TestMatchResourceAttr(dn, "dsn_public", regexp.MustCompile(`^https://`)),
					resource.TestMatchResourceAttr(dn, "dsn_csp", regexp.MustCompile(`^https://`)),
				),
			},
		},
	})
}

func TestAccSentryKeyDataSource_first(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	dn := "data.sentry_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyDataSourceConfig_first(teamName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyDataSourceID(dn),
					resource.TestCheckResourceAttr(dn, "name", "Default"),
				),
			},
		},
	})
}

func TestAccSentryKeyDataSource_name(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	dn := "data.sentry_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyDataSourceConfig_name(teamName, projectName, keyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyDataSourceID(dn),
					resource.TestCheckResourceAttr(dn, "name", keyName),
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

func testAccSentryKeyDataSourceConfig(teamName, projectName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + `
data "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
}
	`
}

// Testing first parameter
func testAccSentryKeyDataSourceConfig_first(teamName, projectName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + `
resource "sentry_key" "test_2" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id

	name = "Test key 2"
}

data "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id

	first = true
}
	`
}

// Testing name parameter
func testAccSentryKeyDataSourceConfig_name(teamName, projectName, keyName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + fmt.Sprintf(`
resource "sentry_key" "test_2" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id

	name = "%[1]s"
}

data "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id

	name = sentry_key.test_2.name
}
	`, keyName)
}
