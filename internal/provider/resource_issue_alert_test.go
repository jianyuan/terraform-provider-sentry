package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccIssueAlertResource_XXX(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition"), knownvalue.ListExact([]knownvalue.Check{
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"first_seen_event": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"name": knownvalue.StringExact("A new issue is created"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"regression_event": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"name": knownvalue.StringExact("The issue changes state from resolved to unresolved"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"event_frequency": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"name":     knownvalue.StringExact("The issue is seen more than 500 times in 1h"),
						"value":    knownvalue.Int64Exact(500),
						"interval": knownvalue.StringExact("1h"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"event_unique_user_frequency": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"name":     knownvalue.StringExact("The issue is seen by more than 1000 users in 15m"),
						"value":    knownvalue.Int64Exact(1000),
						"interval": knownvalue.StringExact("15m"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"event_frequency_percent": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"name":     knownvalue.StringExact("The issue affects more than 50.0 percent of sessions in 10m"),
						"value":    knownvalue.Float64Exact(50.0),
						"interval": knownvalue.StringExact("10m"),
					}),
				}),
			}),
		})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter"), knownvalue.ListExact([]knownvalue.Check{})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action"), knownvalue.ListExact([]knownvalue.Check{
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"notify_email": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":                knownvalue.NotNull(),
						"name":              knownvalue.StringExact("Send a notification to IssueOwners and if none can be found then send a notification to ActiveMembers"),
						"target_type":       knownvalue.StringExact("IssueOwners"),
						"target_identifier": knownvalue.Null(),
						"fallthrough_type":  knownvalue.StringExact("ActiveMembers"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"notify_email": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":                knownvalue.NotNull(),
						"name":              knownvalue.StringExact("Send a notification to Team and if none can be found then send a notification to AllMembers"),
						"target_type":       knownvalue.StringExact("Team"),
						"target_identifier": knownvalue.NotNull(),
						"fallthrough_type":  knownvalue.StringExact("AllMembers"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"slack_notify_service": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":         knownvalue.NotNull(),
						"name":       knownvalue.StringRegexp(regexp.MustCompile(`Send a notification to Slack workspace ".*" in channel #general with tags environment,level and notes "Please <http://example.com|click here> for triage information"`)),
						"workspace":  knownvalue.NotNull(),
						"channel":    knownvalue.StringExact("#general"),
						"channel_id": knownvalue.NotNull(),
						"tags":       knownvalue.StringExact("environment,level"),
						"notes":      knownvalue.StringExact("Please <http://example.com|click here> for triage information"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"github_create_ticket": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":          knownvalue.NotNull(),
						"name":        knownvalue.StringExact("Create a GitHub issue in jianyuan with these "),
						"integration": knownvalue.NotNull(),
						"repo":        knownvalue.StringExact("terraform-provider-sentry"),
						"title":       knownvalue.StringExact("My Test Issue"),
						"body":        knownvalue.Null(),
						"assignee":    knownvalue.StringExact("jianyuan"),
						"labels": knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("bug"),
							knownvalue.StringExact("enhancement"),
						}),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"azure_devops_create_ticket": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":             knownvalue.NotNull(),
						"name":           knownvalue.StringExact("Create an Azure DevOps work item in jianyuanlee with these "),
						"integration":    knownvalue.NotNull(),
						"project":        knownvalue.StringExact("123"),
						"work_item_type": knownvalue.StringExact("Microsoft.VSTS.WorkItemTypes.Task"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"pagerduty_notify_service": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":       knownvalue.NotNull(),
						"name":     knownvalue.StringExact("Send a notification to PagerDuty account terraform-provider-sentry and service issue-alert-service with default severity"),
						"account":  knownvalue.NotNull(),
						"service":  knownvalue.NotNull(),
						"severity": knownvalue.StringExact("default"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"opsgenie_notify_team": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":       knownvalue.NotNull(),
						"name":     knownvalue.StringExact("Send a notification to Opsgenie account terraform-provider-sentry and team issue-alert-team with P1 priority"),
						"account":  knownvalue.NotNull(),
						"team":     knownvalue.NotNull(),
						"priority": knownvalue.StringExact("P1"),
					}),
				}),
			}),
			knownvalue.ObjectPartial(map[string]knownvalue.Check{
				"notify_event": knownvalue.ListExact([]knownvalue.Check{
					knownvalue.ObjectExact(map[string]knownvalue.Check{
						"id":   knownvalue.NotNull(),
						"name": knownvalue.StringExact("Send a notification (for all legacy integrations)"),
					}),
				}),
			}),
		})),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig(team, project, alert),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			// {
			// 	Config: testAccIssueAlertConfig(team, project, alert+"-updated"),
			// 	ConfigStateChecks: append(
			// 		configStateChecks,
			// 		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
			// 	),
			// },
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					return buildThreePartID(rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], rs.Primary.ID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conditions", "filters", "actions"},
			},
		},
	})
}

