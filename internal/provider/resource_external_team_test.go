package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

func TestAccExternalTeamResource_basic(t *testing.T) {
	rn := "sentry_external_team.test"
	teamName := acctest.RandomWithPrefix("tf-team")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			if acctest.TestGitHubInstallationId == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_INSTALLATION_ID environment variable")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExternalTeamResourceConfig(teamName, "@tf-test/"+teamName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("team"), knownvalue.StringExact(teamName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("provider_name"), knownvalue.StringExact("github")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_name"), knownvalue.StringExact("@tf-test/"+teamName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.StringExact(acctest.TestGitHubInstallationId)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_id"), knownvalue.Null()),
				},
			},
			{
				Config: testAccExternalTeamResourceConfig(teamName, "@tf-test/"+teamName+"-renamed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_name"), knownvalue.StringExact("@tf-test/"+teamName+"-renamed")),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: resourceid.ImportState3PartIDFunc(rn, "organization", "team", "id"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccExternalTeamResourceConfig(teamName, externalName string) string {
	return fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = %[1]q
	name         = %[2]q
	slug         = %[2]q
}

resource "sentry_external_team" "test" {
	organization   = sentry_team.test.organization
	team           = sentry_team.test.slug
	provider_name  = "github"
	external_name  = %[3]q
	integration_id = %[4]q
}
`, acctest.TestOrganization, teamName, externalName, acctest.TestGitHubInstallationId)
}
