resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          issue_priority_greater_or_equal = {
            comparison = 75
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
