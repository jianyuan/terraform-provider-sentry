package sentry

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testOrganization = os.Getenv("SENTRY_TEST_ORGANIZATION")

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = NewProvider("dev")()
	testAccProviders = map[string]*schema.Provider{
		"sentry": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sentry": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := NewProvider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SENTRY_AUTH_TOKEN"); v == "" {
		t.Fatal("SENTRY_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SENTRY_TEST_ORGANIZATION"); v == "" {
		t.Fatal("SENTRY_TEST_ORGANIZATION must be set for acceptance tests")
	}
}
