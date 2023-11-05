package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func testAccCheckTeamExists(ctx context.Context, n string, v *sentry.Team) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		team, _, err := acctest.SharedClient.Teams.Get(ctx, rs.Primary.Attributes["organization"], rs.Primary.ID)
		if err != nil {
			return err
		}
		*v = *team
		return nil
	}
}

func testAccCheckTeamDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_team" {
			continue
		}

		ctx := context.Background()
		team, resp, err := acctest.SharedClient.Teams.Get(
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

func testAccTeamImportStateIdFunc(n string) resource.ImportStateIdFunc {
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

func TestAccTeamResource_MigrateFromPluginSDK(t *testing.T) {
	ctx := context.Background()

	var v sentry.Team
	teamName := acctest.RandomWithPrefix("tf-team")
	resourceName := "sentry_team.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.11.2",
					},
				},
				Config: testAccTeamConfig(teamName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamExists(ctx, resourceName, &v),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				Config:                   testAccTeamConfig(teamName),
				PlanOnly:                 true,
			},
		},
	})
}

func TestAccTeamResource(t *testing.T) {
	ctx := context.Background()

	var v sentry.Team
	teamName := acctest.RandomWithPrefix("tf-team")
	resourceName := "sentry_team.test"

	check := func(teamName string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			testAccCheckTeamExists(ctx, resourceName, &v),
			resource.TestCheckResourceAttr(resourceName, "id", teamName),
			resource.TestCheckResourceAttrPair(resourceName, "organization", "data.sentry_organization.test", "id"),
			resource.TestCheckResourceAttr(resourceName, "name", teamName),
			resource.TestCheckResourceAttr(resourceName, "slug", teamName),
			func(s *terraform.State) error {
				return resource.TestCheckResourceAttr(resourceName, "internal_id", sentry.StringValue(v.ID))(s)
			},
			func(s *terraform.State) error {
				return resource.TestCheckResourceAttr(resourceName, "has_access", strconv.FormatBool(sentry.BoolValue(v.HasAccess)))(s)
			},
			func(s *terraform.State) error {
				return resource.TestCheckResourceAttr(resourceName, "is_pending", strconv.FormatBool(sentry.BoolValue(v.IsPending)))(s)
			},
			func(s *terraform.State) error {
				return resource.TestCheckResourceAttr(resourceName, "is_member", strconv.FormatBool(sentry.BoolValue(v.IsMember)))(s)
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig(teamName),
				Check:  check(teamName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTeamImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccTeamConfig(teamName + "-renamed"),
				Check:  check(teamName + "-renamed"),
			},
		},
	})
}

func testAccTeamConfig(teamName string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
  organization = data.sentry_organization.test.id
  name         = "%[1]s"
  slug         = "%[1]s"
}
`, teamName)
}
