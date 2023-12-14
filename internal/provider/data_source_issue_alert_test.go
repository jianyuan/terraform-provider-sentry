package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccIssueAlertDataSource(t *testing.T) {
	rn := "sentry_issue_alert.test"
	rnCopy := "sentry_issue_alert.test_copy"
	dsn := "data.sentry_issue_alert.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string
	var alertIdCopy string

	checkResourceAttrJsonPair := func(a, b, attr string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			resA, ok := s.RootModule().Resources[a]
			if !ok {
				return fmt.Errorf("resource %s not found", a)
			}

			resB, ok := s.RootModule().Resources[b]
			if !ok {
				return fmt.Errorf("resource %s not found", b)
			}

			given := jsontypes.NewNormalizedValue(resA.Primary.Attributes[attr])
			expected := jsontypes.NewNormalizedValue(resB.Primary.Attributes[attr])
			match, diags := expected.StringSemanticEquals(context.Background(), given)
			if !match {
				return fmt.Errorf("expected %s, got %s: %s", expected, given, diags)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertDataSourceConfig(team, project, alert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIssueAlertExists(rn, &alertId),
					resource.TestCheckResourceAttr(dsn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(dsn, "project", project),
					resource.TestCheckResourceAttrPair(dsn, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(dsn, "project", rn, "project"),
					checkResourceAttrJsonPair(dsn, rn, "conditions"),
					checkResourceAttrJsonPair(dsn, rn, "filters"),
					checkResourceAttrJsonPair(dsn, rn, "actions"),
					resource.TestCheckResourceAttrPair(dsn, "action_match", rn, "action_match"),
					resource.TestCheckResourceAttrPair(dsn, "filter_match", rn, "filter_match"),
					resource.TestCheckResourceAttrPair(dsn, "frequency", rn, "frequency"),
					resource.TestCheckResourceAttrPair(dsn, "name", rn, "name"),
					resource.TestCheckResourceAttrPair(dsn, "environment", rn, "environment"),
					testAccCheckIssueAlertExists(rnCopy, &alertIdCopy),
					resource.TestCheckResourceAttrPair(rnCopy, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(rnCopy, "project", rn, "project"),
					checkResourceAttrJsonPair(rnCopy, rn, "conditions"),
					checkResourceAttrJsonPair(rnCopy, rn, "filters"),
					checkResourceAttrJsonPair(rnCopy, rn, "actions"),
					resource.TestCheckResourceAttr(rnCopy, "action_match", "all"),
					resource.TestCheckResourceAttr(rnCopy, "filter_match", "all"),
					resource.TestCheckResourceAttrPair(rnCopy, "frequency", rn, "frequency"),
					resource.TestCheckResourceAttr(rnCopy, "name", alert+"-copy"),
					resource.TestCheckResourceAttrPair(rnCopy, "environment", rn, "environment"),
				),
			},
		},
	})
}

// func TestAccIssueAlertDataSource_MigrateFromPluginSDK(t *testing.T) {
// 	dsn := "data.sentry_issue_alert.test"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() { acctest.PreCheck(t) },
// 		Steps: []resource.TestStep{
// 			{
// 				ExternalProviders: map[string]resource.ExternalProvider{
// 					acctest.ProviderName: {
// 						Source:            "jianyuan/sentry",
// 						VersionConstraint: "0.11.2",
// 					},
// 				},
// 				Config: testAccIssueAlertDataSourceConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrSet(dsn, "id"),
// 					resource.TestCheckResourceAttrPair(dsn, "internal_id", dsn, "id"),
// 					resource.TestCheckResourceAttr(dsn, "organization", acctest.TestOrganization),
// 					resource.TestCheckResourceAttr(dsn, "provider_key", "github"),
// 					resource.TestCheckResourceAttr(dsn, "name", "jianyuan"),
// 				),
// 			},
// 			{
// 				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 				Config:                   testAccIssueAlertDataSourceConfig,
// 				PlanOnly:                 true,
// 			},
// 		},
// 	})
// }

func testAccIssueAlertDataSourceConfig(teamName string, projectName string, alertName string) string {
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
		"id": "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
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

data "sentry_issue_alert" "test" {
	id           = sentry_issue_alert.test.id
	organization = sentry_issue_alert.test.organization
	project      = sentry_issue_alert.test.project
}

resource "sentry_issue_alert" "test_copy" {
	organization = data.sentry_issue_alert.test.organization
	project      = data.sentry_issue_alert.test.project
	name         = "${data.sentry_issue_alert.test.name}-copy"

	action_match = "all"
	filter_match = "all"
	frequency    = data.sentry_issue_alert.test.frequency

	conditions = data.sentry_issue_alert.test.conditions
	filters    = data.sentry_issue_alert.test.filters
	actions    = data.sentry_issue_alert.test.actions
}
`, teamName, projectName, alertName)
}
