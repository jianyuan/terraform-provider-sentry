package provider

import (
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccClientKeyResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "sentry_key" "test" {
						organization = "value"
						project      = "value"
						name 	     = "value"

						rate_limit_window = 1
					}
				`,
				ExpectError: regexp.MustCompile(`These attributes must be configured together:\n\[rate_limit_window,rate_limit_count\]`),
			},
			{
				Config: `
					resource "sentry_key" "test" {
						organization = "value"
						project      = "value"
						name 	     = "value"

						rate_limit_count = 1
					}
				`,
				ExpectError: regexp.MustCompile(`These attributes must be configured together:\n\[rate_limit_window,rate_limit_count\]`),
			},
			{
				Config: `
					resource "sentry_key" "test" {
						organization = "value"
						project      = "value"
						name 	     = "value"

						javascript_loader_script = {
							performance_monitoring_enabled = true
						}
					}
				`,
				ExpectError: regexp.MustCompile(`These attributes must be configured together:\n\[javascript_loader_script.performance_monitoring_enabled,javascript_loader_script.session_replay_enabled,javascript_loader_script.debug_enabled\]`),
			},
			{
				Config: `
					resource "sentry_key" "test" {
						organization = "value"
						project      = "value"
						name 	     = "value"

						javascript_loader_script = {
							session_replay_enabled = true
						}
					}
				`,
				ExpectError: regexp.MustCompile(`These attributes must be configured together:\n\[javascript_loader_script.performance_monitoring_enabled,javascript_loader_script.session_replay_enabled,javascript_loader_script.debug_enabled\]`),
			},
			{
				Config: `
					resource "sentry_key" "test" {
						organization = "value"
						project      = "value"
						name 	     = "value"

						javascript_loader_script = {
							debug_enabled = true
						}
					}
				`,
				ExpectError: regexp.MustCompile(`These attributes must be configured together:\n\[javascript_loader_script.performance_monitoring_enabled,javascript_loader_script.session_replay_enabled,javascript_loader_script.debug_enabled\]`),
			},
		},
	})
}

func TestAccClientKeyResource_upgradeFromVersion(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	config := testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		KeyName:     keyName,
	})

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
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
				Config:            config,
				ConfigStateChecks: checks,
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.14.1",
					},
				},
				Config: config,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
				),
			},
			// Some optional fields are now computed
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   config,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.NotNull(),
						"performance_monitoring_enabled": knownvalue.Bool(true),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(false),
					})),
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
				Config: testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					KeyName:     keyName,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.NotNull(),
						"performance_monitoring_enabled": knownvalue.Bool(true),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(false),
					})),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					KeyName:     keyName + "-renamed",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.NotNull(),
						"performance_monitoring_enabled": knownvalue.Bool(true),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(false),
					})),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
					TeamName:        teamName,
					ProjectName:     projectName,
					KeyName:         keyName + "-renamed",
					RateLimitWindow: ptr.Ptr(1),
					RateLimitCount:  ptr.Ptr(2),
					Extras: `
						javascript_loader_script = {
							browser_sdk_version            = "7.x"
						}
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.StringExact("7.x"),
						"performance_monitoring_enabled": knownvalue.Bool(true),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(false),
					})),
				),
			},
			// Remove browser_sdk_version
			// Add other flags
			{
				Config: testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
					TeamName:        teamName,
					ProjectName:     projectName,
					KeyName:         keyName + "-renamed",
					RateLimitWindow: ptr.Ptr(3),
					RateLimitCount:  ptr.Ptr(4),
					Extras: `
						javascript_loader_script = {
							performance_monitoring_enabled = false
							session_replay_enabled         = false
							debug_enabled                  = true
						}
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.StringExact("7.x"),
						"performance_monitoring_enabled": knownvalue.Bool(false),
						"session_replay_enabled":         knownvalue.Bool(false),
						"debug_enabled":                  knownvalue.Bool(true),
					})),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
					TeamName:        teamName,
					ProjectName:     projectName,
					KeyName:         keyName + "-renamed",
					RateLimitWindow: ptr.Ptr(3),
					RateLimitCount:  ptr.Ptr(4),
					Extras: `
						javascript_loader_script = {
							performance_monitoring_enabled = false
							session_replay_enabled         = true
							debug_enabled                  = true
						}
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.StringExact("7.x"),
						"performance_monitoring_enabled": knownvalue.Bool(false),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(true),
					})),
				),
			},
			// Remove all optional fields
			{
				Config: testAccClientKeyResourceConfig(testAccClientKeyResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					KeyName:     keyName + "-renamed",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("javascript_loader_script"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"browser_sdk_version":            knownvalue.StringExact("7.x"),
						"performance_monitoring_enabled": knownvalue.Bool(false),
						"session_replay_enabled":         knownvalue.Bool(true),
						"debug_enabled":                  knownvalue.Bool(true),
					})),
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

var testAccClientKeyResourceConfigTemplate = template.Must(template.New("config").Parse(`
resource "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "{{ .KeyName }}"

	{{ if ne .RateLimitWindow nil }}
	rate_limit_window = {{ .RateLimitWindow }}
	{{ end }}
	{{ if ne .RateLimitCount nil }}
	rate_limit_count  = {{ .RateLimitCount }}
	{{ end }}
	{{ .Extras }}
}
`))

type testAccClientKeyResourceConfigData struct {
	TeamName        string
	ProjectName     string
	KeyName         string
	RateLimitWindow *int
	RateLimitCount  *int
	Extras          string
}

func testAccClientKeyResourceConfig(data testAccClientKeyResourceConfigData) string {
	var builder strings.Builder

	must.Get(builder.WriteString(testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    data.TeamName,
		ProjectName: data.ProjectName,
	})))
	must.Do(testAccClientKeyResourceConfigTemplate.Execute(&builder, data))

	return builder.String()
}
