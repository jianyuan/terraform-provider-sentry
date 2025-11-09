package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func init() {
	resource.AddTestSweepers("sentry_integration_opsgenie", &resource.Sweeper{
		Name: "sentry_integration_opsgenie",
		F: func(r string) error {
			sweepIntegration := func(ctx context.Context, integrationId string) error {
				httpResp, err := acctest.SharedApiClient.GetOrganizationIntegrationWithResponse(ctx, acctest.TestOrganization, integrationId)
				if err != nil {
					return err
				} else if httpResp.StatusCode() == http.StatusNotFound {
					return nil
				} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Failed to read Opsgenie integration: %s", httpResp.Status())
				}

				integration := *httpResp.JSON200

				specificIntegration, err := integration.AsOrganizationIntegrationOpsgenie()
				if err != nil {
					return err
				}

				filtered := specificIntegration.ConfigData.TeamTable[:0]
				for _, item := range specificIntegration.ConfigData.TeamTable {
					if strings.HasPrefix(item.Team, "tf-opsgenie") {
						log.Printf("[INFO] Deleting Opsgenie team %q", item.Team)
					} else {
						filtered = append(filtered, item)
					}
				}

				specificIntegration.ConfigData.TeamTable = filtered

				configDataJSON, err := json.Marshal(specificIntegration.ConfigData)
				if err != nil {
					return err
				}

				updateHttpResp, err := acctest.SharedApiClient.UpdateOrganizationIntegrationWithBodyWithResponse(
					ctx,
					acctest.TestOrganization,
					integrationId,
					"application/json",
					bytes.NewReader(configDataJSON),
				)
				if err != nil {
					return err
				} else if updateHttpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("[ERROR] Failed to update Opsgenie integration: %s", updateHttpResp.Status())
				}

				return nil
			}

			ctx := context.Background()

			params := &apiclient.ListOrganizationIntegrationsParams{
				ProviderKey: ptr.Ptr("opsgenie"),
			}

			for {
				httpResp, err := acctest.SharedApiClient.ListOrganizationIntegrationsWithResponse(ctx, acctest.TestOrganization, params)
				if err != nil {
					return err
				} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
					return fmt.Errorf("[ERROR] Failed to list Opsgenie integrations: %s", httpResp.Status())
				}

				for _, integration := range *httpResp.JSON200 {
					err := sweepIntegration(ctx, integration.Id)
					if err != nil {
						return err
					}
				}

				if httpResp.StatusCode() != http.StatusOK {
					return fmt.Errorf("[ERROR] Failed to list Opsgenie integrations: %s", httpResp.Status())
				}

				params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
				if params.Cursor == nil {
					break
				}
			}

			return nil
		},
	})
}

func TestAccIntegrationOpsgenieResource(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-opsgenie")
	rn := "sentry_integration_opsgenie.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)

			if acctest.TestOpsgenieOrganization == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_OPSGENIE_ORGANIZATION environment variable")
			}
			if acctest.TestOpsgenieIntegrationKey == "" {
				t.Skip("Skipping test due to missing SENTRY_TEST_OPSGENIE_INTEGRATION_KEY environment variable")
			}
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "sentry_integration_opsgenie" {
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
					return fmt.Errorf("failed to read Opsgenie integration: %s", httpResp.Status())
				}

				integration := *httpResp.JSON200

				specificIntegration, err := integration.AsOrganizationIntegrationOpsgenie()
				if err != nil {
					return err
				}

				for _, i := range specificIntegration.ConfigData.TeamTable {
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
					return tfutils.BuildThreePartId(organization, integrationId, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIntegrationOpsgenieResourceConfig(teamName string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
data "sentry_organization_integration" "opsgenie" {
	organization = data.sentry_organization.test.slug
	provider_key = "opsgenie"
	name         = "%[1]s"
}

resource "sentry_integration_opsgenie" "test" {
	organization    = data.sentry_organization.test.slug
	integration_id  = data.sentry_organization_integration.opsgenie.id
	team            = "%[2]s"
	integration_key = "%[3]s"
}
`, acctest.TestOpsgenieOrganization, teamName, acctest.TestOpsgenieIntegrationKey)
}
