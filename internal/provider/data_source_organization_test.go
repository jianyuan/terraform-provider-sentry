package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func testAccCheckOrganizationExists(ctx context.Context, n string, v *sentry.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		org, _, err := acctest.SharedClient.Organizations.Get(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		*v = *org

		return nil
	}
}

func TestAccOrganizationDataSource_MigrateFromPluginSDK(t *testing.T) {
	ctx := context.Background()

	var v sentry.Organization
	dsn := "data.sentry_organization.test"

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
				Config: testAccOrganizationDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrganizationExists(ctx, dsn, &v),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				Config:                   testAccOrganizationDataSourceConfig,
				PlanOnly:                 true,
			},
		},
	})
}

func TestAccOrganizationDataSource(t *testing.T) {
	ctx := context.Background()

	var v sentry.Organization
	dsn := "data.sentry_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrganizationExists(ctx, dsn, &v),
					resource.TestCheckResourceAttr(dsn, "id", acctest.TestOrganization),
					resource.TestCheckResourceAttr(dsn, "slug", acctest.TestOrganization),
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(dsn, "internal_id", sentry.StringValue(v.ID))(s)
					},
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(dsn, "name", sentry.StringValue(v.Name))(s)
					},
				),
			},
		},
	})
}

var testAccOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
  slug = "%s"
}
`, acctest.TestOrganization)
