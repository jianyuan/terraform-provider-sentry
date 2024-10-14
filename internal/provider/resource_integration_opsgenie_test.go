package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccIntegrationOpsgenieResource(t *testing.T) {
	if acctest.TestOpsgenieOrganization == "" {
		t.Skip("Skipping test due to missing SENTRY_TEST_OPSGENIE_ORGANIZATION environment variable")
	}
	if acctest.TestOpsgenieIntegrationKey == "" {
		t.Skip("Skipping test due to missing SENTRY_TEST_OPSGENIE_INTEGRATION_KEY environment variable")
	}

	teamName := acctest.RandomWithPrefix("tf-opsgenie-service")
	rn := "sentry_integration_opsgenie.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "sentry_integration_opsgenie" {
					continue
				}

				ctx := context.Background()
				integration, apiResp, err := acctest.SharedClient.OrganizationIntegrations.Get(
					ctx,
					rs.Primary.Attributes["organization"],
					rs.Primary.Attributes["integration_id"],
				)

				if apiResp.StatusCode == 404 {
					return nil
				}
				if err != nil {
					return err
				}
				var configData IntegrationOpsgenieConfigData
				if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
					return err
				}

				for _, i := range configData.TeamTable {
					if i.Id == rs.Primary.ID {
						return fmt.Errorf("Opsgenie service %q still exists", rs.Primary.ID)
					}
				}

				return nil
			}

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationOpsgenieResourceConfig(teamName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("team"), knownvalue.StringExact(teamName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.StringExact(acctest.TestOpsgenieIntegrationKey)),
				},
			},
			{
				Config: testAccIntegrationOpsgenieResourceConfig(teamName + "-changed"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("team"), knownvalue.StringExact(teamName+"-changed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_key"), knownvalue.StringExact(acctest.TestOpsgenieIntegrationKey)),
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
					return buildThreePartID(organization, integrationId, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIntegrationOpsgenieResourceConfig(teamName string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
data "sentry_organization_integration" "opsgenie" {
	organization = data.sentry_organization.test.id
	provider_key = "opsgenie"
	name         = "%[1]s"
}

resource "sentry_integration_opsgenie" "test" {
	organization    = data.sentry_organization.test.id
	integration_id  = data.sentry_organization_integration.opsgenie.id
	team            = "%[2]s"
	integration_key = "%[3]s"
}
`, acctest.TestOpsgenieOrganization, teamName, acctest.TestOpsgenieIntegrationKey)
}
