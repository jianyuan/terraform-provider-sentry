resource "sentry_issue_alert" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "My issue alert"

  action_match = "any"
  filter_match = "any"
  frequency    = 30

  conditions = [
    # A new issue is created
    {
      id = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
    },

    # The issue changes state from resolved to unresolved
    {
      id = "sentry.rules.conditions.regression_event.RegressionEventCondition"
    },

    # The issue is seen more than 100 times in 1h
    {
      id             = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
      value          = 100
      comparisonType = "count"
      interval       = "1h"
    },

    # The issue is seen by more than 100 users in 1h
    {
      id             = "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
      value          = 100
      comparisonType = "count"
      interval       = "1h"
    },

    # The issue affects more than 50.0 percent of sessions in 1h
    {
      id             = "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
      value          = 50.0
      comparisonType = "count"
      interval       = "1h"
    },
  ]

  filters = [
    # The issue is older than 10 minute
    {
      id              = "sentry.rules.filters.age_comparison.AgeComparisonFilter"
      value           = 10
      time            = "minute"
      comparison_type = "older"
    },

    # The issue has happened at least 10 times
    {
      id    = "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter"
      value = 10
    },

    # The issue is assigned to Team
    {
      id               = "sentry.rules.filters.assigned_to.AssignedToFilter"
      targetType       = "Team"
      targetIdentifier = sentry_team.main.team_id
    },

    # The event is from the latest release
    {
      id = "sentry.rules.filters.latest_release.LatestReleaseFilter"
    },

    # The event's message value contains test
    {
      id        = "sentry.rules.filters.event_attribute.EventAttributeFilter"
      attribute = "message"
      match     = "co"
      value     = "test"
    },

    # The event's tags match test contains test
    {
      id    = "sentry.rules.filters.tagged_event.TaggedEventFilter"
      key   = "test"
      match = "co"
      value = "test"
    },

    # The event's level is equal to fatal
    {
      id    = "sentry.rules.filters.level.LevelFilter"
      match = "eq"
      level = "50"
    }
  ]

  actions = [
    # Send a notification to IssueOwners
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      targetType       = "IssueOwners"
      targetIdentifier = ""
    },

    # Send a notification to Team
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      targetType       = "Team"
      targetIdentifier = sentry_team.main.team_id
    },

    # Send a notification (for all legacy integrations)
    {
      id = "sentry.rules.actions.notify_event.NotifyEventAction"
    },

    # Send a notification to the Slack workspace to #general
    {
      id      = "sentry.integrations.slack.notify_action.SlackNotifyServiceAction"
      channel = "#general"

      # From: https://sentry.io/settings/[org-slug]/integrations/slack/[slack-integration-id]/
      # Or use the sentry_organization_integration data source to retrieve the integration ID:
      workspace = data.sentry_organization_integration.slack.internal_id
    },
  ]
}

# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = sentry_project.test.organization

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}
