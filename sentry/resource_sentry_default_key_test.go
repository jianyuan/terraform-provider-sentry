package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryDefaultKey_basic(t *testing.T) {
	var key sentry.ProjectKey

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryDefaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryDefaultKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryDefaultKeyExists("sentry_default_key.test_key", &key),
					resource.TestCheckResourceAttr("sentry_default_key.test_key", "name", "Test key"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "public"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "secret"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "dsn_secret"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "dsn_public"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "dsn_csp"),
				),
			},
			{
				Config: testAccSentryDefaultKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryDefaultKeyExists("sentry_default_key.test_key", &key),
					resource.TestCheckResourceAttr("sentry_default_key.test_key", "name", "Test key changed"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "public"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "secret"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "dsn_secret"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "dsn_public"),
					resource.TestCheckResourceAttrSet("sentry_default_key.test_key", "dsn_csp"),
				),
			},
		},
	})
}

func TestAccSentryDefaultKey_RateLimit(t *testing.T) {
	var key sentry.ProjectKey

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryDefaultKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryDefaultKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryDefaultKeyExists("sentry_default_key.test_key", &key),
					resource.TestCheckResourceAttr("sentry_default_key.test_key", "rate_limit_window", "86400"),
					resource.TestCheckResourceAttr("sentry_default_key.test_key", "rate_limit_count", "1000"),
				),
			},
			{
				Config: testAccSentryDefaultKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryDefaultKeyExists("sentry_default_key.test_key", &key),
					resource.TestCheckResourceAttr("sentry_default_key.test_key", "rate_limit_window", "100"),
					resource.TestCheckResourceAttr("sentry_default_key.test_key", "rate_limit_count", "100"),
				),
			},
		},
	})
}

func testAccCheckSentryDefaultKeyDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckSentryDefaultKeyExists(n string, projectKey *sentry.ProjectKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No key ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		keys, _, err := client.ProjectKeys.List(
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
		)
		if err != nil {
			return err
		}

		for _, key := range keys {
			if key.ID == rs.Primary.ID {
				*projectKey = key
				break
			}
		}
		return nil
	}
}

var testAccSentryDefaultKeyConfig = fmt.Sprintf(`
	resource "sentry_team" "test_team" {
		organization = "%[1]s"
		name = "Test team"
	}

	resource "sentry_project" "test_project" {
		organization = "%[1]s"
		team = "${sentry_team.test_team.id}"
		name = "Test project"
	}

	resource "sentry_default_key" "test_key" {
		organization = "%[1]s"
		project = "${sentry_project.test_project.id}"
		name = "Test key"
		rate_limit_window = 86400
		rate_limit_count = 1000
	}
`, testOrganization)

var testAccSentryDefaultKeyUpdateConfig = fmt.Sprintf(`
	resource "sentry_team" "test_team" {
		organization = "%[1]s"
		name = "Test team"
	}

	resource "sentry_project" "test_project" {
		organization = "%[1]s"
		team = "${sentry_team.test_team.id}"
		name = "Test project"
	}

	resource "sentry_default_key" "test_key" {
		organization = "%[1]s"
		project = "${sentry_project.test_project.id}"
		name = "Test key changed"
		rate_limit_window = 100
		rate_limit_count = 100
	}
`, testOrganization)
