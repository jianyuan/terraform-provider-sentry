package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccIssueAlertResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig("value", "value", "value", `
					actions = "[]"

					conditions    = "[]"
					conditions_v2 = [
						{ first_seen_event = {} },
					]
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute "conditions" cannot be specified when "conditions_v2" is specified`),
			},
			{
				Config: testAccIssueAlertConfig("value", "value", "value", `
					actions = "[]"

					conditions_v2 = []
				`),
				ExpectError: acctest.ExpectLiteralError(`Attribute conditions_v2 list must contain at least 1 elements, got: 0`),
			},
			{
				Config: testAccIssueAlertConfig("value", "value", "value", `
					actions = "[]"

					conditions_v2 = [
						{ first_seen_event = {}, regression_event = {} },
					]
				`),
				ExpectError: acctest.ExpectLiteralError(
					`Attribute "conditions_v2[0].first_seen_event" cannot be specified when`,
					`"conditions_v2[0].regression_event" is specified`,
				),
			},
		},
	})
}

func TestAccIssueAlertResource_basic(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
		// statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.NotNull()),
		// statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.NotNull()),
		// statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig(team, project, alert, `
					actions = <<EOT
						[{"id": "sentry.rules.actions.notify_event.NotifyEventAction"}]
					EOT

					conditions_v2 = [
						{ first_seen_event = {} },
						{ regression_event = {} },
						{ reappeared_event = {} },
						{ new_high_priority_issue = {} },
						{ existing_high_priority_issue = {} },
						{ event_frequency = { comparison_type = "count", value = 100, interval = "1h" } },
						{ event_frequency = { comparison_type = "percent", comparison_interval = "1w", value = 100, interval = "1h" } },
						{ event_unique_user_frequency = { comparison_type = "count", value = 100, interval = "1h" } },
						{ event_unique_user_frequency = { comparison_type = "percent", comparison_interval = "1w", value = 100, interval = "1h" } },
						{ event_frequency_percent = { comparison_type = "count", value = 100, interval = "1h" } },
						{ event_frequency_percent = { comparison_type = "percent", comparison_interval = "1w", value = 100, interval = "1h" } }
					]
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"first_seen_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("A new issue is created"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"regression_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("The issue changes state from resolved to unresolved"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"reappeared_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("The issue changes state from ignored to unresolved"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"new_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Sentry marks a new issue as high priority"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"existing_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Sentry marks an existing issue as high priority"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("The issue is seen more than 100 times in 1h"),
								"comparison_type":     knownvalue.StringExact("count"),
								"comparison_interval": knownvalue.Null(),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("The issue is seen more than 100 times in 1h"),
								"comparison_type":     knownvalue.StringExact("percent"),
								"comparison_interval": knownvalue.StringExact("1w"),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_unique_user_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("The issue is seen by more than 100 users in 1h"),
								"comparison_type":     knownvalue.StringExact("count"),
								"comparison_interval": knownvalue.Null(),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_unique_user_frequency": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("The issue is seen by more than 100 users in 1h"),
								"comparison_type":     knownvalue.StringExact("percent"),
								"comparison_interval": knownvalue.StringExact("1w"),
								"value":               knownvalue.Int64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency_percent": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("The issue affects more than 100.0 percent of sessions in 1h"),
								"comparison_type":     knownvalue.StringExact("count"),
								"comparison_interval": knownvalue.Null(),
								"value":               knownvalue.Float64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"event_frequency_percent": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("The issue affects more than 100.0 percent of sessions in 1h"),
								"comparison_type":     knownvalue.StringExact("percent"),
								"comparison_interval": knownvalue.StringExact("1w"),
								"value":               knownvalue.Float64Exact(100),
								"interval":            knownvalue.StringExact("1h"),
							}),
						}),
					})),
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert+"-updated", `
					actions = <<EOT
						[{"id": "sentry.rules.actions.notify_event.NotifyEventAction"}]
					EOT

					conditions_v2 = [
						{ reappeared_event = {} },
						{ new_high_priority_issue = {} },
						{ existing_high_priority_issue = {} },
					]
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions_v2"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"reappeared_event": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("The issue changes state from ignored to unresolved"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"new_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Sentry marks a new issue as high priority"),
							}),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"existing_high_priority_issue": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Sentry marks an existing issue as high priority"),
							}),
						}),
					})),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
			},
		},
	})
}

func TestAccIssueAlertResource_jsonValues(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert-with-a-very-looooong-name-greater-than-64-characters")
	var alertId string

	check := func(alert string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckIssueAlertExists(rn, &alertId),
			resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
				if alertId != value {
					return fmt.Errorf("expected %s, got %s", alertId, value)
				}
				return nil
			}),
		)
	}

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("conditions"), knownvalue.StringExact("[]")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filters"), knownvalue.StringExact("[]")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("actions"), knownvalue.StringExact("[]")),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig_jsonValues(team, project, alert),
				Check:  check(alert),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertConfig_jsonValues(team, project, alert+"-updated"),
				Check:  check(alert + "-updated"),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert+"-updated")),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
			},
		},
	})
}

func TestAccIssueAlertResource_upgradeFromVersion(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "id"),
					resource.TestCheckResourceAttrSet(rn, "internal_id"),
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "name", alert),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIssueAlertExists(rn, &alertId),
					resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
						if alertId != value {
							return fmt.Errorf("expected %s, got %s", alertId, value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "name", alert),
					resource.TestCheckResourceAttr(rn, "action_match", "any"),
					resource.TestCheckResourceAttr(rn, "filter_match", "any"),
					resource.TestCheckResourceAttr(rn, "frequency", "30"),
					resource.TestCheckResourceAttrSet(rn, "conditions"),
					resource.TestCheckResourceAttrSet(rn, "actions"),
				),
			},
		},
	})
}

func TestAccIssueAlertResource_emptyArray(t *testing.T) {
	rn := "sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

	check := func(alert string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckIssueAlertExists(rn, &alertId),
			resource.TestCheckResourceAttrWith(rn, "id", func(value string) error {
				if alertId != value {
					return fmt.Errorf("expected %s, got %s", alertId, value)
				}
				return nil
			}),
			resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
			resource.TestCheckResourceAttr(rn, "project", project),
			resource.TestCheckResourceAttr(rn, "name", alert),
			resource.TestCheckResourceAttr(rn, "action_match", "any"),
			resource.TestCheckResourceAttr(rn, "filter_match", "any"),
			resource.TestCheckResourceAttr(rn, "frequency", "30"),
			resource.TestCheckResourceAttrSet(rn, "conditions"),
			resource.TestCheckResourceAttrSet(rn, "actions"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfigEmptyArray(team, project, alert),
				Check:  check(alert),
			},
			{
				Config: testAccIssueAlertConfigEmptyArray(team, project, alert+"-updated"),
				Check:  check(alert + "-updated"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateIdFunc: acctest.ThreePartImportStateIdFunc(rn, "organization", "project"),
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

func testAccIssueAlertConfigEmptyArray(teamName string, projectName string, alertName string) string {
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
`, teamName, projectName, alertName)
}

func testAccIssueAlertConfig(team string, project string, alert string, extras string) string {
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
`, team, project, alert, extras)
}

func testAccIssueAlertConfig_jsonValues(team string, project string, alert string) string {
	return testAccIssueAlertConfig(team, project, alert, `
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
