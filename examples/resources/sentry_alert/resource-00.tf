resource "sentry_cron_monitor" "default" {
  # ...
}

resource "sentry_metric_monitor" "default" {
  # ...
}

resource "sentry_uptime_monitor" "default" {
  # ...
}


resource "sentry_alert" "default" {
  organization = "my-organization"
  name         = "My Alert"
  environment  = "production"

  # Trigger when any of the monitors are triggered
  monitor_ids = [
    sentry_cron_monitor.default.id,
    sentry_metric_monitor.default.id,
    sentry_uptime_monitor.default.id,
  ]

  frequency_minutes = 1440

  trigger_conditions = [
    "first_seen_event",
    "issue_resolved_trigger",
    "reappeared_event",
    "regression_event",
  ]

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          event_frequency_percent = {
            comparison_interval = "1w"
            interval            = "1h"
            value               = 100
          }
        }
      ]
      actions = [
        {
          email = {
            target_type      = "issue_owners"
            fallthrough_type = "AllMembers"
          }
        }
      ]
    }
  ]
}
