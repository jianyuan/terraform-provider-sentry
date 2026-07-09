resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          issue_occurrences = {
            value = 1
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
