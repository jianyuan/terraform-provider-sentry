package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSentryKey_basic(t *testing.T) {
	var key Key

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

func testAccCheckSentryKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_key" {
			continue
		}

		key, resp, err := client.GetKey(
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
			rs.Primary.ID,
		)
		if err == nil {
			if key != nil {
				return errors.New("Key still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryKeyExists(n string, key *Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No key ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		sentryKey, _, err := client.GetKey(
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		}
		*key = *sentryKey
		return nil
	}
}

var testAccSentryKeyConfig = fmt.Sprintf(`
	resource "sentry_team" "test_team" {
		organization = "%s"
		name = "Test team"
	}

	resource "sentry_project" "test_project" {
		organization = "%s"
		team = "${sentry_team.test_team.id}"
		name = "Test project"
	}

	resource "sentry_key" "test_key" {
		organization = "%s"
		project = "${sentry_project.test_project.id}"
		name = "Test key"
	}
`, testOrganization, testOrganization, testOrganization)

var testAccSentryKeyUpdateConfig = fmt.Sprintf(`
	resource "sentry_team" "test_team" {
		organization = "%s"
		name = "Test team"
	}

	resource "sentry_project" "test_project" {
		organization = "%s"
		team = "${sentry_team.test_team.id}"
		name = "Test project"
	}

	resource "sentry_key" "test_key" {
		organization = "%s"
		project = "${sentry_project.test_project.id}"
		name = "Test key changed"
	}
`, testOrganization, testOrganization, testOrganization)
