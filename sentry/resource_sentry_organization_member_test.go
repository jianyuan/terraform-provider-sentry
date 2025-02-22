package sentry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

func TestAccSentryOrganizationMember_basic(t *testing.T) {
	memberEmail := acctest.RandomWithPrefix("tf-team") + "@example.com"
	rn := "sentry_organization_member.john_doe"

	check := func(role string) resource.TestCheckFunc {
		var member apiclient.OrganizationMemberWithRoles
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryOrganizationMemberExists(rn, &member),
			resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
			resource.TestCheckResourceAttr(rn, "email", memberEmail),
			resource.TestCheckResourceAttr(rn, "role", role),
			resource.TestCheckResourceAttrSet(rn, "internal_id"),
			resource.TestCheckResourceAttrSet(rn, "pending"),
			resource.TestCheckResourceAttrSet(rn, "expired"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSentryOrganizationMemberDestroy,
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization_member" {
			continue
		}

		ctx := context.Background()
		httpResp, err := acctest.SharedApiClient.GetOrganizationMemberWithResponse(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		} else if httpResp.StatusCode() != http.StatusNotFound {
			return errors.New("organization member still exists")
		}
		return nil
	}
	return nil
}

func testAccCheckSentryOrganizationMemberExists(n string, member *apiclient.OrganizationMemberWithRoles) resource.TestCheckFunc {
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
		ctx := context.Background()
		httpResp, err := acctest.SharedApiClient.GetOrganizationMemberWithResponse(ctx, org, id)
		if err != nil {
			return err
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			return fmt.Errorf("organization member not found")
		}
		*member = *httpResp.JSON200
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
