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
    alert_threshold = 100
    label           = "warning"
    threshold_type  = 0
  }
}
