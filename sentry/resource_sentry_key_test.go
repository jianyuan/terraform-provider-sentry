package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryKey_basic(t *testing.T) {
	var key sentry.ProjectKey

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyExists("sentry_key.test_key", &key),
					resource.TestCheckResourceAttr("sentry_key.test_key", "name", "Test key"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "public"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "secret"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "dsn_secret"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "dsn_public"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "dsn_csp"),
				),
			},
			{
				Config: testAccSentryKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyExists("sentry_key.test_key", &key),
					resource.TestCheckResourceAttr("sentry_key.test_key", "name", "Test key changed"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "public"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "secret"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "dsn_secret"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "dsn_public"),
					resource.TestCheckResourceAttrSet("sentry_key.test_key", "dsn_csp"),
				),
			},
		},
	})
}

func TestAccSentryKey_RateLimit(t *testing.T) {
	var key sentry.ProjectKey

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyExists("sentry_key.test_key_rate_limit", &key),
					resource.TestCheckResourceAttr("sentry_key.test_key_rate_limit", "rate_limit_window", "86400"),
					resource.TestCheckResourceAttr("sentry_key.test_key_rate_limit", "rate_limit_count", "1000"),
				),
			},
			{
				Config: testAccSentryKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryKeyExists("sentry_key.test_key_rate_limit", &key),
					resource.TestCheckResourceAttr("sentry_key.test_key_rate_limit", "rate_limit_window", "100"),
					resource.TestCheckResourceAttr("sentry_key.test_key_rate_limit", "rate_limit_count", "100"),
				),
			},
		},
	})
}

func testAccCheckSentryKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_key" {
			continue
		}

		keys, resp, err := client.ProjectKeys.List(
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
		)
		if err == nil {
			for _, key := range keys {
				if key.ID == rs.Primary.ID {
					return errors.New("Key still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryKeyExists(n string, projectKey *sentry.ProjectKey) resource.TestCheckFunc {
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

var testAccSentryKeyConfig = fmt.Sprintf(`
	resource "sentry_team" "test_team" {
		organization = "%[1]s"
		name = "Test team"
	}

	resource "sentry_project" "test_project" {
		organization = "%[1]s"
		team = "${sentry_team.test_team.id}"
		name = "Test project"
	}

	resource "sentry_key" "test_key" {
		organization = "%[1]s"
		project = "${sentry_project.test_project.id}"
		name = "Test key"
	}

	resource "sentry_key" "test_key_rate_limit" {
		organization = "%[1]s"
		project = "${sentry_project.test_project.id}"
		name = "Test key"
		rate_limit_window = 86400
		rate_limit_count = 1000
	}
`, testOrganization)

var testAccSentryKeyUpdateConfig = fmt.Sprintf(`
	resource "sentry_team" "test_team" {
		organization = "%[1]s"
		name = "Test team"
	}

	resource "sentry_project" "test_project" {
		organization = "%[1]s"
		team = "${sentry_team.test_team.id}"
		name = "Test project"
	}

	resource "sentry_key" "test_key" {
		organization = "%[1]s"
		project = "${sentry_project.test_project.id}"
		name = "Test key changed"
	}

	resource "sentry_key" "test_key_rate_limit" {
		organization = "%[1]s"
		project = "${sentry_project.test_project.id}"
		name = "Test key"
		rate_limit_window = 100
		rate_limit_count = 100
	}
`, testOrganization)
