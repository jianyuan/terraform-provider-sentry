resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          age_comparison = {
            comparison_type = "older"
            time            = "minute"
            value           = 1
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
