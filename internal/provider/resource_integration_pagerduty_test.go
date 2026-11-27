package provider

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

func TestAccIntegrationPagerDutyResource(t *testing.T) {
	serviceName := acctest.RandomWithPrefix("tf-pagerduty-service")
	integrationKey := acctest.RandomWithPrefix("tf-integration-key")
	rn := "sentry_integration_pagerduty.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestPagerDutyOrganization == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_PAGERDUTY_ORGANIZATION environment variable")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIntegrationPagerDutyResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationPagerDutyResourceConfig(serviceName, integrationKey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service"), knownvalue.StringExact(serviceName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.StringExact(integrationKey)),
				},
			},
			{
				Config: testAccIntegrationPagerDutyResourceConfig(serviceName+"-changed", integrationKey+"-changed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service"), knownvalue.StringExact(serviceName+"-changed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.StringExact(integrationKey+"-changed")),
				},
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: resourceid.ImportState3PartIDFunc(rn, "organization", "integration_id", "id"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIntegrationPagerDutyResourceConfig(serviceName, integrationKey string) string {
	return fmt.Sprintf(`
data "sentry_organization_integration" "pagerduty" {
	organization = "%[1]s"
	provider_key = "pagerduty"
	name         = "%[2]s"
}

resource "sentry_integration_pagerduty" "test" {
	organization    = data.sentry_organization_integration.pagerduty.organization
	integration_id  = data.sentry_organization_integration.pagerduty.id
	service         = "%[3]s"
	integration_key = "%[4]s"
}
`, acctest.TestOrganization, acctest.TestPagerDutyOrganization, serviceName, integrationKey)
}

func testAccCheckIntegrationPagerDutyResourceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_integration_pagerduty" {
			continue
		}

		ctx := context.Background()
		httpResp, err := acctest.SharedApiClient.GetOrganizationIntegrationWithResponse(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["integration_id"],
		)

		if err != nil {
			return err
		} else if httpResp.StatusCode() == http.StatusNotFound {
			return nil
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			return fmt.Errorf("failed to read PagerDuty integration: %s", httpResp.Status())
		}

		integration := *httpResp.JSON200

		specificIntegration, err := integration.AsOrganizationIntegrationPagerDuty()
		if err != nil {
			return err
		}

		for _, i := range specificIntegration.ConfigData.ServiceTable {
			if i.Id.String() == rs.Primary.ID {
				return fmt.Errorf("PagerDuty service %q still exists", rs.Primary.ID)
			}
		}

		return nil
	}

	return nil
}

func TestAccIntegrationPagerDutyResource_writeOnly(t *testing.T) {
	serviceName := acctest.RandomWithPrefix("tf-pagerduty-service")
	integrationKey := acctest.RandomWithPrefix("tf-integration-key")
	rn := "sentry_integration_pagerduty.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestPagerDutyOrganization == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_PAGERDUTY_ORGANIZATION environment variable")
			}
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIntegrationPagerDutyResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationPagerDutyResourceWriteOnlyConfig(serviceName, integrationKey, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("service"), knownvalue.StringExact(serviceName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key_wo_version"), knownvalue.Int64Exact(1)),
				},
			},
			// Bumping the version sends the new key.
			{
				Config: testAccIntegrationPagerDutyResourceWriteOnlyConfig(serviceName, integrationKey+"-changed", 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key_wo_version"), knownvalue.Int64Exact(2)),
				},
			},
			// Import reads integration_key from Sentry; the write-only attributes are not imported.
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       resourceid.ImportState3PartIDFunc(rn, "organization", "integration_id", "id"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"integration_key", "integration_key_wo", "integration_key_wo_version"},
			},
		},
	})
}

func testAccIntegrationPagerDutyResourceWriteOnlyConfig(serviceName, integrationKey string, integrationKeyVersion int) string {
	return fmt.Sprintf(`
data "sentry_organization_integration" "pagerduty" {
	organization = "%[1]s"
	provider_key = "pagerduty"
	name         = "%[2]s"
}

resource "sentry_integration_pagerduty" "test" {
	organization               = data.sentry_organization_integration.pagerduty.organization
	integration_id             = data.sentry_organization_integration.pagerduty.id
	service                    = "%[3]s"
	integration_key_wo         = "%[4]s"
	integration_key_wo_version = %[5]d
}
`, acctest.TestOrganization, acctest.TestPagerDutyOrganization, serviceName, integrationKey, integrationKeyVersion)
}
