package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccOrganizationResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "sentry_organization" "test" {
						name = "tf-org"
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`The argument "agree_terms" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccOrganizationResource_basic(t *testing.T) {
	orgName := acctest.RandomWithPrefix("tf-org")
	rn := "sentry_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if os.Getenv("SENTRY_RUN_ORGANIZATION_TEST") == "" {
				// Organization creation is rate limited. Only run the test once in a while.
				t.Skip("Skipping Organization tests. Set SENTRY_RUN_ORGANIZATION_TEST=true to enable.")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfig(orgName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(orgName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("agree_terms"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccOrganizationResourceConfig(orgName + "-renamed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(orgName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("agree_terms"), knownvalue.Bool(true)),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOrganizationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization" {
			continue
		}

		ctx := context.Background()
		httpResp, err := acctest.SharedApiClient.GetOrganizationWithResponse(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if httpResp.StatusCode() == http.StatusNotFound {
			continue
		}
		if httpResp.StatusCode() == http.StatusOK {
			return fmt.Errorf("organization %q still exists", rs.Primary.ID)
		}
		return fmt.Errorf("unexpected status checking organization %q: %s", rs.Primary.ID, httpResp.Status())
	}

	return nil
}

func testAccOrganizationResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "sentry_organization" "test" {
	name        = "%[1]s"
	agree_terms = true
}
`, name)
}
