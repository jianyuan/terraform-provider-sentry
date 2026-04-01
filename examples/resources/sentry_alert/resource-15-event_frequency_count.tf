resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          event_frequency_count = {
            interval = "1m"
            value    = 100
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
