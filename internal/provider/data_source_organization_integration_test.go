package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccOrganizationIntegrationDataSource(t *testing.T) {
	rn := "data.sentry_organization_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationIntegrationDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("provider_key"), knownvalue.StringExact("github")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("jianyuan")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("raw_json"), knownvalue.StringFunc(func(v string) error {
						var data struct {
							Id       string `json:"id"`
							Name     string `json:"name"`
							Provider struct {
								Key string `json:"key"`
							} `json:"provider"`
						}
						if err := json.Unmarshal([]byte(v), &data); err != nil {
							return err
						}
						if data.Id == "" {
							return fmt.Errorf("expected id to be non-empty")
						}
						if data.Name != "jianyuan" {
							return fmt.Errorf("expected name to be jianyuan, got %s", data.Name)
						}
						if data.Provider.Key != "github" {
							return fmt.Errorf("expected provider key to be github, got %s", data.Provider.Key)
						}
						return nil
					})),
				},
			},
		},
	})
}

var testAccOrganizationIntegrationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization_integration" "test" {
	organization = "%[1]s"
	provider_key = "github"
	name         = "jianyuan"
}
`, acctest.TestOrganization)
