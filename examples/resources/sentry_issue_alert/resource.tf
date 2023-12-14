resource "sentry_issue_alert" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "My issue alert"

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
    "value": 500,
    "interval": "1h"
  },
  {
    "id": "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition",
    "value": 1000,
    "interval": "15m"
  },
  {
    "id": "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition",
    "value": 50,
    "interval": "10m"
  }
]
EOT

  actions = <<EOT
[
  {
    "id" - "sentry.mail.actions.NotifyEmailAction",
    "targetType" - "IssueOwners",
    "fallthroughType" - "ActiveMembers"
  },
  {
    "id": "sentry.mail.actions.NotifyEmailAction",
    "targetType": "Team"
    "fallthroughType": "AllMembers"
    "targetIdentifier": 4524986223
  },
  {
    "id": "sentry.integrations.slack.notify_action.SlackNotifyServiceAction",
    "workspace": ${parseint(data.sentry_organization_integration.slack.id, 10)},
    "channel": "#warning",
    "tags": "environment,level"
  },
  {
    "id": "sentry.integrations.msteams.notify_action.MsTeamsNotifyServiceAction",
    "team": 23465424,
    "channel": "General"
  },
  {
    "id": "sentry.integrations.discord.notify_action.DiscordNotifyServiceAction",
    "server": 63408298,
    "channel_id": 94732897,
    "tags": "browser,user"
  },
  {
    "id": "sentry.integrations.jira.notify_action.JiraCreateTicketAction",
    "integration": 321424,
    "project": "349719"
    "issueType": "1"
  },
  {
    "id": "sentry.integrations.jira_server.notify_action.JiraServerCreateTicketAction",
    "integration": 321424,
    "project": "349719"
    "issueType": "1"
  },
  {
    "id": "sentry.integrations.github.notify_action.GitHubCreateTicketAction",
    "integration": 93749,
    "repo": default,
    "title": "My Test Issue",
    "assignee": "Baxter the Hacker",
    "labels": ["bug", "p1"]
  },
  {
    "id": "sentry.integrations.vsts.notify_action.AzureDevopsCreateTicketAction",
    "integration": 294838,
    "project": "0389485",
    "work_item_type": "Microsoft.VSTS.WorkItemTypes.Task",
  },
  {
    "id": "sentry.integrations.pagerduty.notify_action.PagerDutyNotifyServiceAction",
    "account": 92385907,
    "service": 9823924
  },
  {
    "id": "sentry.integrations.opsgenie.notify_action.OpsgenieNotifyTeamAction",
    "account": 8723897589,
    "team": "9438930258-fairy"
  },
  {
    "id": "sentry.rules.actions.notify_event_service.NotifyEventServiceAction",
    "service": "mail"
  },
  {
    "id": "sentry.rules.actions.notify_event_sentry_app.NotifyEventSentryAppAction",
    "settings": [
        {"name": "title", "value": "Team Rocket"},
        {"name": "summary", "value": "We're blasting off again."},
    ],
    "sentryAppInstallationUuid": 643522
    "hasSchemaFormConfig": true
  },
  {
    "id": "sentry.rules.actions.notify_event.NotifyEventAction"
  }
]
EOT

  filters = <<EOT
[
  {
    "id": "sentry.rules.filters.age_comparison.AgeComparisonFilter",
    "comparison_type": "older",
    "value": 3,
    "time": "week"
  },
  {
    "id": "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter",
    "value": 120
  },
  {
    "id": "sentry.rules.filters.assigned_to.AssignedToFilter",
    "targetType": "Unassigned"
  },
  {
    "id": "sentry.rules.filters.assigned_to.AssignedToFilter",
    "targetType": "Member",
    "targetIdentifier": 895329789
  },
  {
    "id": "sentry.rules.filters.latest_release.LatestReleaseFilter"
  },
  {
    "id": "sentry.rules.filters.issue_category.IssueCategoryFilter",
    "value": 2
  },
  {
    "id": "sentry.rules.conditions.event_attribute.EventAttributeCondition",
    "attribute": "http.url",
    "match": "nc",
    "value": "localhost"
  },
  {
    "id": "sentry.rules.filters.tagged_event.TaggedEventFilter",
    "key": "level",
    "match": "eq"
    "value": "error"
  },
  {
    "id": "sentry.rules.filters.level.LevelFilter",
    "match": "gte"
    "level": "50"
  }
]
EOT
}

# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = sentry_project.test.organization

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}
