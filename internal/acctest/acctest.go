package acctest

import (
	"context"
	"os"
	"testing"

	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

const (
	// ProviderName is the name of the Terraform provider.
	ProviderName = "sentry"

	// ProviderVersion is the version of the Terraform provider.
	ProviderVersion = "test"
)

var (
	// TestOrganization is the organization used for acceptance tests.
	TestOrganization = os.Getenv("SENTRY_TEST_ORGANIZATION")

	// SharedClient is a shared Sentry client for acceptance tests.
	SharedClient *sentry.Client
)

func init() {
	var err error
	var token string
	if v := os.Getenv("SENTRY_AUTH_TOKEN"); v != "" {
		token = v
	} else if v := os.Getenv("SENTRY_TOKEN"); v != "" {
		token = v
	}

	var baseUrl string
	if v := os.Getenv("SENTRY_BASE_URL"); v != "" {
		baseUrl = v
	} else {
		baseUrl = "https://sentry.io/api/"
	}

	config := sentryclient.Config{
		Token:   token,
		BaseURL: baseUrl,
	}
	SharedClient, err = config.Client(context.Background())
	if err != nil {
		panic(err)
	}
}

func PreCheck(t *testing.T) {
	if v := os.Getenv("SENTRY_AUTH_TOKEN"); v == "" {
		t.Fatal("SENTRY_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SENTRY_TEST_ORGANIZATION"); v == "" {
		t.Fatal("SENTRY_TEST_ORGANIZATION must be set for acceptance tests")
	}
}
