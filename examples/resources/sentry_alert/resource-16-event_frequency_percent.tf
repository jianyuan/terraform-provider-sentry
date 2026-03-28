resource "sentry_alert" "default" {
  # ...

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
        # ...
      ]
    }
  ]
}
