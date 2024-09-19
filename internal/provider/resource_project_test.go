package provider

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/pkg/must"
)

func init() {
	resource.AddTestSweepers("sentry_project", &resource.Sweeper{
		Name: "sentry_project",
		F: func(r string) error {
			ctx := context.Background()

			listParams := &sentry.ListOrganizationProjectsParams{}

			for {
				projects, resp, err := acctest.SharedClient.OrganizationProjects.List(ctx, acctest.TestOrganization, listParams)
				if err != nil {
					return err
				}

				for _, project := range projects {
					if !strings.HasPrefix(project.Slug, "tf-project") {
						continue
					}

					log.Printf("[INFO] Destroying Project: %s", project.Slug)

					_, err := acctest.SharedClient.Projects.Delete(ctx, acctest.TestOrganization, project.Slug)
					if err != nil {
						log.Printf("[ERROR] Failed to destroy Project %q: %s", project.Slug, err)
						continue
					}

					log.Printf("[INFO] Project %q has been destroyed.", project.Slug)
				}

				if resp.Cursor == "" {
					break
				}
				listParams.Cursor = resp.Cursor
			}

			return nil
		},
	})
}

func TestAccProjectResource_basic(t *testing.T) {
	teamName1 := acctest.RandomWithPrefix("tf-team")
	teamName2 := acctest.RandomWithPrefix("tf-team")
	teamName3 := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					AllTeamNames: []string{teamName1, teamName2, teamName3},
					TeamIds:      []int{0, 1},
					ProjectName:  projectName,
					Platform:     "go",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(teamName1),
						knownvalue.StringExact(teamName2),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.StringExact("go")),
				),
			},
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					AllTeamNames: []string{teamName1, teamName2, teamName3},
					TeamIds:      []int{0},
					ProjectName:  projectName + "-renamed",
					Platform:     "python",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(teamName1),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.StringExact("python")),
				),
			},
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					AllTeamNames: []string{teamName1, teamName2, teamName3},
					TeamIds:      []int{2},
					ProjectName:  projectName + "-renamed-again",
					Platform:     "python",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName+"-renamed-again")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(teamName3),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.StringExact("python")),
				),
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					project := rs.Primary.ID
					return buildTwoPartID(organization, project), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectResource_filters(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(teamName)})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Extras: `
						filters = {
							blacklisted_ips = ["127.0.0.1", "0.0.0.0/8"]
							error_messages  = ["TypeError*", "*: integer division or modulo by zero"]
						}
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.MapExact(map[string]knownvalue.Check{
						"blacklisted_ips": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("127.0.0.1"),
							knownvalue.StringExact("0.0.0.0/8"),
						}),
						"releases": knownvalue.Null(),
						"error_messages": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("TypeError*"),
							knownvalue.StringExact("*: integer division or modulo by zero"),
						}),
					})),
				),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Extras: `
						filters = {
							blacklisted_ips = ["0.0.0.0/8"]
							releases        = ["1.*", "[!3].[0-9].*"]
							error_messages  = ["TypeError*"]
						}
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.MapExact(map[string]knownvalue.Check{
						"blacklisted_ips": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("0.0.0.0/8"),
						}),
						"releases": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("1.*"),
							knownvalue.StringExact("[!3].[0-9].*"),
						}),
						"error_messages": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("TypeError*"),
						}),
					})),
				),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Extras: `
						filters = {}
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.MapExact(map[string]knownvalue.Check{
						"blacklisted_ips": knownvalue.Null(),
						"releases":        knownvalue.Null(),
						"error_messages":  knownvalue.Null(),
					})),
				),
			},
		},
	})
}

func TestAccProjectResource_noDefaultKeyOnCreate(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Extras: `
						default_key = false
					`,
				}) + `
					data "sentry_all_keys" "test" {
						organization = sentry_project.test.organization
						project      = sentry_project.test.slug
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(teamName)})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.sentry_all_keys.test", tfjsonpath.New("keys"), knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}

func TestAccProjectResource_noDefaultKeyOnUpdate(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(teamName)})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
				}) + `
					data "sentry_all_keys" "test" {
						organization = sentry_project.test.organization
						project      = sentry_project.test.slug
					}
				`,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.sentry_all_keys.test", tfjsonpath.New("keys"), knownvalue.ListSizeExact(1)),
				),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Extras: `
						default_key = false
					`,
				}) + `
					data "sentry_all_keys" "test" {
						organization = sentry_project.test.organization
						project      = sentry_project.test.slug
					}
				`,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("data.sentry_all_keys.test", tfjsonpath.New("keys"), knownvalue.ListSizeExact(0)),
				),
			},
		},
	})
}

func TestAccProjectResource_invalidPlatform(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Platform:    "invalid",
				}),
				ExpectError: regexp.MustCompile(`Attribute platform value must be one of`),
			},
		},
	})
}

func TestAccProjectResource_noTeam(t *testing.T) {
	projectName := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					ProjectName: projectName,
				}),
				ExpectError: regexp.MustCompile(`Attribute teams set must contain at least 1 elements, got: 0`),
			},
		},
	})
}

func TestAccProjectResource_UpgradeFromVersion(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact(teamName)})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.StringExact("go")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		CheckDestroy: testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.12.3",
					},
				},
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Platform:    "go",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.NotNull()),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Platform:    "go",
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
				),
			},
		},
	})
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_project" {
			continue
		}

		ctx := context.Background()
		project, resp, err := acctest.SharedClient.Projects.Get(
			ctx,
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err == nil {
			if project != nil {
				return fmt.Errorf("project %q still exists", rs.Primary.ID)
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

var testAccProjectResourceConfigTemplate = template.Must(template.New("config").Parse(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "{{ .ProjectName }}"
	{{ if .Platform }}
	platform = "{{ .Platform }}"
	{{ end }}
	{{ .Extras }}
}
`))

type testAccProjectResourceConfigData struct {
	TeamName    string
	ProjectName string
	Platform    string
	Extras      string
}

func testAccProjectResourceConfig(data testAccProjectResourceConfigData) string {
	var builder strings.Builder

	must.Get(builder.WriteString(testAccTeamResourceConfig(data.TeamName)))
	must.Do(testAccProjectResourceConfigTemplate.Execute(&builder, data))

	return builder.String()
}

var testAccProjectResourceConfig_teamsTemplate = template.Must(template.New("config").Parse(`
{{ range $i, $teamName := .AllTeamNames }}
resource "sentry_team" "team_{{ $i }}" {
	organization = data.sentry_organization.test.id
	name         = "{{ $teamName }}"
	slug         = "{{ $teamName }}"
}
{{ end }}

resource "sentry_project" "test" {
	organization = data.sentry_organization.test.id
	teams        = [
		{{ range $i, $TeamId := .TeamIds }}
		sentry_team.team_{{ $TeamId }}.slug,
		{{ end }}
	]
	name         = "{{ .ProjectName }}"
	{{ if .Platform }}
	platform     = "{{ .Platform }}"
	{{ end }}
}
`))

type testAccProjectResourceConfig_teamsData struct {
	AllTeamNames []string
	TeamIds      []int
	ProjectName  string
	Platform     string
}

func testAccProjectResourceConfig_teams(data testAccProjectResourceConfig_teamsData) string {
	var builder strings.Builder

	must.Get(builder.WriteString(testAccOrganizationDataSourceConfig))
	must.Do(testAccProjectResourceConfig_teamsTemplate.Execute(&builder, data))

	return builder.String()
}
