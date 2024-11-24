package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func init() {
	resource.AddTestSweepers("sentry_organization_repository", &resource.Sweeper{
		Name: "sentry_organization_repository",
		F: func(r string) error {
			return testAccOrganizationRepositoryResourcePreCheck()
		},
	})
}

func testAccOrganizationRepositoryResourcePreCheck() error {
	ctx := context.Background()

	listParams := &sentry.ListOrganizationRepositoriesParams{
		Status: "active",
	}

	for {
		repos, resp, err := acctest.SharedClient.OrganizationRepositories.List(ctx, acctest.TestOrganization, listParams)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			var externalSlugStr string
			var externalSlugNum json.Number
			if err := json.Unmarshal([]byte(repo.ExternalSlug), &externalSlugStr); err != nil {
				if err := json.Unmarshal([]byte(repo.ExternalSlug), &externalSlugNum); err != nil {
					log.Printf("[ERROR] Failed to unmarshal ExternalSlug %q: %s", repo.ExternalSlug, err)
					continue
				}
			}

			found := false
			if acctest.TestGitHubInstallationId != "" && repo.Provider.ID == "integrations:github" && repo.IntegrationId == acctest.TestGitHubInstallationId && externalSlugStr == acctest.TestGitHubRepositoryIdentifier {
				found = true
			} else if acctest.TestGitLabInstallationId != "" && repo.Provider.ID == "integrations:gitlab" && repo.IntegrationId == acctest.TestGitLabInstallationId && externalSlugNum.String() == acctest.TestGitLabRepositoryIdentifier {
				found = true
			} else if acctest.TestVSTSInstallationId != "" && repo.Provider.ID == "integrations:vsts" && repo.IntegrationId == acctest.TestVSTSInstallationId && externalSlugStr == acctest.TestVSTSRepositoryIdentifier {
				found = true
			}

			if !found {
				continue
			}

			log.Printf("[INFO] Destroying OrganizationRepository: %s", repo.ID)

			_, _, err := acctest.SharedClient.OrganizationRepositories.Delete(ctx, acctest.TestOrganization, repo.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to destroy OrganizationRepository %q: %s", repo.ID, err)
				continue
			}

			log.Printf("[INFO] OrganizationRepository %q has been destroyed.", repo.ID)
		}

		if resp.Cursor == "" {
			break
		}
		listParams.Cursor = resp.Cursor
	}

	return nil
}

func TestAccOrganizationRepositoryResource_GitHub(t *testing.T) {
	rn := "sentry_organization_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestGitHubInstallationId == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_INSTALLATION_ID environment variable")
			}
			if acctest.TestGitHubRepositoryIdentifier == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_REPOSITORY_IDENTIFIER environment variable")
			}

			testAccOrganizationRepositoryResourcePreCheck()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationRepositoryResourceConfig(testAccOrganizationRepositoryResourceConfigData{
					IntegrationType: "github",
					IntegrationId:   acctest.TestGitHubInstallationId,
					Identifier:      acctest.TestGitHubRepositoryIdentifier,
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_type"), knownvalue.StringExact("github")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.StringExact(acctest.TestGitHubInstallationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("identifier"), knownvalue.StringExact(acctest.TestGitHubRepositoryIdentifier)),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					integrationType := rs.Primary.Attributes["integration_type"]
					integrationId := rs.Primary.Attributes["integration_id"]
					id := rs.Primary.ID
					return buildFourPartID(organization, integrationType, integrationId, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationRepositoryResource_GitLab(t *testing.T) {
	rn := "sentry_organization_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestGitLabInstallationId == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITLAB_INSTALLATION_ID environment variable")
			}
			if acctest.TestGitLabRepositoryIdentifier == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITLAB_REPOSITORY_IDENTIFIER environment variable")
			}

			testAccOrganizationRepositoryResourcePreCheck()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationRepositoryResourceConfig(testAccOrganizationRepositoryResourceConfigData{
					IntegrationType: "gitlab",
					IntegrationId:   acctest.TestGitLabInstallationId,
					Identifier:      acctest.TestGitLabRepositoryIdentifier,
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_type"), knownvalue.StringExact("gitlab")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.StringExact(acctest.TestGitLabInstallationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("identifier"), knownvalue.StringExact(acctest.TestGitLabRepositoryIdentifier)),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					integrationType := rs.Primary.Attributes["integration_type"]
					integrationId := rs.Primary.Attributes["integration_id"]
					id := rs.Primary.ID
					return buildFourPartID(organization, integrationType, integrationId, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationRepositoryResource_VSTS(t *testing.T) {
	rn := "sentry_organization_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestVSTSInstallationId == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_VSTS_INSTALLATION_ID environment variable")
			}
			if acctest.TestVSTSRepositoryIdentifier == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_VSTS_REPOSITORY_IDENTIFIER environment variable")
			}

			testAccOrganizationRepositoryResourcePreCheck()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationRepositoryResourceConfig(testAccOrganizationRepositoryResourceConfigData{
					IntegrationType: "vsts",
					IntegrationId:   acctest.TestVSTSInstallationId,
					Identifier:      acctest.TestVSTSRepositoryIdentifier,
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_type"), knownvalue.StringExact("vsts")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.StringExact(acctest.TestVSTSInstallationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("identifier"), knownvalue.StringExact(acctest.TestVSTSRepositoryIdentifier)),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					integrationType := rs.Primary.Attributes["integration_type"]
					integrationId := rs.Primary.Attributes["integration_id"]
					id := rs.Primary.ID
					return buildFourPartID(organization, integrationType, integrationId, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

var testAccOrganizationRepositoryResourceConfigTemplate = template.Must(template.New("config").Parse(`
resource "sentry_organization_repository" "test" {
	organization     = data.sentry_organization.test.slug
	integration_type = "{{ .IntegrationType }}"
	integration_id   = "{{ .IntegrationId }}"
	identifier	     = "{{ .Identifier }}"
}
`))

type testAccOrganizationRepositoryResourceConfigData struct {
	IntegrationType string
	IntegrationId   string
	Identifier      string
}

func testAccOrganizationRepositoryResourceConfig(data testAccOrganizationRepositoryResourceConfigData) string {
	var builder strings.Builder

	must.Get(builder.WriteString(testAccOrganizationDataSourceConfig))
	must.Do(testAccOrganizationRepositoryResourceConfigTemplate.Execute(&builder, data))

	return builder.String()
}
