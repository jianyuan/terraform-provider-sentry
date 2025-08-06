package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/go-utils/sliceutils"

	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
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

					_, err := acctest.SharedApiClient.DeleteOrganizationProjectWithResponse(
						ctx,
						acctest.TestOrganization,
						project.Slug,
					)
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

	checkProperties := func(data testAccProjectResourceConfig_teamsData) func(apiclient.Project) error {
		return func(project apiclient.Project) error {
			if project.Name != data.ProjectName {
				return fmt.Errorf("unexpected project name %v", project.Name)
			}

			if v, err := project.Platform.Get(); err == nil {
				if data.Platform != v {
					return fmt.Errorf("unexpected platform %v", v)
				}
			} else {
				return fmt.Errorf("unexpected platform: %s", err)
			}

			if len(project.Teams) != len(data.TeamIds) {
				return fmt.Errorf("unexpected number of teams %v", len(project.Teams))
			}

			for _, teamId := range data.TeamIds {
				if !slices.ContainsFunc(project.Teams, func(team apiclient.Team) bool {
					return team.Slug == data.AllTeamNames[teamId]
				}) {
					return fmt.Errorf("team %v not found", data.AllTeamNames[teamId])
				}
			}

			if len(project.AllowedDomains) != len(ptr.Value(data.AllowedDomains)) {
				return fmt.Errorf("unexpected allowed domains %v", project.AllowedDomains)
			}

			for _, allowedDomain := range *data.AllowedDomains {
				if !slices.Contains(project.AllowedDomains, allowedDomain) {
					return fmt.Errorf("allowed domain %v not found", allowedDomain)
				}
			}

			if project.ScrapeJavaScript != ptr.Value(data.ScrapeJavascript) {
				return fmt.Errorf("unexpected scrape javascript %v", project.ScrapeJavaScript)
			}

			if v, err := project.SecurityTokenHeader.Get(); err == nil {
				if v != ptr.Value(data.SecurityTokenHeader) {
					return fmt.Errorf("unexpected security token header %v", v)
				}
			} else {
				if ptr.Value(data.SecurityTokenHeader) != "" {
					return fmt.Errorf("unexpected security token header")
				}
			}

			if project.VerifySSL != ptr.Value(data.VerifyTlsSsl) {
				return fmt.Errorf("unexpected verify tls ssl %v", project.VerifySSL)
			}

			if data.HighlightTags != nil {
				if project.HighlightTags == nil {
					return fmt.Errorf("highlight tags is nil")
				}

				if len(*project.HighlightTags) != len(*data.HighlightTags) {
					return fmt.Errorf("unexpected highlight tags %v", *project.HighlightTags)
				}

				for _, tag := range *data.HighlightTags {
					if !slices.Contains(*project.HighlightTags, tag) {
						return fmt.Errorf("highlight tag %v not found", tag)
					}
				}
			}

			return nil
		}
	}

	configStateChecks := func(data testAccProjectResourceConfig_teamsData) []statecheck.StateCheck {
		return []statecheck.StateCheck{
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("teams"), knownvalue.SetExact(sliceutils.Map(func(teamId int) knownvalue.Check {
				return knownvalue.StringExact(data.AllTeamNames[teamId])
			}, data.TeamIds))),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(data.ProjectName)),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("slug"), knownvalue.NotNull()),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("platform"), knownvalue.StringExact(data.Platform)),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("features"), knownvalue.NotNull()),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Int64Exact(300)),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Int64Exact(1800)),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Int64Exact(0)),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
				"blacklisted_ips": knownvalue.SetSizeExact(0),
				"releases":        knownvalue.SetSizeExact(0),
				"error_messages":  knownvalue.SetSizeExact(0),
			})),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("")),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("")),
			statecheck.ExpectKnownValue(rn, tfjsonpath.New("client_security"), knownvalue.ObjectExact(map[string]knownvalue.Check{
				"allowed_domains": knownvalue.SetExact(sliceutils.Map(func(v string) knownvalue.Check {
					return knownvalue.StringExact(v)
				}, ptr.Value(data.AllowedDomains))),
				"scrape_javascript":     knownvalue.Bool(ptr.Value(data.ScrapeJavascript)),
				"security_token":        knownvalue.NotNull(),
				"security_token_header": knownvalue.StringExact(ptr.Value(data.SecurityTokenHeader)),
				"verify_tls_ssl":        knownvalue.Bool(ptr.Value(data.VerifyTlsSsl)),
			})),
			func() statecheck.StateCheck {
				if data.HighlightTags == nil {
					return statecheck.ExpectKnownValue(rn, tfjsonpath.New("highlight_tags"), knownvalue.Null())
				}
				return statecheck.ExpectKnownValue(rn, tfjsonpath.New("highlight_tags"), knownvalue.SetExact(sliceutils.Map(func(v string) knownvalue.Check {
					return knownvalue.StringExact(v)
				}, *data.HighlightTags)))
			}(),
		}
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
				Check: testAccCheckProject(rn, checkProperties(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{0, 1},
					ProjectName:         projectName,
					Platform:            "go",
					AllowedDomains:      ptr.Ptr([]string{"*"}),
					ScrapeJavascript:    ptr.Ptr(true),
					SecurityTokenHeader: ptr.Ptr(""),
					VerifyTlsSsl:        ptr.Ptr(false),
				})),
				ConfigStateChecks: configStateChecks(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{0, 1},
					ProjectName:         projectName,
					Platform:            "go",
					AllowedDomains:      ptr.Ptr([]string{"*"}),
					ScrapeJavascript:    ptr.Ptr(true),
					SecurityTokenHeader: ptr.Ptr(""),
					VerifyTlsSsl:        ptr.Ptr(false),
				}),
			},
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					AllTeamNames:     []string{teamName1, teamName2, teamName3},
					TeamIds:          []int{0},
					ProjectName:      projectName + "-renamed",
					Platform:         "python",
					ScrapeJavascript: ptr.Ptr(false),
					VerifyTlsSsl:     ptr.Ptr(true),
				}),
				Check: testAccCheckProject(rn, checkProperties(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{0},
					ProjectName:         projectName + "-renamed",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"*"}),
					ScrapeJavascript:    ptr.Ptr(false),
					SecurityTokenHeader: ptr.Ptr(""),
					VerifyTlsSsl:        ptr.Ptr(true),
				})),
				ConfigStateChecks: configStateChecks(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{0},
					ProjectName:         projectName + "-renamed",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"*"}),
					ScrapeJavascript:    ptr.Ptr(false),
					SecurityTokenHeader: ptr.Ptr(""),
					VerifyTlsSsl:        ptr.Ptr(true),
				}),
			},
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{2},
					ProjectName:         projectName + "-renamed-again",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"jianyuan.io", "*.jianyuan.io"}),
					SecurityTokenHeader: ptr.Ptr("x-my-security-token"),
					HighlightTags:       ptr.Ptr([]string{"release", "environment"}),
				}),
				Check: testAccCheckProject(rn, checkProperties(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{2},
					ProjectName:         projectName + "-renamed-again",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"jianyuan.io", "*.jianyuan.io"}),
					ScrapeJavascript:    ptr.Ptr(false),
					SecurityTokenHeader: ptr.Ptr("x-my-security-token"),
					VerifyTlsSsl:        ptr.Ptr(true),
					HighlightTags:       ptr.Ptr([]string{"release", "environment"}),
				})),
				ConfigStateChecks: configStateChecks(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{2},
					ProjectName:         projectName + "-renamed-again",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"jianyuan.io", "*.jianyuan.io"}),
					ScrapeJavascript:    ptr.Ptr(false),
					SecurityTokenHeader: ptr.Ptr("x-my-security-token"),
					VerifyTlsSsl:        ptr.Ptr(true),
					HighlightTags:       ptr.Ptr([]string{"release", "environment"}),
				}),
			},
			// Remove all optional attributes
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					AllTeamNames: []string{teamName1, teamName2, teamName3},
					TeamIds:      []int{2},
					ProjectName:  projectName + "-renamed-again",
					Platform:     "python",
				}),
				Check: testAccCheckProject(rn, checkProperties(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{2},
					ProjectName:         projectName + "-renamed-again",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"jianyuan.io", "*.jianyuan.io"}),
					ScrapeJavascript:    ptr.Ptr(false),
					SecurityTokenHeader: ptr.Ptr("x-my-security-token"),
					VerifyTlsSsl:        ptr.Ptr(true),
				})),
				ConfigStateChecks: configStateChecks(testAccProjectResourceConfig_teamsData{
					AllTeamNames:        []string{teamName1, teamName2, teamName3},
					TeamIds:             []int{2},
					ProjectName:         projectName + "-renamed-again",
					Platform:            "python",
					AllowedDomains:      ptr.Ptr([]string{"jianyuan.io", "*.jianyuan.io"}),
					ScrapeJavascript:    ptr.Ptr(false),
					SecurityTokenHeader: ptr.Ptr("x-my-security-token"),
					VerifyTlsSsl:        ptr.Ptr(true),
				}),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.TwoPartImportStateIdFunc(rn, "organization"),
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
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Int64Exact(300)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Int64Exact(1800)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Int64Exact(0)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("client_security"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"allowed_domains": knownvalue.SetExact([]knownvalue.Check{
				knownvalue.StringExact("*"),
			}),
			"scrape_javascript":     knownvalue.Bool(true),
			"security_token":        knownvalue.NotNull(),
			"security_token_header": knownvalue.StringExact(""),
			"verify_tls_ssl":        knownvalue.Bool(false),
		})),
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"blacklisted_ips": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("127.0.0.1"),
							knownvalue.StringExact("0.0.0.0/8"),
						}),
						"releases": knownvalue.SetSizeExact(0),
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
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
		},
	})
}

