# Issue Detection: Threshold
# Absolute-valued thresholds, for non-seasonal data.
resource "sentry_metric_monitor" "threshold" {
  organization = data.sentry_organization.default.slug
  project      = sentry_project.default.slug

  name = "New threshold metric monitor"

  owner = {
    team_id = sentry_team.default.internal_id
  }

  aggregate   = "count()"
  dataset     = "events"
  event_types = ["default", "error"]

  condition_group = {
    conditions = [
      {
        type             = "gt"
        comparison       = 100
        condition_result = 75
      },
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
