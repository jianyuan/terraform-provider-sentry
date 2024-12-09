package acctest

import (
	"context"
	"os"
	"testing"

	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/go-utils/must"
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

	// TestPagerDutyOrganization is the PagerDuty organization used for acceptance tests.
	TestPagerDutyOrganization = os.Getenv("SENTRY_TEST_PAGERDUTY_ORGANIZATION")

	// TestOpsgenieOrganization is the Opsgenie organization used for acceptance tests.
	TestOpsgenieOrganization = os.Getenv("SENTRY_TEST_OPSGENIE_ORGANIZATION")
	// TestOpsgenieIntegrationKey is the Opsgenie integration key used for acceptance tests.
	TestOpsgenieIntegrationKey = os.Getenv("SENTRY_TEST_OPSGENIE_INTEGRATION_KEY")

	// TestGitHubInstallationId is the GitHub installation ID used for acceptance tests.
	TestGitHubInstallationId = os.Getenv("SENTRY_TEST_GITHUB_INSTALLATION_ID")

	// TestGitHubRepositoryIdentifier is the GitHub repository identifier used for acceptance tests.
	TestGitHubRepositoryIdentifier = os.Getenv("SENTRY_TEST_GITHUB_REPOSITORY_IDENTIFIER")

	// TestGitLabInstallationId is the GitLab installation ID used for acceptance tests.
	TestGitLabInstallationId = os.Getenv("SENTRY_TEST_GITLAB_INSTALLATION_ID")

	// TestGitLabRepositoryIdentifier is the GitLab repository identifier used for acceptance tests.
	TestGitLabRepositoryIdentifier = os.Getenv("SENTRY_TEST_GITLAB_REPOSITORY_IDENTIFIER")

	// TestVSTSInstallationId is the VSTS installation ID used for acceptance tests.
	TestVSTSInstallationId = os.Getenv("SENTRY_TEST_VSTS_INSTALLATION_ID")

	// TestVSTSRepositoryIdentifier is the VSTS repository identifier used for acceptance tests.
	TestVSTSRepositoryIdentifier = os.Getenv("SENTRY_TEST_VSTS_REPOSITORY_IDENTIFIER")

	// SharedClient is a shared Sentry client for acceptance tests.
	SharedClient *sentry.Client
)

func init() {
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
		UserAgent: "Terraform/" + ProviderVersion + " (+https://www.terraform.io) terraform-provider-sentry/" + ProviderVersion,
		Token:     token,
	}
	httpClient := config.HttpClient(context.Background())

	if baseUrl == "" {
		SharedClient = sentry.NewClient(httpClient)
	} else {
		SharedClient = must.Get(sentry.NewOnPremiseClient(baseUrl, httpClient))
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