func TestAccIssueAlertResource_JsonValues(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert-with-a-very-looooong-name-greater-than-64-characters")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition"), knownvalue.SetExact([]knownvalue.Check{
			// knownvalue
		})),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertJsonConfig(team, project, alert),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertJsonConfig(team, project, alert+"-updated"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
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
					return buildThreePartID(rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], rs.Primary.ID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conditions", "filters", "actions"},
			},
		},
	})
}

func TestAccIssueAlertResource_Validation(t *testing.T) {
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert-with-a-very-looooong-name-greater-than-64-characters")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertBaseConfig(team, project, alert, `
actions = <<EOT
[
	{
		"id": "sentry.rules.actions.notify_event.NotifyEventAction"
	}
]
EOT

condition {
	first_seen_event {}
}
action {
	notify_event {}
}
`),
				ExpectError: regexp.MustCompile(`Attribute "action" cannot be specified when "actions" is specified`),
			},
			{
				Config: testAccIssueAlertBaseConfig(team, project, alert, `
conditions = <<EOT
[
	{
		"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
		"name": "ignored"
	}
]
EOT

condition {
	first_seen_event {}
}
action {
	notify_event {}
}
`),
				ExpectError: regexp.MustCompile(`Attribute "condition" cannot be specified when "conditions" is specified`),
			},
			{
				Config: testAccIssueAlertBaseConfig(team, project, alert, `
condition {
	first_seen_event {}
	regression_event {}
}

condition {
	first_seen_event {}
	regression_event {}
}

action {
	notify_event {}
}
`),
				ExpectError: regexp.MustCompile(`Duplicate condition block`),
			},
		},
	})
}

func TestAccIssueAlertResource_UpgradeFromVersion(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.11.2",
					},
				},
				Config: testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = [
		{
			id = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
		},
		{
			id = "sentry.rules.conditions.regression_event.RegressionEventCondition"
		}
	]

	actions = [
		{
			id = "sentry.rules.actions.notify_event.NotifyEventAction"
		}
	]
}
`, team, project, alert),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("internal_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
				},
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config: testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = jsonencode(
		[
			{
				id = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
			},
			{
				id = "sentry.rules.conditions.regression_event.RegressionEventCondition"
			}
		]
	)

	actions = jsonencode(
		[
			{
				"id": "sentry.rules.actions.notify_event.NotifyEventAction"
			}
		]
	)
}
`, team, project, alert),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccIssueAlertResource_EmptyArray(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("condition"), knownvalue.SetSizeExact(0)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter"), knownvalue.SetSizeExact(0)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action"), knownvalue.SetSizeExact(0)),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfigEmptyArray(team, project, alert),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertConfigEmptyArray(team, project, alert+"-updated"),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
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
					return buildThreePartID(rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], rs.Primary.ID), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"actions"},
			},
		},
	})
}

func testAccCheckIssueAlertDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_issue_alert" {
			continue
		}

		ctx := context.Background()
		alert, resp, err := acctest.SharedClient.IssueAlerts.Get(ctx, rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], rs.Primary.ID)
		if err == nil {
			if alert != nil {
				return fmt.Errorf("issue alert %q still exists", rs.Primary.ID)
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckIssueAlertExists(n string, alertId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No project ID is set")
		}

		var resolvedAlertId string
		// Support schema v1 and below
		if value, ok := rs.Primary.Attributes["internal_id"]; ok {
			resolvedAlertId = value
		} else {
			resolvedAlertId = rs.Primary.ID
		}

		ctx := context.Background()
		gotAlert, _, err := acctest.SharedClient.IssueAlerts.Get(ctx, rs.Primary.Attributes["organization"], rs.Primary.Attributes["project"], resolvedAlertId)
		if err != nil {
			return err
		}
		*alertId = sentry.StringValue(gotAlert.ID)
		return nil
	}
}

func testAccIssueAlertConfigEmptyArray(team, project, alert string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30

	conditions = "[]"

	actions = <<EOT
[
	{
		"id": "sentry.mail.actions.NotifyEmailAction",
		"targetType": "IssueOwners"
	}
]
EOT
}
`, team, project, alert)
}

func testAccIssueAlertBaseConfig(team, project, alert, body string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[3]s"

	action_match = "any"
	filter_match = "any"
	frequency    = 30
	%[4]s
}
`, team, project, alert, body)
}

func testAccIssueAlertConfig(team, project, alert string) string {
	return testAccIssueAlertBaseConfig(team, project, alert, `
condition {
	first_seen_event {}
}

condition {
	regression_event {}
}

condition {
	event_frequency {
		value    = 500
		interval = "1h"
	}
}

condition {
	event_unique_user_frequency {
		value    = 1000
		interval = "15m"
	}
}

condition {
	event_frequency_percent {
		value    = 50.0
		interval = "10m"
	}
}

// filter {
// 	age_comparison {
// 		comparison_type = "older"
// 		value           = 10
// 		time            = "minute"
// 	}
// }

action {
	notify_email {
		target_type      = "IssueOwners"
		fallthrough_type = "ActiveMembers"
	}
}

