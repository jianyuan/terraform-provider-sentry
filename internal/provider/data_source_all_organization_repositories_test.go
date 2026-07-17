package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mzglinski/terraform-provider-sentry/internal/must"
	"github.com/mzglinski/terraform-provider-sentry/internal/acctest"
)

func TestAccAllOrganizationRepositoriesDataSource_GitHub(t *testing.T) {
	rn := "data.sentry_all_organization_repositories.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestGitHubInstallationId == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_INSTALLATION_ID environment variable")
			}
			if acctest.TestGitHubRepositoryIdentifier == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_REPOSITORY_IDENTIFIER environment variable")
			}

			must.Do(testAccOrganizationRepositoryResourcePreCheck())
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAllOrganizationRepositoriesConfig(testAccOrganizationRepositoryResourceConfigData{
					IntegrationType: "github",
					IntegrationId:   acctest.TestGitHubInstallationId,
					Identifier:      acctest.TestGitHubRepositoryIdentifier,
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("repositories"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"id":               knownvalue.StringRegexp(regexp.MustCompile(`^\d+$`)),
							"integration_type": knownvalue.StringExact("github"),
							"integration_id":   knownvalue.StringExact(acctest.TestGitHubInstallationId),
							"identifier":       knownvalue.StringExact(acctest.TestGitHubRepositoryIdentifier),
							"import_id": knownvalue.StringRegexp(regexp.MustCompile(
								fmt.Sprintf(`^%s/github/%s/\d+$`, acctest.TestOrganization, regexp.QuoteMeta(acctest.TestGitHubInstallationId)),
							)),
						}),
					})),
				},
			},
		},
	})
}

func testAccAllOrganizationRepositoriesConfig(data testAccOrganizationRepositoryResourceConfigData) string {
	return testAccOrganizationRepositoryResourceConfig(data) + `

data "sentry_all_organization_repositories" "test" {
	organization = data.sentry_organization.test.slug

	depends_on = [sentry_organization_repository.test]
}
`
}