func TestAccProjectResource_issueGrouping(t *testing.T) {
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
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Int64Exact(300)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Int64Exact(1800)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Int64Exact(0)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"blacklisted_ips": knownvalue.SetSizeExact(0),
			"releases":        knownvalue.SetSizeExact(0),
			"error_messages":  knownvalue.SetSizeExact(0),
		})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("client_security"), knownvalue.ObjectExact(map[string]knownvalue.Check{
			"allowed_domains": knownvalue.SetExact([]knownvalue.Check{
				knownvalue.StringExact("*"),
			}),
			"scrape_javascript":     knownvalue.Bool(true),
			"security_token":        knownvalue.NotNull(),
			"security_token_header": knownvalue.StringExact(""),
			"verify_tls_ssl":        knownvalue.Bool(false),
		})),
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
						fingerprinting_rules = <<-EOT
							# force all errors of the same type to have the same fingerprint
							error.type:DatabaseUnavailable -> system-down
						EOT
						grouping_enhancements = <<-EOT
							# remove all frames above a certain function from grouping
							stack.function:panic_handler ^-group
						EOT
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("# force all errors of the same type to have the same fingerprint\nerror.type:DatabaseUnavailable -> system-down\n")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("# remove all frames above a certain function from grouping\nstack.function:panic_handler ^-group\n")),
				),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Extras: `
						fingerprinting_rules = <<-EOT
							# force all errors of the same type to have the same fingerprint
							error.type:DatabaseUnavailable -> system-down
							# force all memory allocation errors to be grouped together
							stack.function:malloc -> memory-allocation-error
						EOT
						grouping_enhancements = <<-EOT
							# remove all frames above a certain function from grouping
							stack.function:panic_handler ^-group
							# mark all functions following a prefix in-app
							stack.function:mylibrary_* +app
						EOT
					`,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("# force all errors of the same type to have the same fingerprint\nerror.type:DatabaseUnavailable -> system-down\n# force all memory allocation errors to be grouped together\nstack.function:malloc -> memory-allocation-error\n")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("# remove all frames above a certain function from grouping\nstack.function:panic_handler ^-group\n# mark all functions following a prefix in-app\nstack.function:mylibrary_* +app\n")),
				),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
				}),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("# force all errors of the same type to have the same fingerprint\nerror.type:DatabaseUnavailable -> system-down\n# force all memory allocation errors to be grouped together\nstack.function:malloc -> memory-allocation-error\n")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("# remove all frames above a certain function from grouping\nstack.function:panic_handler ^-group\n# mark all functions following a prefix in-app\nstack.function:mylibrary_* +app\n")),
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
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Int64Exact(300)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Int64Exact(1800)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"blacklisted_ips": knownvalue.SetSizeExact(0),
						"releases":        knownvalue.SetSizeExact(0),
						"error_messages":  knownvalue.SetSizeExact(0),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("client_security"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"allowed_domains": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("*"),
						}),
						"scrape_javascript":     knownvalue.Bool(true),
						"security_token":        knownvalue.NotNull(),
						"security_token_header": knownvalue.StringExact(""),
						"verify_tls_ssl":        knownvalue.Bool(false),
					})),
					statecheck.ExpectKnownValue("data.sentry_all_keys.test", tfjsonpath.New("keys"), knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}