action {
	notify_email {
		target_type       = "Team"
		target_identifier = sentry_team.test.internal_id
		fallthrough_type  = "AllMembers"
	}
}

action {
	slack_notify_service {
		workspace = data.sentry_organization_integration.slack.id
		channel   = "#general"
		tags      = "environment,level"
		notes     = "Please <http://example.com|click here> for triage information"
	}
}

action {
	github_create_ticket {
		integration = data.sentry_organization_integration.github.id
		repo        = "terraform-provider-sentry"
		title       = "My Test Issue"
		assignee    = "jianyuan"
		labels      = ["bug", "enhancement"]
	}
}

action {
	azure_devops_create_ticket {
		integration    = data.sentry_organization_integration.vsts.id
		project        = "123"
		work_item_type = "Microsoft.VSTS.WorkItemTypes.Task"
	}
}

action {
	pagerduty_notify_service {
		account  = sentry_integration_pagerduty.pagerduty.integration_id
		service  = sentry_integration_pagerduty.pagerduty.id
		severity = "default"
	}
}

action {
	opsgenie_notify_team {
		account  = sentry_integration_opsgenie.opsgenie.integration_id
		team     = sentry_integration_opsgenie.opsgenie.id
		priority = "P1"
	}
}

action {
	notify_event {}
}
`) + fmt.Sprintf(`
data "sentry_organization_integration" "slack" {
	organization = sentry_project.test.organization

	provider_key = "slack"
	name         = "A2 Marketing"  # TODO: Use a real integration name
}

data "sentry_organization_integration" "github" {
	organization = sentry_project.test.organization

	provider_key = "github"
	name         = "jianyuan"
}

data "sentry_organization_integration" "vsts" {
	organization = sentry_project.test.organization

	provider_key = "vsts"
	name         = "jianyuanlee"
}

data "sentry_organization_integration" "pagerduty" {
	organization = sentry_project.test.organization

	provider_key = "pagerduty"
	name         = "terraform-provider-sentry"
}

resource "sentry_integration_pagerduty" "pagerduty" {
	organization   = data.sentry_organization_integration.pagerduty.organization
	integration_id = data.sentry_organization_integration.pagerduty.id

	service         = "issue-alert-service"
	integration_key = "issue-alert-integration-key"
}

data "sentry_organization_integration" "opsgenie" {
	organization = sentry_project.test.organization

	provider_key = "opsgenie"
	name         = "terraform-provider-sentry"
}

resource "sentry_integration_opsgenie" "opsgenie" {
	organization    = data.sentry_organization_integration.opsgenie.organization
	integration_id  = data.sentry_organization_integration.opsgenie.id
	team            = "issue-alert-team"
	integration_key = "%[1]s"
}
`, acctest.TestOpsgenieIntegrationKey)
}

func testAccIssueAlertJsonConfig(team, project, alert string) string {
	return testAccIssueAlertBaseConfig(team, project, alert, `
	conditions = <<EOT
[
	{
		"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition",
		"name": "ignored"
	},
	{
		"id": "sentry.rules.conditions.regression_event.RegressionEventCondition"
	},
	{
		"id": "sentry.rules.conditions.event_frequency.EventFrequencyCondition",
		"value": 100,
		"comparisonType": "count",
		"interval": "1h"
	},
	{
		"id": "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
		"value": 100,
		"comparisonType": "count",
		"interval": "1h"
	},
	{
		"id": "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
		"value": 50.0,
		"comparisonType": "count",
		"interval": "1h"
	}
]
EOT

	filters = <<EOT
[
	{
		"id": "sentry.rules.filters.age_comparison.AgeComparisonFilter",
		"value": 10,
		"time": "minute",
		"comparison_type": "older"
	},
	{
		"id": "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
		"value": 10
	},
	{
		"id": "sentry.rules.filters.assigned_to.AssignedToFilter",
		"targetType": "Team",
		"targetIdentifier": ${parseint(sentry_team.test.team_id, 10)}
	},
	{
		"id": "sentry.rules.filters.latest_release.LatestReleaseFilter"
	},
	{
		"id": "sentry.rules.filters.event_attribute.EventAttributeFilter",
		"attribute": "message",
		"match": "co",
		"value": "test"
	},
	{
		"id": "sentry.rules.filters.tagged_event.TaggedEventFilter",
		"key": "test",
		"match": "co",
		"value": "test"
	},
	{
		"id": "sentry.rules.filters.level.LevelFilter",
		"match": "eq",
		"level": "50"
	}
]
EOT

	actions = <<EOT
[
	{
		"id": "sentry.mail.actions.NotifyEmailAction",
		"targetType": "IssueOwners",
		"targetIdentifier": ""
	},
	{
		"id": "sentry.mail.actions.NotifyEmailAction",
		"targetType": "Team",
		"targetIdentifier": ${parseint(sentry_team.test.team_id, 10)}
	},
	{
		"id": "sentry.rules.actions.notify_event.NotifyEventAction"
	}
]
EOT
`)
}
