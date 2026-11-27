package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

// testAccProtoV6ProviderFactoriesWithEcho adds the echo provider, which copies
// its `data` argument into the `echo` resource so that ephemeral values can be
// checked.
var testAccProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	acctest.ProviderName: testAccProtoV6ProviderFactories[acctest.ProviderName],
	"echo":               echoprovider.NewProviderServer(),
}

var _ statecheck.StateCheck = expectNoEphemeralResourcesInState{}

// expectNoEphemeralResourcesInState fails if any ephemeral resource was
// persisted to state.
type expectNoEphemeralResourcesInState struct{}

func (expectNoEphemeralResourcesInState) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	if req.State == nil || req.State.Values == nil || req.State.Values.RootModule == nil {
		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if r.Mode == "ephemeral" || strings.HasPrefix(r.Address, "ephemeral.") {
			resp.Error = fmt.Errorf("ephemeral resource %s found in state", r.Address)
			return
		}
	}
}

func testAccClientKeyEphemeralResourceChecks(name knownvalue.Check, first knownvalue.Check) []statecheck.StateCheck {
	rn := "echo.test"
	return []statecheck.StateCheck{
		expectNoEphemeralResourcesInState{},
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("name"), name),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("first"), first),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("rate_limit_window"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("rate_limit_count"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("javascript_loader_script"), knownvalue.MapExact(map[string]knownvalue.Check{
			"browser_sdk_version":            knownvalue.NotNull(),
			"performance_monitoring_enabled": knownvalue.Bool(true),
			"session_replay_enabled":         knownvalue.Bool(true),
			"debug_enabled":                  knownvalue.Bool(false),
		})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("data").AtMapKey("dsn"), knownvalue.MapPartial(map[string]knownvalue.Check{
			"public": knownvalue.NotNull(),
			"secret": knownvalue.NotNull(),
			"csp":    knownvalue.NotNull(),
		})),
	}
}

func TestAccClientKeyEphemeralResource_id(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyEphemeralResourceConfig(teamName, projectName, keyName, `
					id = sentry_key.test.id
				`),
				ConfigStateChecks: append(
					testAccClientKeyEphemeralResourceChecks(knownvalue.StringExact(keyName), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("project"), knownvalue.StringExact(projectName)),
				),
			},
		},
	})
}

func TestAccClientKeyEphemeralResource_name(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyEphemeralResourceConfig(teamName, projectName, keyName, `
					name = sentry_key.test.name
				`),
				ConfigStateChecks: append(
					testAccClientKeyEphemeralResourceChecks(knownvalue.StringExact(keyName), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("project"), knownvalue.StringExact(projectName)),
				),
			},
		},
	})
}

func TestAccClientKeyEphemeralResource_first(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyEphemeralResourceConfig(teamName, projectName, keyName, `
					first = true
				`),
				ConfigStateChecks: append(
					testAccClientKeyEphemeralResourceChecks(knownvalue.StringExact("Default"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("project"), knownvalue.StringExact(projectName)),
				),
			},
		},
	})
}

func testAccClientKeyEphemeralResourceConfig(teamName, projectName, keyName, extras string) string {
	return testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		KeyName:     keyName,
	}) + fmt.Sprintf(`
ephemeral "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	%[1]s
}

provider "echo" {
	data = ephemeral.sentry_key.test
}

resource "echo" "test" {}
`, extras)
}
