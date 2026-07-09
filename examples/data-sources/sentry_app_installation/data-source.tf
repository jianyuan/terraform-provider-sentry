# Look up a Sentry App Installation by app slug
data "sentry_app_installation" "example" {
  organization = "my-org"
  slug         = "my-sentry-app"
}

# Use in an issue alert notify_event_sentry_app action
resource "sentry_issue_alert" "example" {
  # ...
  actions_v2 = [{
    notify_event_sentry_app = {
      sentry_app_installation_uuid = data.sentry_app_installation.example.uuid
      settings                     = {}
    }
  }]
}

# Use in a metric alert sentry_app trigger action
resource "sentry_metric_alert" "example" {
  # ...
  trigger {
    action {
      type              = "sentry_app"
      target_type       = "sentry_app"
      target_identifier = tostring(data.sentry_app_installation.example.sentry_app_id)
      sentry_app_id     = data.sentry_app_installation.example.sentry_app_id
      integration_id    = 0
    }
    alert_threshold = 100
    label           = "critical"
    threshold_type  = 0
  }
}
