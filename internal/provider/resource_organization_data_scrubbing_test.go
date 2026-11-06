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

	const relayPiiConfig = `{"applications":{"$message":["0"]},"rules":{"0":{"redaction":{"method":"mask"},"type":"creditcard"}}}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckOrganizationDataScrubbingDestroy,
		Steps: []resource.TestStep{
			{
				// Mixed case checks that Sentry stores field names unchanged; the heredoc checks that the stored layout causes no diff.
				Config: testAccOrganizationDataScrubbingResourceConfig(`
					data_scrubber          = true
					data_scrubber_defaults = true
					scrub_ip_addresses     = true
					sensitive_fields       = ["email", "IBAN"]
					safe_fields            = ["order_id"]
					relay_pii_config       = <<-EOT
						{
							"rules": {"0": {"type": "creditcard", "redaction": {"method": "mask"}}},
							"applications": {"$message": ["0"]}
						}
					EOT
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber_defaults"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scrub_ip_addresses"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("sensitive_fields"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("email"),
						knownvalue.StringExact("IBAN"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("safe_fields"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("order_id"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("relay_pii_config"), knownvalue.NotNull()),
				},
				Check: testAccCheckOrganizationDataScrubbing(func(settings OrganizationDataScrubbingResourceModel) error {
					if got := settings.RelayPiiConfig.ValueString(); got != relayPiiConfig {
						return fmt.Errorf("unexpected relayPiiConfig in Sentry: %s", got)
					}
					return nil
				}),
			},
			{
				// Omitted attributes reset to the Sentry defaults, including the rules.
				Config: testAccOrganizationDataScrubbingResourceConfig(""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber_defaults"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("scrub_ip_addresses"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("sensitive_fields"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("safe_fields"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("relay_pii_config"), knownvalue.Null()),
				},
			},
			{
				Config: testAccOrganizationDataScrubbingResourceConfig(`
					data_scrubber    = true
					sensitive_fields = ["phone"]
					relay_pii_config = jsonencode({
						rules        = { "0" = { type = "creditcard", redaction = { method = "mask" } } }
						applications = { "$message" = ["0"] }
					})
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("data_scrubber_defaults"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("sensitive_fields"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("phone"),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("safe_fields"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("relay_pii_config"), knownvalue.StringExact(relayPiiConfig)),
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

func testAccCheckOrganizationDataScrubbing(check func(OrganizationDataScrubbingResourceModel) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		settings, err := testAccReadOrganizationDataScrubbing(acctest.TestOrganization)
		if err != nil {
			return err
		}
		return check(settings)
	}
}

func testAccCheckOrganizationDataScrubbingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization_data_scrubbing" {
			continue
		}

		settings, err := testAccReadOrganizationDataScrubbing(rs.Primary.Attributes["organization"])
		if err != nil {
			return err
		}
		settings.Id = types.StringNull()

		if diff := cmp.Diff(organizationDataScrubbingDefaults(), settings, dataScrubbingCmpOptions); diff != "" {
			return fmt.Errorf("data scrubbing settings of organization %q were not reset (-want +got):\n%s", rs.Primary.ID, diff)
		}
	}

	return nil
}

func testAccReadOrganizationDataScrubbing(organization string) (OrganizationDataScrubbingResourceModel, error) {
	var settings OrganizationDataScrubbingResourceModel

	httpResp, err := acctest.SharedApiClient.GetOrganizationWithResponse(context.Background(), organization)
	if err != nil {
		return settings, err
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		return settings, fmt.Errorf("unexpected status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
	}

	if diags := settings.Fill(context.Background(), *httpResp.JSON200); diags.HasError() {
		return settings, diagutils.DiagnosticsError(diags)
	}

	return settings, nil
}

func testAccOrganizationDataScrubbingResourceConfig(settings string) string {
	return fmt.Sprintf(`
resource "sentry_organization_data_scrubbing" "test" {
	organization = "%[1]s"
	%[2]s
}
`, acctest.TestOrganization, settings)
}
