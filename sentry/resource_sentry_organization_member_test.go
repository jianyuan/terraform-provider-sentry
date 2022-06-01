package sentry

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/sentry"
)

func TestAccSentryOrganizationMember_basic(t *testing.T) {
	var member sentry.OrganizationMember

	testAccSentrySentryOrganizationMemberUpdateConfig := fmt.Sprintf(`
    resource "sentry_organization_member" "john_doe" {
      organization = "%s"
      email = "test2@example.com"
      role = "manager"
    }
	`, testOrganization)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSentryOrganizationMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryOrganizationMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryOrganizationMemberExists("sentry_organization_member.john_doe", &member),
					testAccCheckSentryOrganizationMemberAttributes(&member, &testAccSentryOrganizationMemberExpectedAttributes{
						Email: "test2@example.com",
						Role:  sentry.RoleMember,
					}),
				),
			},
			{
				Config: testAccSentrySentryOrganizationMemberUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSentryOrganizationMemberExists("sentry_organization_member.john_doe", &member),
					testAccCheckSentryOrganizationMemberAttributes(&member, &testAccSentryOrganizationMemberExpectedAttributes{
						Email: "test2@example.com",
						Role:  sentry.RoleManager,
					}),
				),
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

		member, resp, err := client.OrganizationMembers.Get(
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

		client := testAccProvider.Meta().(*sentry.Client)
		sentryOrganizationMember, _, err := client.OrganizationMembers.Get(
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)

		if err != nil {
			return err
		}
		*member = *sentryOrganizationMember
		return nil
	}
}

type testAccSentryOrganizationMemberExpectedAttributes struct {
	Email string
	Role  string
	Teams []string
}

func testAccCheckSentryOrganizationMemberAttributes(member *sentry.OrganizationMember, want *testAccSentryOrganizationMemberExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if member.Email != want.Email {
			return fmt.Errorf("got email %q; want %q", member.Email, want.Email)
		}

		if member.Role != want.Role {
			return fmt.Errorf("got role %q; want %q", member.Role, want.Role)
		}

		if len(member.Teams) != len(want.Teams) {
			return fmt.Errorf("got total teams %d; want %d", len(member.Teams), len(want.Teams))
		}

		return nil
	}
}

var testAccSentryOrganizationMemberConfig = fmt.Sprintf(`
  resource "sentry_organization_member" "john_doe" {
    organization = "%s"
    email = "test2@example.com"
	role = "member"
  }
`, testOrganization)
