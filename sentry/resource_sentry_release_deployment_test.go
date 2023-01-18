package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryReleaseDeployment_basic(t *testing.T) {
	rn := "sentry_organization_member.john_doe"
	environment := "test"

	check := func(role string) resource.TestCheckFunc {
		var deploy sentry.ReleaseDeployment
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryReleaseDeploymentExists(rn, &deploy),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "version", "0.1.0"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryReleaseDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryReleaseDeploymentConfig("0.1.0", environment),
				Check:  check("member"),
			},
			{
				Config: testAccSentryReleaseDeploymentConfig("0.2.0", environment),
				Check:  check("manager"),
			},
		},
	})
}

func testAccCheckSentryReleaseDeploymentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_release_deployment" {
			continue
		}

		ctx := context.Background()
		deploy, resp, err := client.ReleaseDeployments.Get(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.Attributes["version"],
			rs.Primary.ID,
		)
		if err == nil {
			if deploy != nil {
				return errors.New("release deployment still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryReleaseDeploymentExists(n string, deploy *sentry.ReleaseDeployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no deploy ID is set")
		}

		org, version, id, err := splitSentryReleaseDeploymentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotDeployment, _, err := client.ReleaseDeployments.Get(ctx, org, version, id)

		if err != nil {
			return err
		}
		*deploy = *gotDeployment
		return nil
	}
}

func testAccSentryReleaseDeploymentConfig(version, environment string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_release_deployment" "my-release" {
	organization = data.sentry_organization.test.id
	version      = "%[1]s"
	environment  = "%[2]s"
}
	`, version, environment)
}
