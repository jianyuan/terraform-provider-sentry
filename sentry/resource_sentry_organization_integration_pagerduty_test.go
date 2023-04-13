package sentry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryOrganizationIntegrationPagerduty_basic(t *testing.T) {
	integrationId := os.Getenv("SENTRY_TEST_PAGERDUTY_INTEGRATION_ID")
	if integrationId == "" {
		t.Skip("Skipping because SENTRY_TEST_PAGERDUTY_INTEGRATION_ID is not set")
	}

	serviceName := acctest.RandomWithPrefix("tf-pagerduty-service")
	integrationKey := acctest.RandomWithPrefix("tf-pagerduty-integration-key")
	rn := "sentry_organization_integration_pagerduty.test"

	var serviceID string

	check := func(serviceName, integrationKey string) resource.TestCheckFunc {
		fs := resource.ComposeTestCheckFunc(
			testAccCheckSentryOrganizationIntegrationPagerdutyExists(rn, &serviceID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "integration_id", integrationId),
			resource.TestCheckResourceAttr(rn, "service_name", serviceName),
			resource.TestCheckResourceAttr(rn, "integration_key", integrationKey),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &serviceID),
		)
		return fs
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryOrganizationIntegrationPagerdutyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryOrganizationIntegrationPagerdutyConfig(integrationId, serviceName, integrationKey),
				Check:  check(serviceName, integrationKey),
			},
			{
				Config: testAccSentryOrganizationIntegrationPagerdutyConfig(integrationId, serviceName+"-renamed", integrationKey),
				Check:  check(serviceName+"-renamed", integrationKey),
			},
			{
				Config: testAccSentryOrganizationIntegrationPagerdutyConfig(integrationId, serviceName+"-renamed", integrationKey+"-renamed"),
				Check:  check(serviceName+"-renamed", integrationKey+"-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryOrganizationIntegrationPagerdutyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization_integration_pagerduty" {
			continue
		}

		ctx := context.Background()
		//proj, resp, err := client.Projects.Get(ctx, testOrganization, rs.Primary.ID)
		org, integrationId, id, err := splitThreePartID(rs.Primary.ID, "organization-slug", "integration-id", "service-id")

		orgIntegration, resp, err := client.OrganizationIntegrations.Get(ctx, org, integrationId)
		configData := *orgIntegration.ConfigData
		var foundServiceRow map[string]interface{}
		if serviceTable, found := configData["service_table"]; found {
			for _, row := range serviceTable.([]interface{}) {
				serviceRow := row.(map[string]interface{})
				if string(serviceRow["id"].(json.Number)) == id {
					foundServiceRow = serviceRow
					break
				}
			}
		} else {
			return errors.New("unable to find PagerDuty service_table information")
		}

		if err == nil {
			if foundServiceRow != nil {
				return errors.New("pagerduty service id still exists")
			}
		}
		if resp.StatusCode != 403 && resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryOrganizationIntegrationPagerdutyExists(n string, serviceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no pagerduty service ID is set")
		}

		org, integrationId, id, err := splitThreePartID(rs.Primary.ID, "organization-slug", "integration-id", "service-id")
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		orgIntegration, _, err := client.OrganizationIntegrations.Get(ctx, org, integrationId)
		configData := *orgIntegration.ConfigData
		var foundServiceRow map[string]interface{}
		if serviceTable, found := configData["service_table"]; found {
			for _, row := range serviceTable.([]interface{}) {
				serviceRow := row.(map[string]interface{})
				if string(serviceRow["id"].(json.Number)) == id {
					foundServiceRow = serviceRow
					break
				}
			}
		} else {
			return errors.New("unable to find PagerDuty service_table information")
		}

		if err != nil {
			return err
		}
		*serviceID = string(foundServiceRow["id"].(json.Number))
		return nil
	}
}

func testAccSentryOrganizationIntegrationPagerdutyConfig(integrationId, serviceName, integrationKey string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_organization_integration_pagerduty" "test" {
	organization    = data.sentry_organization.test.id
	integration_id  = "%s"
	service_name    = "%s"
	integration_key = "%s"
}
	`, integrationId, serviceName, integrationKey)
}
