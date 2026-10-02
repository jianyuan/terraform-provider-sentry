package provider

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

func TestAccOrganizationDataScrubbingResource(t *testing.T) {
	rn := "sentry_organization_data_scrubbing.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckOrganizationDataScrubbingDestroy,
		Steps: []resource.TestStep{
			{
				// Mixed case checks that Sentry stores the field names unchanged.
				Config: testAccOrganizationDataScrubbingResourceConfig(`
					data_scrubber          = true
					data_scrubber_defaults = true
					sensitive_fields       = ["email", "IBAN"]
					safe_fields            = ["order_id"]
					scrub_ip_addresses     = true
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber_defaults"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("sensitive_fields"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("email"),
						knownvalue.StringExact("IBAN"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("safe_fields"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("order_id"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scrub_ip_addresses"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccOrganizationDataScrubbingResourceConfig(""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber_defaults"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("sensitive_fields"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("safe_fields"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scrub_ip_addresses"), knownvalue.Bool(false)),
				},
			},
			{
				Config: testAccOrganizationDataScrubbingResourceConfig(`
					data_scrubber    = true
					sensitive_fields = ["phone"]
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber_defaults"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("sensitive_fields"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("phone"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("safe_fields"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scrub_ip_addresses"), knownvalue.Bool(false)),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateId:     acctest.TestOrganization,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOrganizationDataScrubbingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization_data_scrubbing" {
			continue
		}

		httpResp, err := acctest.SharedApiClient.GetOrganizationWithResponse(context.Background(), rs.Primary.Attributes["organization"])
		if err != nil {
			return err
		} else if httpResp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unexpected status code %d", httpResp.StatusCode())
		}

		var settings OrganizationDataScrubbingResourceModel
		if diags := settings.Fill(context.Background(), *httpResp.JSON200); diags.HasError() {
			return diagutils.DiagnosticsError(diags)
		}
		settings.Id = types.StringNull()

		if diff := cmp.Diff(organizationDataScrubbingDefaults(), settings, organizationDataScrubbingCmpOptions); diff != "" {
			return fmt.Errorf("data scrubbing settings of organization %q were not reset (-want +got):\n%s", rs.Primary.ID, diff)
		}
	}

	return nil
}

func testAccOrganizationDataScrubbingResourceConfig(settings string) string {
	return fmt.Sprintf(`
resource "sentry_organization_data_scrubbing" "test" {
	organization = "%[1]s"
	%[2]s
}
`, acctest.TestOrganization, settings)
}
