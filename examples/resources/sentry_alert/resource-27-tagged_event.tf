resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          tagged_event = {
            key   = "level"
            match = "eq"
            value = "error"
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
