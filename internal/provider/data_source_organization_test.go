package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccOrganizationDataSource(t *testing.T) {
	rn := "data.sentry_organization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccOrganizationDataSource_UpgradeFromVersion(t *testing.T) {
	rn := "data.sentry_organization.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.12.3",
					},
				},
				Config:            testAccOrganizationDataSourceConfig,
				ConfigStateChecks: checks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccOrganizationDataSourceConfig,
				ConfigStateChecks:        checks,
			},
		},
	})
}

var testAccOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%s"
}
`, acctest.TestOrganization)
