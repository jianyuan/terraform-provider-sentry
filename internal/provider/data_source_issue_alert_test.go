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

func TestAccIssueAlertDataSource(t *testing.T) {
	rn := "sentry_issue_alert.test"
	dsn := "data.sentry_issue_alert.test"
	project := acctest.RandomWithPrefix("tf-project")
	alert := acctest.RandomWithPrefix("tf-issue-alert")
	var alertId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIssueAlertDataSourceConfig(project, alert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIssueAlertExists(rn, &alertId),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dsn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
					statecheck.ExpectKnownValue(dsn, tfjsonpath.New("project"), knownvalue.StringExact(project)),
					statecheck.CompareValuePairs(dsn, tfjsonpath.New("action_match"), rn, tfjsonpath.New("action_match"), compare.ValuesSame()),
					statecheck.CompareValuePairs(dsn, tfjsonpath.New("filter_match"), rn, tfjsonpath.New("filter_match"), compare.ValuesSame()),
					statecheck.CompareValuePairs(dsn, tfjsonpath.New("frequency"), rn, tfjsonpath.New("frequency"), compare.ValuesSame()),
					statecheck.CompareValuePairs(dsn, tfjsonpath.New("name"), rn, tfjsonpath.New("name"), compare.ValuesSame()),
					statecheck.CompareValuePairs(dsn, tfjsonpath.New("environment"), rn, tfjsonpath.New("environment"), compare.ValuesSame()),
				},
			},
		},
	})
}

func testAccIssueAlertDataSourceConfig(projectName string, alertName string) string {
	return fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = "%[1]s"
	teams        = ["%[3]s"]
	name         = "%[4]s"
	platform     = "go"
}

resource "sentry_issue_alert" "test" {
	organization = "%[1]s"
	project      = sentry_project.test.id
	name         = "%[5]s"

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
		"targetIdentifier": %[2]s
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
		"targetIdentifier": %[2]s
	}
]
EOT
}

data "sentry_issue_alert" "test" {
	id           = sentry_issue_alert.test.id
	organization = sentry_issue_alert.test.organization
	project      = sentry_issue_alert.test.project
}
`, acctest.TestOrganization, acctest.TestTeam.Id, acctest.TestTeam.Slug, projectName, alertName)
}