func TestAccProjectResource_validation(t *testing.T) {
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
			{
				Config: testAccProjectResourceConfig_teams(testAccProjectResourceConfig_teamsData{
					ProjectName: projectName,
				}),
				ExpectError: regexp.MustCompile(`Attribute teams set must contain at least 1 elements, got: 0`),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Platform:    "go",
					Extras: `
						client_security = {
							allowed_domains = []
						}
					`,
				}),
				ExpectError: regexp.MustCompile(`Attribute client_security.allowed_domains set must contain at least 1\nelements, got: 0`),
			},
			{
				Config: testAccProjectResourceConfig(testAccProjectResourceConfigData{
					TeamName:    teamName,
					ProjectName: projectName,
					Platform:    "go",
					Extras: `
						client_security = {
							security_token_header = "012345678901234567890"
						}
					`,
				}),
				ExpectError: regexp.MustCompile(`Attribute client_security.security_token_header string length must be at most\n20, got: 21`),
			},
		},
	})
}

func TestAccProjectResource_upgradeFromVersion(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_project.test"

	config := testAccProjectResourceConfig(testAccProjectResourceConfigData{
		TeamName:    teamName,
		ProjectName: projectName,
		Platform:    "go",
	})

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
				Config: config,
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
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.14.1",
					},
				},
				Config: config,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.Null()),
				),
			},
			// Some optional fields are now computed
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   config,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_rules"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("default_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_min_delay"), knownvalue.Int64Exact(300)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("digests_max_delay"), knownvalue.Int64Exact(1800)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("resolve_age"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"blacklisted_ips": knownvalue.SetSizeExact(0),
						"releases":        knownvalue.SetSizeExact(0),
						"error_messages":  knownvalue.SetSizeExact(0),
					})),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("fingerprinting_rules"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("grouping_enhancements"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("client_security"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"allowed_domains": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("*"),
						}),
						"scrape_javascript":     knownvalue.Bool(true),
						"security_token":        knownvalue.NotNull(),
						"security_token_header": knownvalue.StringExact(""),
						"verify_tls_ssl":        knownvalue.Bool(false),
					})),
				),
			},
		},
	})
}

