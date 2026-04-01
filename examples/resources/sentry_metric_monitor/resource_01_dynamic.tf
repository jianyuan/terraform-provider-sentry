# Issue Detection: Dynamic
# Auto-detect anomalies and mean deviation, for seasonal/noisy data.
resource "sentry_metric_monitor" "dynamic" {
  organization = data.sentry_organization.default.slug
  project      = sentry_project.default.slug

  name = "New dynamic metric monitor"

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
      {
        type                      = "anomaly_detection"
        comparison_sensitivity    = "high"
        comparison_threshold_type = "above_and_below"
        condition_result          = 75
      },
    ]
  }

  issue_detection = {
    type = "dynamic"
  }
}
