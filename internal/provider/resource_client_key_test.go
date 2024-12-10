package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccClientKeyResource_UpgradeFromVersion(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
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
				Config:            testAccClientKeyResourceConfig(teamName, projectName, keyName, ""),
				ConfigStateChecks: checks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccClientKeyResourceConfig(teamName, projectName, keyName, ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn"), knownvalue.MapPartial(map[string]knownvalue.Check{
						"public": knownvalue.NotNull(),
						"secret": knownvalue.NotNull(),
						"csp":    knownvalue.NotNull(),
					})),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("dsn_public"), rn, tfjsonpath.New("dsn").AtMapKey("public"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("dsn_secret"), rn, tfjsonpath.New("dsn").AtMapKey("secret"), compare.ValuesSame()),
					statecheck.CompareValuePairs(rn, tfjsonpath.New("dsn_csp"), rn, tfjsonpath.New("dsn").AtMapKey("csp"), compare.ValuesSame()),
				),
			},
		},
	})
}

func TestAccClientKeyResource(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn"), knownvalue.MapPartial(map[string]knownvalue.Check{
			"public": knownvalue.NotNull(),
			"secret": knownvalue.NotNull(),
			"csp":    knownvalue.NotNull(),
		})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
		statecheck.CompareValuePairs(rn, tfjsonpath.New("dsn_public"), rn, tfjsonpath.New("dsn").AtMapKey("public"), compare.ValuesSame()),
		statecheck.CompareValuePairs(rn, tfjsonpath.New("dsn_secret"), rn, tfjsonpath.New("dsn").AtMapKey("secret"), compare.ValuesSame()),
		statecheck.CompareValuePairs(rn, tfjsonpath.New("dsn_csp"), rn, tfjsonpath.New("dsn").AtMapKey("csp"), compare.ValuesSame()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName, ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.Null()),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.Null()),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", `
					rate_limit_window = 1
					rate_limit_count  = 2

					javascript_loader_script {
						performance_monitoring_enabled = true
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.MapExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.Null(),
						"performance_monitoring_enabled": knownvalue.Bool(true),
						"session_replay_enabled":         knownvalue.Null(),
						"debug_enabled":                  knownvalue.Null(),
					})),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", `
					rate_limit_window = 3
					rate_limit_count  = 4

					javascript_loader_script {
						performance_monitoring_enabled = false
						session_replay_enabled         = true
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.MapExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.Null(),
						"performance_monitoring_enabled": knownvalue.Bool(false),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Null(),
					})),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", `
					rate_limit_window = 3
					rate_limit_count  = 4

					javascript_loader_script {
						performance_monitoring_enabled = false
						session_replay_enabled         = true
						debug_enabled                  = true
					}
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.MapExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.Null(),
						"performance_monitoring_enabled": knownvalue.Bool(false),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(true),
					})),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.Null()),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClientKeyResourceConfig(teamName, projectName, keyName, extras string) string {
	return testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
	}) + fmt.Sprintf(`
resource "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[1]s"
	%[2]s
}
`, keyName, extras)
}
