package sentry

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testOrganization = os.Getenv("SENTRY_TEST_ORGANIZATION")

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"sentry": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SENTRY_TOKEN"); v == "" {
		t.Fatal("SENTRY_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SENTRY_TEST_ORGANIZATION"); v == "" {
		t.Fatal("SENTRY_TEST_ORGANIZATION must be set for acceptance tests")
	}
}
