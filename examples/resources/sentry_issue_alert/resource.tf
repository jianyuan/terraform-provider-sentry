resource "sentry_issue_alert" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "My issue alert"

  action_match = "any"
  filter_match = "any"
  frequency    = 30

  conditions = [
    {
      id   = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
      name = "A new issue is created"
    },
    {
      id   = "sentry.rules.conditions.regression_event.RegressionEventCondition"
      name = "The issue changes state from resolved to unresolved"
    },
    {
      id             = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
      name           = "The issue is seen more than 100 times in 1h"
      value          = 100
      comparisonType = "count"
      interval       = "1h"
    },
    {
      id             = "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
      name           = "The issue is seen by more than 100 users in 1h"
      value          = 100
      comparisonType = "count"
      interval       = "1h"
    },
    {
      id             = "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
      name           = "The issue affects more than 50.0 percent of sessions in 1h"
      value          = 50.0
      comparisonType = "count"
      interval       = "1h"
    },
  ]

  filters = [
    {
      id              = "sentry.rules.filters.age_comparison.AgeComparisonFilter"
      name            = "The issue is older than 10 minute"
      value           = 10
      time            = "minute"
      comparison_type = "older"
    },
    {
      id    = "sentry.rules.filters.issue_occurrences.IssueOccurrencesFilter"
      name  = "The issue has happened at least 10 times"
      value = 10
    },
    {
      id               = "sentry.rules.filters.assigned_to.AssignedToFilter"
      name             = "The issue is assigned to Team"
      targetType       = "Team"
      targetIdentifier = sentry_team.main.team_id
    },
    {
      id   = "sentry.rules.filters.latest_release.LatestReleaseFilter"
      name = "The event is from the latest release"
    },
    {
      id        = "sentry.rules.filters.event_attribute.EventAttributeFilter"
      name      = "The event's message value contains test"
      attribute = "message"
      match     = "co"
      value     = "test"
    },
    {
      id    = "sentry.rules.filters.tagged_event.TaggedEventFilter"
      name  = "The event's tags match test contains test"
      key   = "test"
      match = "co"
      value = "test"
    },
    {
      id    = "sentry.rules.filters.level.LevelFilter"
      name  = "The event's level is equal to fatal"
      match = "eq"
      level = "50"
    }
  ]

  actions = [
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      name             = "Send a notification to IssueOwners"
      targetType       = "IssueOwners"
      targetIdentifier = ""
    },
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      name             = "Send a notification to Team"
      targetType       = "Team"
      targetIdentifier = sentry_team.main.team_id
    },
    {
      id   = "sentry.rules.actions.notify_event.NotifyEventAction"
      name = "Send a notification (for all legacy integrations)"
    }
  ]
}
