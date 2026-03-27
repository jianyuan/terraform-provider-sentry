resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          event_unique_user_frequency_count = {
            interval = "5m"
            value    = 1
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
