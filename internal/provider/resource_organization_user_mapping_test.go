package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func TestAccOrganizationUserMappingResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true,
				Config: `
					resource "sentry_organization_user_mapping" "test" {
						organization      = "my-org"
						user_id           = 1
						external_provider = "github"
						external_id       = "2"
						external_name     = "@octocat"
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`The argument "integration_id" is required, but no definition was found.`),
			},
			{
				PlanOnly: true,
				Config: `
					resource "sentry_organization_user_mapping" "test" {
						organization      = "my-org"
						user_id           = 1
						integration_id    = 2
						external_provider = "not-a-provider"
						external_id       = "3"
						external_name     = "@octocat"
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Attribute external_provider value must be one of`),
			},
		},
	})
}

func TestAccOrganizationUserMappingResource_basic(t *testing.T) {
	acctest.PreCheck(t)

	if acctest.TestGitHubInstallationId == "" {
		t.Skip("Skipping test due to missing SENTRY_TEST_GITHUB_INSTALLATION_ID environment variable")
	}

	integrationID, err := strconv.ParseInt(acctest.TestGitHubInstallationId, 10, 64)
	if err != nil {
		t.Fatalf("SENTRY_TEST_GITHUB_INSTALLATION_ID must be an integer: %s", err)
	}

	userID := testAccOrganizationUserMappingFindUserID(t)
	externalID := strconv.FormatInt(time.Now().UnixNano()%1_000_000_000, 10)
	externalName := "@" + acctest.RandomWithPrefix("tf-user")
	externalNameUpdated := externalName + "-updated"
	rn := "sentry_organization_user_mapping.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckOrganizationUserMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationUserMappingResourceConfig(userID, integrationID, externalID, externalName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.Int64Exact(userID)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.Int64Exact(integrationID)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_provider"), knownvalue.StringExact("github")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_id"), knownvalue.StringExact(externalID)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_name"), knownvalue.StringExact(externalName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccOrganizationUserMappingResourceConfig(userID, integrationID, externalID, externalNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_name"), knownvalue.StringExact(externalNameUpdated)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("user_id"), knownvalue.Int64Exact(userID)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integration_id"), knownvalue.Int64Exact(integrationID)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("external_id"), knownvalue.StringExact(externalID)),
				},
			},
			{
				ResourceName: rn,
				ImportState:  true,
				// Sentry has no GET for external user mappings, so import only restores
				// organization + internal_id (and id via read).
				ImportStateVerify: false,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}
					attrs := states[0].Attributes
					if attrs["organization"] != acctest.TestOrganization {
						return fmt.Errorf("organization = %q, want %q", attrs["organization"], acctest.TestOrganization)
					}
					if attrs["internal_id"] == "" {
						return fmt.Errorf("internal_id is empty")
					}
					return nil
				},
			},
		},
	})
}

func testAccOrganizationUserMappingFindUserID(t *testing.T) int64 {
	t.Helper()

	ctx := context.Background()
	params := &sentry.ListCursorParams{}

	for {
		members, resp, err := acctest.SharedClient.OrganizationMembers.List(ctx, acctest.TestOrganization, params)
		if err != nil {
			t.Fatalf("failed to list organization members: %s", err)
		}

		for _, member := range members {
			if member == nil || member.Pending || member.User.ID == "" {
				continue
			}
			userID, err := strconv.ParseInt(member.User.ID, 10, 64)
			if err != nil {
				continue
			}
			return userID
		}

		if resp.Cursor == "" {
			break
		}
		params.Cursor = resp.Cursor
	}

	t.Skip("no active organization member with a user id found")
	return 0
}

func testAccCheckOrganizationUserMappingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_organization_user_mapping" {
			continue
		}

		organization := rs.Primary.Attributes["organization"]
		internalID := rs.Primary.Attributes["internal_id"]
		if organization == "" || internalID == "" {
			organization, internalID, _ = tfutils.SplitTwoPartId(rs.Primary.ID, "organization", "internal_id")
		}
		if organization == "" || internalID == "" {
			return fmt.Errorf("unable to determine organization/internal_id for %s", rs.Primary.ID)
		}

		// There is no GET API; a second delete should 404 once the mapping is gone.
		ctx := context.Background()
		httpResp, err := acctest.SharedApiClient.DeleteOrganizationExternalUserWithResponse(ctx, organization, internalID)
		if err != nil {
			return err
		}
		if httpResp.StatusCode() == http.StatusNotFound {
			continue
		}
		if httpResp.StatusCode() == http.StatusNoContent {
			return fmt.Errorf("organization user mapping %q still exists", rs.Primary.ID)
		}
		return fmt.Errorf("unexpected status checking organization user mapping %q: %s", rs.Primary.ID, httpResp.Status())
	}

	return nil
}

func testAccOrganizationUserMappingResourceConfig(userID, integrationID int64, externalID, externalName string) string {
	return fmt.Sprintf(`
resource "sentry_organization_user_mapping" "test" {
	organization      = "%[1]s"
	user_id           = %[2]d
	integration_id    = %[3]d
	external_provider = "github"
	external_id       = "%[4]s"
	external_name     = "%[5]s"
}
`, acctest.TestOrganization, userID, integrationID, externalID, externalName)
}
