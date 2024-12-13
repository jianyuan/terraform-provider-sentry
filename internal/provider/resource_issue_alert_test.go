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

func TestAccIssueAlertResource(t *testing.T) {
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
			resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
			resource.TestCheckResourceAttr(rn, "project", project),
			resource.TestCheckResourceAttr(rn, "name", alert),
			resource.TestCheckResourceAttr(rn, "action_match", "any"),
			resource.TestCheckResourceAttr(rn, "filter_match", "any"),
			resource.TestCheckResourceAttr(rn, "frequency", "30"),
			resource.TestCheckResourceAttrSet(rn, "conditions"),
			resource.TestCheckResourceAttrSet(rn, "filters"),
			resource.TestCheckResourceAttrSet(rn, "actions"),
		)
	}

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("action_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("filter_match"), knownvalue.StringExact("any")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("frequency"), knownvalue.Int64Exact(30)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("environment"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertConfig(team, project, alert),
				Check:  check(alert),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(alert)),
				),
			},
			{
				Config: testAccIssueAlertConfig(team, project, alert+"-updated"),
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

func testAccIssueAlertConfig(teamName string, projectName string, alertName string) string {
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
}
`, teamName, projectName, alertName)
}
