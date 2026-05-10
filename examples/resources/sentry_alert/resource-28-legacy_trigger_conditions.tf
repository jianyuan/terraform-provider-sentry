resource "sentry_metric_monitor" "default" {
  # ...
}

resource "sentry_alert" "legacy" {
  organization = "my-organization"
  name         = "My Alert (with legacy trigger conditions)"
  environment  = "production"

  monitor_ids       = [sentry_metric_monitor.default.id]
  frequency_minutes = 1440

  # All natively supported trigger condition types
  trigger_conditions = [
    { first_seen_event = {} },
    { issue_resolved_trigger = {} },
    { reappeared_event = {} },
    { regression_event = {} },
  ]

  # Trigger condition types that are not configurable in the latest Sentry alerts
  # UI similar to trigger_conditions. Set explicitly to preserve them;
  # omit to remove them on the next apply.
  legacy_trigger_conditions = [
    "new_high_priority_issue",
    "existing_high_priority_issue",
  ]

  action_filters = [
    {
      logic_type = "all"
      conditions = []
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
