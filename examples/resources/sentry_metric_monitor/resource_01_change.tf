# Issue Detection: Change
# Percentage changes over defined time windows.
resource "sentry_metric_monitor" "change" {
  organization = data.sentry_organization.default.slug
  project      = sentry_project.default.slug

  name = "New change metric monitor"

  owner = {
    team_id = sentry_team.default.internal_id
  }

  aggregate           = "count()"
  dataset             = "events"
  event_types         = ["default", "error"]
  query               = "is:unresolved"
  query_type          = "error"
  time_window_seconds = 3600

  condition_group = {
    conditions = [
      # A high priority issue will be created when query value is 50% lower than the previous 1 hour (issue_detection.comparison_delta).
      {
        type             = "lt"
        comparison       = 50
        condition_result = 75
      },
      # A medium priority issue will be created when query value is 100% lower than the previous 1 hour (issue_detection.comparison_delta).
      {
        type             = "lt"
        comparison       = 100
        condition_result = 50
      },
      # Issue will be resolved when the query value is below or equal to 100% lower than the previous 1 hour.
      {
        type             = "lte"
        comparison       = 100
        condition_result = 0
      },
    ]
  }

  issue_detection = {
    type             = "percent"
    comparison_delta = 3600
  }
}
