# Issue Detection: Threshold
# Absolute-valued thresholds, for non-seasonal data.
resource "sentry_metric_monitor" "threshold" {
  organization = data.sentry_organization.default.slug
  project      = sentry_project.default.slug

  name = "New threshold metric monitor"

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
      # A high priority issue will be created when query value is greater than 100.
      {
        type             = "gt"
        comparison       = 100
        condition_result = 75
      },
      # A medium priority issue will be created when query value is greater than 50.
      {
        type             = "gt"
        comparison       = 50
        condition_result = 50
      },
      # Issue will be resolved when the query value is below or equal to 50.
      {
        type             = "lte"
        comparison       = 50
        condition_result = 0
      },
    ]
  }

  issue_detection = {
    type = "static"
  }
}
