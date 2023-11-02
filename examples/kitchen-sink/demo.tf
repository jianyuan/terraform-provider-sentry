terraform {
  required_providers {
    sentry = {
      source = "jianyuan/sentry"
    }
  }
}

provider "sentry" {
  token = "xxx"
}

data "sentry_organization" "main" {
  # Taken from URL: https://sentry.io/organizations/[slug]/issues/
  slug = "my-sentry-organization"
}

output "organization" {
  value = data.sentry_organization.main
}

#
# Team
#

resource "sentry_team" "main" {
  organization = data.sentry_organization.main.id
  name         = "My team"
}

output "team" {
  value = sentry_team.main.id
}

#
# Project
#

resource "sentry_project" "main" {
  organization = sentry_team.main.organization
  team         = sentry_team.main.id
  name         = "My project"
  platform     = "python"
}

output "project" {
  value = sentry_project.main.id
}

#
# Project key
#

data "sentry_key" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id

  first = true
}

output "project_key" {
  value = data.sentry_key.main
}

#
# Issue alert
#

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
      name           = "The issue is seen more than 200 times in 1h"
      value          = 200
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

output "issue_alert_id" {
  value = sentry_issue_alert.main.internal_id
}

output "issue_alert_url" {
  value = "https://sentry.io/organizations/${sentry_issue_alert.main.organization}/alerts/rules/${sentry_issue_alert.main.project}/${sentry_issue_alert.main.internal_id}/details/"
}

#
# Metric alert
#

data "sentry_organization_integration" "slack" {
  organization = sentry_project.main.organization
  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}

resource "sentry_metric_alert" "main" {
  organization      = sentry_project.main.organization
  project           = sentry_project.main.id
  name              = "My metric alert"
  dataset           = "events"
  query             = ""
  aggregate         = "count()"
  time_window       = 60
  threshold_type    = 0
  resolve_threshold = 0

  trigger {
    action {
      type              = "email"
      target_type       = "team"
      target_identifier = sentry_team.main.team_id
    }
    alert_threshold = 300
    label           = "critical"
    threshold_type  = 0
  }

  trigger {
    action {
      type              = "slack"
      target_type       = "specific"
      target_identifier = "#slack-channel"
      integration_id    = data.sentry_organization_integration.slack.id
    }
    alert_threshold = 300
    label           = "critical"
    threshold_type  = 0
  }

  trigger {
    alert_threshold = 100
    label           = "warning"
    threshold_type  = 0
  }
}

output "metric_alert_id" {
  value = sentry_metric_alert.main.internal_id
}

output "metric_alert_url" {
  value = "https://sentry.io/organizations/${sentry_metric_alert.main.organization}/alerts/rules/details/${sentry_metric_alert.main.internal_id}/"
}

#
# Dashboard
#

resource "sentry_dashboard" "main" {
  organization = data.sentry_organization.main.id
  title        = "My dashboard"

  widget {
    title        = "Number of Errors"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "!event.type:transaction"
      order_by   = "count()"
    }

    layout {
      x     = 0
      y     = 0
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "Number of Issues"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count_unique(issue)"]
      aggregates = ["count_unique(issue)"]
      conditions = "!event.type:transaction"
      order_by   = "count_unique(issue)"
    }

    layout {
      x     = 1
      y     = 0
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "Events"
    display_type = "line"
    interval     = "5m"
    widget_type  = "discover"

    query {
      name       = "Events"
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "!event.type:transaction"
      order_by   = "count()"
    }

    layout {
      x     = 2
      y     = 0
      w     = 4
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Affected Users"
    display_type = "line"
    interval     = "5m"
    widget_type  = "discover"

    query {
      name       = "Known Users"
      fields     = ["count_unique(user)"]
      aggregates = ["count_unique(user)"]
      conditions = "has:user.email !event.type:transaction"
      order_by   = "count_unique(user)"
    }

    query {
      name       = "Anonymous Users"
      fields     = ["count_unique(user)"]
      aggregates = ["count_unique(user)"]
      conditions = "!has:user.email !event.type:transaction"
      order_by   = "count_unique(user)"
    }

    layout {
      x     = 1
      y     = 2
      w     = 1
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Handled vs. Unhandled"
    display_type = "line"
    interval     = "5m"
    widget_type  = "discover"

    query {
      name       = "Handled"
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "error.handled:true"
      order_by   = "count()"
    }

    query {
      name       = "Unhandled"
      fields     = ["count()"]
      aggregates = ["count()"]
      conditions = "error.handled:false"
      order_by   = "count()"
    }

    layout {
      x     = 0
      y     = 2
      w     = 1
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Errors by Country"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["geo.country_code", "geo.region", "count()"]
      aggregates = ["count()"]
      conditions = "!event.type:transaction has:geo.country_code"
      order_by   = "count()"
    }

    layout {
      x     = 4
      y     = 6
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "High Throughput Transactions"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["count()", "transaction"]
      aggregates = ["count()"]
      columns    = ["transaction"]
      conditions = "!event.type:error"
      order_by   = "-count()"
    }

    layout {
      x     = 0
      y     = 6
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Errors by Browser"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["browser.name", "count()"]
      aggregates = ["count()"]
      columns    = ["browser.name"]
      conditions = "!event.type:transaction has:browser.name"
      order_by   = "-count()"
    }

    layout {
      x     = 5
      y     = 2
      w     = 1
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Overall User Misery"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["user_misery(300)"]
      aggregates = ["user_misery(300)"]
    }

    layout {
      x     = 0
      y     = 1
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "Overall Apdex"
    display_type = "big_number"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["apdex(300)"]
      aggregates = ["apdex(300)"]
    }

    layout {
      x     = 1
      y     = 1
      w     = 1
      h     = 1
      min_h = 1
    }
  }

  widget {
    title        = "High Throughput Transactions"
    display_type = "top_n"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["transaction", "count()"]
      aggregates = ["count()"]
      columns    = ["transaction"]
      conditions = "!event.type:error"
      order_by   = "-count()"
    }

    layout {
      x     = 0
      y     = 4
      w     = 2
      h     = 2
      min_h = 2
    }
  }

  widget {
    title        = "Issues Assigned to Me or My Teams"
    display_type = "table"
    interval     = "5m"
    widget_type  = "issue"

    query {
      fields     = ["assignee", "issue", "title"]
      columns    = ["assignee", "issue", "title"]
      conditions = "assigned_or_suggested:me is:unresolved"
      order_by   = "priority"
    }

    layout {
      x     = 2
      y     = 2
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Transactions Ordered by Misery"
    display_type = "table"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["transaction", "user_misery(300)"]
      aggregates = ["user_misery(300)"]
      columns    = ["transaction"]
      order_by   = "-user_misery(300)"
    }

    layout {
      x     = 2
      y     = 6
      w     = 2
      h     = 4
      min_h = 2
    }
  }

  widget {
    title        = "Errors by Browser Over Time"
    display_type = "top_n"
    interval     = "5m"
    widget_type  = "discover"

    query {
      fields     = ["browser.name", "count()"]
      aggregates = ["count()"]
      columns    = ["browser.name"]
      conditions = "event.type:error has:browser.name"
      order_by   = "-count()"
    }

    layout {
      x     = 4
      y     = 2
      w     = 1
      h     = 4
      min_h = 2
    }
  }
}

output "dashboard_id" {
  value = sentry_dashboard.main.internal_id
}

output "dashboard_url" {
  value = "https://sentry.io/organizations/${sentry_dashboard.main.organization}/dashboard/${sentry_dashboard.main.internal_id}/"
}
