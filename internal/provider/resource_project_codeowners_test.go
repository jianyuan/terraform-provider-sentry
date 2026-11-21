package provider

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

// testAccProjectCodeOwnersCodeMapping creates a repository and a code mapping on the shared
// project, because the provider cannot manage code mappings. Both are removed when the test ends.
func testAccProjectCodeOwnersCodeMapping(t *testing.T) apiclient.OrganizationCodeMapping {
	t.Helper()
	ctx := context.Background()

	// Remove a repository left behind by an earlier run.
	// shortcut: reads only the first page of repositories, paginate if the test integration has many.
	listHttpResp, err := acctest.SharedApiClient.ListOrganizationRepositoriesWithResponse(ctx, acctest.TestOrganization, &apiclient.ListOrganizationRepositoriesParams{
		IntegrationId: &acctest.TestGitHubInstallationId,
		Status:        new("active"),
	})
	if err != nil {
		t.Fatalf("failed to list repositories: %s", err)
	} else if listHttpResp.StatusCode() != http.StatusOK || listHttpResp.JSON200 == nil {
		t.Fatalf("failed to list repositories: status code=%d, body=%s", listHttpResp.StatusCode(), listHttpResp.Body)
	}
	for _, repo := range *listHttpResp.JSON200 {
		if repo.Name == acctest.TestGitHubRepositoryIdentifier {
			_, _ = acctest.SharedApiClient.DeleteOrganizationRepositoryWithResponse(ctx, acctest.TestOrganization, repo.Id)
		}
	}

	repoHttpResp, err := acctest.SharedApiClient.CreateOrganizationRepositoryWithResponse(ctx, acctest.TestOrganization, apiclient.CreateOrganizationRepositoryJSONRequestBody{
		Provider:     "integrations:github",
		Installation: acctest.TestGitHubInstallationId,
		Identifier:   acctest.TestGitHubRepositoryIdentifier,
	})
	if err != nil {
		t.Fatalf("failed to create repository: %s", err)
	} else if repoHttpResp.StatusCode() != http.StatusCreated || repoHttpResp.JSON201 == nil {
		t.Fatalf("failed to create repository: status code=%d, body=%s", repoHttpResp.StatusCode(), repoHttpResp.Body)
	}
	repo := *repoHttpResp.JSON201
	t.Cleanup(func() {
		_, _ = acctest.SharedApiClient.DeleteOrganizationRepositoryWithResponse(ctx, acctest.TestOrganization, repo.Id)
	})

	codeMappingHttpResp, err := acctest.SharedApiClient.CreateOrganizationCodeMappingWithResponse(ctx, acctest.TestOrganization, apiclient.CreateOrganizationCodeMappingJSONRequestBody{
		DefaultBranch: "main",
		StackRoot:     "app/",
		SourceRoot:    "src/app/",
		RepositoryId:  repo.Id,
		IntegrationId: acctest.TestGitHubInstallationId,
		ProjectId:     acctest.TestProject.Id,
	})
	if err != nil {
		t.Fatalf("failed to create code mapping: %s", err)
	} else if codeMappingHttpResp.StatusCode() != http.StatusCreated || codeMappingHttpResp.JSON201 == nil {
		t.Fatalf("failed to create code mapping: status code=%d, body=%s", codeMappingHttpResp.StatusCode(), codeMappingHttpResp.Body)
	}
	codeMapping := *codeMappingHttpResp.JSON201
	t.Cleanup(func() {
		_, _ = acctest.SharedApiClient.DeleteOrganizationCodeMappingWithResponse(ctx, acctest.TestOrganization, codeMapping.Id)
	})

	return codeMapping
}

func TestAccProjectCodeOwnersResource_basic(t *testing.T) {
	rn := "sentry_project_codeowners.test"

	codeMappingId := testAccProjectCodeOwnersCodeMapping(t).Id

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			if acctest.TestGitHubInstallationId == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_INSTALLATION_ID environment variable")
			}
			if acctest.TestGitHubRepositoryIdentifier == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_REPOSITORY_IDENTIFIER environment variable")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectCodeOwnersResourceConfig(codeMappingId, "*.go admin@example.com\n"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(acctest.TestProject.Slug)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("code_mapping_id"), knownvalue.StringExact(codeMappingId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("raw"), knownvalue.StringExact("*.go admin@example.com\n")),
				},
			},
			{
				Config: testAccProjectCodeOwnersResourceConfig(codeMappingId, "*.go admin@example.com\n*.md docs@example.com\n"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("raw"), knownvalue.StringExact("*.go admin@example.com\n*.md docs@example.com\n")),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: resourceid.ImportState3PartIDFunc(rn, "organization", "project", "id"),
				ImportStateVerify: true,
				// Sentry trims raw, so the imported value lacks the trailing newline from the config.
				ImportStateVerifyIgnore: []string{"raw"},
			},
		},
	})
}

func testAccProjectCodeOwnersResourceConfig(codeMappingId string, raw string) string {
	return fmt.Sprintf(`
resource "sentry_project_codeowners" "test" {
	organization    = "%[1]s"
	project         = "%[2]s"
	code_mapping_id = "%[3]s"
	raw             = %[4]q
}
`, acctest.TestOrganization, acctest.TestProject.Slug, codeMappingId, raw)
}
