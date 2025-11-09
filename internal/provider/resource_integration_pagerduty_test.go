package provider

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func TestAccIntegrationPagerDutyResource(t *testing.T) {
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
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
		},
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
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					integrationId := rs.Primary.Attributes["integration_id"]
					id := rs.Primary.ID
					return tfutils.BuildThreePartId(organization, integrationId, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIntegrationPagerDutyResourceConfig(serviceName, integrationKey string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
data "sentry_organization_integration" "pagerduty" {
	organization = data.sentry_organization.test.slug
	provider_key = "pagerduty"
	name         = "%[1]s"
}

resource "sentry_integration_pagerduty" "test" {
	organization    = data.sentry_organization.test.slug
	integration_id  = data.sentry_organization_integration.pagerduty.id
	service         = "%[2]s"
	integration_key = "%[3]s"
}
`, acctest.TestPagerDutyOrganization, serviceName, integrationKey)
}
