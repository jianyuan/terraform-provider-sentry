package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func testAccCheckSentryTeamExists(n string, teamID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		org := rs.Primary.Attributes["organization"]
		teamSlug := rs.Primary.ID
		ctx := context.Background()
		gotTeam, _, err := acctest.SharedClient.Teams.Get(ctx, org, teamSlug)
		if err != nil {
			return err
		}
		*teamID = sentry.StringValue(gotTeam.ID)
		return nil
	}
}

func testAccSentryTeamConfig(teamName string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}
	`, teamName)
}

func TestValidatePlatform(t *testing.T) {
	for _, tc := range []string{
		"javascript-react",
		"other",
		"python-aiohttp",
		"python",
		"react-native",
	} {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			diag := validatePlatform(tc, nil)
			if diag.HasError() {
				t.Errorf("platform should be valid: %v", tc)
			}
		})
	}
}
