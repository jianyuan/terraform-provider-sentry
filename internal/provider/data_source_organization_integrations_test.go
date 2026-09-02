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

func TestAccOrganizationIntegrationsDataSource(t *testing.T) {
	rn := "data.sentry_organization_integrations.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationIntegrationsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("provider_key"), knownvalue.StringExact("github")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("integrations"), knownvalue.SetPartial([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id":           knownvalue.NotNull(),
							"name":         knownvalue.NotNull(),
							"provider_key": knownvalue.StringExact("github"),
							"raw_json": knownvalue.StringFunc(func(v string) error {
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
								if data.Name == "" {
									return fmt.Errorf("expected name to be non-empty")
								}
								if data.Provider.Key != "github" {
									return fmt.Errorf("expected provider key to be github, got %s", data.Provider.Key)
								}
								return nil
							}),
						}),
					})),
				},
			},
		},
	})
}

var testAccOrganizationIntegrationsDataSourceConfig = fmt.Sprintf(`
data "sentry_organization_integrations" "test" {
	organization = "%[1]s"
	provider_key = "github"
}

check "raw_json" {
	assert {
		condition = alltrue([
			for i in data.sentry_organization_integrations.test.integrations :
			can(jsondecode(i.raw_json))
		])
		error_message = "One or more integrations contain invalid raw_json format."
	}
}
`, acctest.TestOrganization)
