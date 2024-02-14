package sentry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryOrganization_basic(t *testing.T) {
	if os.Getenv("SENTRY_RUN_ORGANIZATION_TEST") == "" {
		// Organization creation is rate limited. Only run the test once in a while.
		t.Skip("Skipping Organization tests. Set SENTRY_RUN_ORGANIZATION_TEST=true to enable.")
	}

	orgName := acctest.RandomWithPrefix("tf-org")
	rn := "sentry_organization.test_organization"

	check := func(orgName string) resource.TestCheckFunc {
		var organization sentry.Organization

		return resource.ComposeTestCheckFunc(
			testAccCheckSentryOrganizationExists(rn, &organization),
			resource.TestCheckResourceAttr(rn, "name", orgName),
			resource.TestCheckResourceAttrSet(rn, "slug"),
			resource.TestCheckResourceAttrWith(rn, "internal_id", func(v string) error {
				want := sentry.StringValue(organization.ID)
				if v != want {
					return fmt.Errorf("got organization ID %s; want %s", v, want)
				}
				return nil
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryOrganizationConfig(orgName),
				Check:  check(orgName),
			},
			{
				Config: testAccSentryOrganizationConfig(orgName + "-renamed"),
				Check:  check(orgName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryOrganizationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization" {
			continue
		}

		ctx := context.Background()
		organization, resp, err := client.Organizations.Get(ctx, rs.Primary.ID)
		if err == nil {
			if organization != nil && *organization.Status.ID == "active" {
				return errors.New("organization still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckSentryOrganizationExists(n string, organization *sentry.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		org := rs.Primary.ID
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotOrganization, _, err := client.Organizations.Get(ctx, org)
		if err != nil {
			return err
		}
		*organization = *gotOrganization
		return nil
	}
}

func testAccSentryOrganizationConfig(orgName string) string {
	return fmt.Sprintf(`
resource "sentry_organization" "test_organization" {
	name = "%[1]s"

	agree_terms = true
}
	`, orgName)
}
