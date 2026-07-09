resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          percent_sessions_count = {
            interval = "1h"
            value    = 10
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
