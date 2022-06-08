resource "sentry_team" "test_team" {
  organization = "%s"
  name         = "Test team"
}

resource "sentry_project" "test_project" {
  organization = "%s"
  team         = sentry_team.test_team.id
  name         = "Test project"
  platform     = "go"
}

resource "sentry_issue_alert" "test_issue_alert" {
  organization = "%s"
  project      = sentry_project.test_project.id
  name         = "Test rule"

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
      interval       = "1h"
      id             = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
      comparisonType = "count"
      value          = 100
      name           = "The issue is seen more than 100 times in 1h"
    },
    {
      interval       = "1h"
      id             = "sentry.rules.conditions.event_frequency.EventUniqueUserFrequencyCondition"
      comparisonType = "count"
      value          = 100
      name           = "The issue is seen by more than 100 users in 1h"
    },
    {
      interval       = "1h"
      id             = "sentry.rules.conditions.event_frequency.EventFrequencyPercentCondition"
      comparisonType = "count"
      value          = 100
      name           = "The issue affects more than 100.0 percent of sessions in 1h"
    },
  ]

  filters = [
    {
      id               = "sentry.rules.filters.assigned_to.AssignedToFilter"
      name             = "The issue is assigned to Unassigned"
      targetIdentifier = ""
      targetType       = "Unassigned"
    }
  ]

  actions = [
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      name             = "Send a notification to IssueOwners"
      targetIdentifier = ""
      targetType       = "IssueOwners"
    }
  ]
}
