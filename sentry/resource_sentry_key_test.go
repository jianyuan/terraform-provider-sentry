package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryKey_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	check := func(keyName string) resource.TestCheckFunc {
		var keyID string

		return resource.ComposeTestCheckFunc(
			testAccCheckSentryKeyExists(rn, &keyID),
			resource.TestCheckResourceAttrPtr(rn, "id", &keyID),
			resource.TestCheckResourceAttr(rn, "name", keyName),
			resource.TestCheckResourceAttrSet(rn, "public"),
			resource.TestCheckResourceAttrSet(rn, "secret"),
			resource.TestCheckResourceAttrSet(rn, "dsn_secret"),
			resource.TestCheckResourceAttrSet(rn, "dsn_public"),
			resource.TestCheckResourceAttrSet(rn, "dsn_csp"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyConfig(teamName, projectName, keyName),
				Check:  check(keyName),
			},
			{
				Config: testAccSentryKeyConfig(teamName, projectName, keyName+"-renamed"),
				Check:  check(keyName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccSentryKeyImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSentryKey_RateLimit(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	check := func(keyName, rateLimitWindow, rateLimitCount string) resource.TestCheckFunc {
		var keyID string

		return resource.ComposeTestCheckFunc(
			testAccCheckSentryKeyExists(rn, &keyID),
			resource.TestCheckResourceAttrPtr(rn, "id", &keyID),
			resource.TestCheckResourceAttr(rn, "name", keyName),
			resource.TestCheckResourceAttrSet(rn, "public"),
			resource.TestCheckResourceAttrSet(rn, "secret"),
			resource.TestCheckResourceAttrSet(rn, "dsn_secret"),
			resource.TestCheckResourceAttrSet(rn, "dsn_public"),
			resource.TestCheckResourceAttrSet(rn, "dsn_csp"),
			resource.TestCheckResourceAttr(rn, "rate_limit_window", rateLimitWindow),
			resource.TestCheckResourceAttr(rn, "rate_limit_count", rateLimitCount),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryKeyConfig_rateLimit(teamName, projectName, keyName, "86400", "1000"),
				Check:  check(keyName, "86400", "1000"),
			},
			{
				Config: testAccSentryKeyConfig_rateLimit(teamName, projectName, keyName, "100", "100"),
				Check:  check(keyName, "100", "100"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccSentryKeyImportStateIdFunc(rn),
				ImportStateVerify: true,
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

		ctx := context.Background()
		keys, resp, err := client.ProjectKeys.List(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
			nil,
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

func testAccCheckSentryKeyExists(n string, keyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no key ID is set")
		}

		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		keys, _, err := client.ProjectKeys.List(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["project"],
			nil,
		)
		if err != nil {
			return err
		}

		for _, key := range keys {
			if key.ID == rs.Primary.ID {
				*keyID = key.ID
				return nil
			}
		}
		return fmt.Errorf("not found: %s", n)
	}
}

func testAccSentryKeyImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		org := rs.Primary.Attributes["organization"]
		proj := rs.Primary.Attributes["project"]
		keyID := rs.Primary.ID
		return buildThreePartID(org, proj, keyID), nil
	}
}

func testAccSentryKeyConfig(teamName, projectName, keyName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + fmt.Sprintf(`
resource "sentry_key" "test" {
	organization      = sentry_project.test.organization
	project           = sentry_project.test.id
	name              = "%[1]s"
}
	`, keyName)
}

func testAccSentryKeyConfig_rateLimit(teamName, projectName, keyName, rateLimitWindow, rateLimitCount string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + fmt.Sprintf(`
resource "sentry_key" "test" {
	organization      = sentry_project.test.organization
	project           = sentry_project.test.id
	name              = "%[1]s"
	rate_limit_window = %[2]s
	rate_limit_count  = %[3]s
}
	`, keyName, rateLimitWindow, rateLimitCount)
}
