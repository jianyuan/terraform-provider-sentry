resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          issue_priority_deescalating = {}
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
