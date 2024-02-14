package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryOrganizationMember_basic(t *testing.T) {
	memberEmail := acctest.RandomWithPrefix("tf-team") + "@example.com"
	rn := "sentry_organization_member.john_doe"

	check := func(role string) resource.TestCheckFunc {
		var member sentry.OrganizationMember
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryOrganizationMemberExists(rn, &member),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "email", memberEmail),
			resource.TestCheckResourceAttr(rn, "role", role),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrSet(rn, "pending"),
			resource.TestCheckResourceAttrSet(rn, "expired"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryOrganizationMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryOrganizationMemberConfig(memberEmail, "member"),
				Check:  check("member"),
			},
			{
				Config: testAccSentryOrganizationMemberConfig(memberEmail, "manager"),
				Check:  check("manager"),
			},
		},
	})
}

func testAccCheckSentryOrganizationMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization_member" {
			continue
		}

		ctx := context.Background()
		member, resp, err := client.OrganizationMembers.Get(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err == nil {
			if member != nil {
				return errors.New("organization member still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckSentryOrganizationMemberExists(n string, member *sentry.OrganizationMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no member ID is set")
		}

		org, id, err := splitSentryOrganizationMemberID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotMember, _, err := client.OrganizationMembers.Get(ctx, org, id)

		if err != nil {
			return err
		}
		*member = *gotMember
		return nil
	}
}

func testAccSentryOrganizationMemberConfig(email, role string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_organization_member" "john_doe" {
	organization = data.sentry_organization.test.id
	email        = "%[1]s"
	role         = "%[2]s"
}
	`, email, role)
}
