package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

// testAccOrganizationAuthTokenManagerTokenEnv names the environment variable that holds an OAuth
// application token with the `org:tokens`, `org:read`, `org:write` and `org:ci` scopes. The token
// used by the other acceptance tests can't manage org auth tokens, so these tests are opt-in.
const testAccOrganizationAuthTokenManagerTokenEnv = "SENTRY_TEST_ORG_AUTH_TOKEN_MANAGER_TOKEN"

func testAccOrganizationAuthTokenPreCheck(t *testing.T) {
	t.Helper()
	acctest.PreCheck(t)
	if os.Getenv(testAccOrganizationAuthTokenManagerTokenEnv) == "" {
		t.Skipf("Skipping test due to missing %s environment variable", testAccOrganizationAuthTokenManagerTokenEnv)
	}
}

func testAccOrganizationAuthTokenResourceConfig(managerToken, name string) string {
	return fmt.Sprintf(`
provider "sentry" {
	token = %[1]q
}

resource "sentry_organization_auth_token" "test" {
	organization = %[2]q
	name         = %[3]q
}
`, managerToken, acctest.TestOrganization, name)
}

func TestAccOrganizationAuthTokenResource_basic(t *testing.T) {
	rn := "sentry_organization_auth_token.test"
	name := acctest.RandomWithPrefix("tf-org-auth-token")
	managerToken := os.Getenv(testAccOrganizationAuthTokenManagerTokenEnv)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccOrganizationAuthTokenPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_4_6),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationAuthTokenResourceConfig(managerToken, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scopes"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("org:ci")})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token"), knownvalue.NotNull()),
					statecheck.ExpectSensitiveValue(rn, tfjsonpath.New("token")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token_last_characters"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("date_created"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccOrganizationAuthTokenResourceConfig(managerToken, name+"-renamed"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(name+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("token"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", rn)
					}
					return acctest.TestOrganization + "/" + rs.Primary.ID, nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}
