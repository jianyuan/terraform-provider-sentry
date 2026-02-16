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
      target_identifier = sentry_team.main.internal_id
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
      input_channel_id  = "C0XXXXXXXXX"
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

# Example: Metric Alert with Sentry App Action
data "sentry_organization_integration" "rootly" {
  organization = "my-organization"
  name         = "Rootly"
}

resource "sentry_metric_alert" "main" {
  organization   = "my-organization"
  project        = "my-project"
  name           = "My Alert with Rootly"
  dataset        = "events"
  event_types    = ["error"]
  query          = ""
  aggregate      = "count()"
  time_window    = 60
  threshold_type = 0

  trigger {
    action {
      type              = "sentry_app"
      target_type       = "sentry_app"
      target_identifier = data.sentry_organization_integration.rootly.internal_id
      sentry_app_id     = tonumber(data.sentry_organization_integration.rootly.internal_id)
      integration_id    = 0
    }
    alert_threshold = 100
    label           = "critical"
    threshold_type  = 0
  }
}