func testAccCheckProject(rn string, check func(apiclient.Project) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[rn]

		httpResp, err := acctest.SharedApiClient.GetOrganizationProjectWithResponse(
			context.Background(),
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		} else if httpResp.StatusCode() == http.StatusNotFound {
			return fmt.Errorf("project %q not found", rs.Primary.ID)
		} else if httpResp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unexpected status code %d: %s", httpResp.StatusCode(), string(httpResp.Body))
		} else if httpResp.JSON200 == nil {
			return fmt.Errorf("response body is empty")
		}

		return check(*httpResp.JSON200)
	}
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_project" {
			continue
		}

		httpResp, err := acctest.SharedApiClient.GetOrganizationProjectWithResponse(
			context.Background(),
			rs.Primary.Attributes["organization"],
			rs.Primary.ID,
		)
		if err != nil {
			return err
		} else if httpResp.StatusCode() != http.StatusNotFound {
			return fmt.Errorf("project %q still exists", rs.Primary.ID)
		}

		return nil
	}

	return nil
}

var testAccProjectResourceConfigTemplate = template.Must(template.New("config").Parse(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.slug]
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
	organization = data.sentry_organization.test.slug
	name         = "{{ $teamName }}"
	slug         = "{{ $teamName }}"
}
{{ end }}

resource "sentry_project" "test" {
	organization = data.sentry_organization.test.slug
	teams        = [
		{{ range $i, $TeamId := .TeamIds }}
		sentry_team.team_{{ $TeamId }}.slug,
		{{ end }}
	]
	name         = "{{ .ProjectName }}"
	{{ if .Platform }}
	platform     = "{{ .Platform }}"
	{{ end }}

	client_security = {
		{{ if ne .AllowedDomains nil }}
		allowed_domains = [
			{{ range $i, $domain := .AllowedDomains }}
			"{{ $domain }}",
			{{ end }}
		],
		{{ end }}
		{{ if ne .ScrapeJavascript nil }}
		scrape_javascript = {{ .ScrapeJavascript }}
		{{ end }}
		{{ if ne .SecurityToken nil }}
		security_token = "{{ .SecurityToken }}"
		{{ end }}
		{{ if ne .SecurityTokenHeader nil }}
		security_token_header = "{{ .SecurityTokenHeader }}"
		{{ end }}
		{{ if ne .VerifyTlsSsl nil }}
		verify_tls_ssl = {{ .VerifyTlsSsl }}
		{{ end }}
	}

	{{ if ne .HighlightTags nil }}
	highlight_tags = [
		{{ range $i, $tag := .HighlightTags }}
		"{{ $tag }}",
		{{ end }}
	]
	{{ end }}
}
`))

type testAccProjectResourceConfig_teamsData struct {
	AllTeamNames        []string
	TeamIds             []int
	ProjectName         string
	Platform            string
	AllowedDomains      *[]string
	ScrapeJavascript    *bool
	SecurityToken       *string
	SecurityTokenHeader *string
	VerifyTlsSsl        *bool
	HighlightTags       *[]string
}

func testAccProjectResourceConfig_teams(data testAccProjectResourceConfig_teamsData) string {
	var builder strings.Builder

	must.Get(builder.WriteString(testAccOrganizationDataSourceConfig))
	must.Do(testAccProjectResourceConfig_teamsTemplate.Execute(&builder, data))

	return builder.String()
}
