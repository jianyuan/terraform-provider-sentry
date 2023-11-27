package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/deste-org/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryTeam_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	rn := "sentry_team.test"

	var teamID string

	check := func(teamName string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryTeamExists(rn, &teamID),
			resource.TestCheckResourceAttr(rn, "id", teamName),
			resource.TestCheckResourceAttrPair(rn, "organization", "data.sentry_organization.test", "id"),
			resource.TestCheckResourceAttr(rn, "name", teamName),
			resource.TestCheckResourceAttr(rn, "slug", teamName),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &teamID),
			resource.TestCheckResourceAttrPair(rn, "internal_id", rn, "team_id"),
			resource.TestCheckResourceAttrSet(rn, "has_access"),
			resource.TestCheckResourceAttrSet(rn, "is_pending"),
			resource.TestCheckResourceAttrSet(rn, "is_member"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryTeamConfig(teamName),
				Check:  check(teamName),
			},
			{
				Config: testAccSentryTeamConfig(teamName + "-renamed"),
				Check:  check(teamName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: testAccSentryTeamImportStateIdFunc(rn),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sentry.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_team" {
			continue
		}

		ctx := context.Background()
		team, resp, err := client.Teams.Get(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err == nil {
			if team != nil {
				return errors.New("team still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

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
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotTeam, _, err := client.Teams.Get(ctx, org, teamSlug)
		if err != nil {
			return err
		}
		*teamID = sentry.StringValue(gotTeam.ID)
		return nil
	}
}

func testAccSentryTeamImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		org := rs.Primary.Attributes["organization"]
		teamSlug := rs.Primary.ID
		return buildTwoPartID(org, teamSlug), nil
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
