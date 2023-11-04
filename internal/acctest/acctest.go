package acctest

import (
	"context"
	"os"

	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

// SharedClient is a shared Sentry client for acceptance tests.
var SharedClient *sentry.Client

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
