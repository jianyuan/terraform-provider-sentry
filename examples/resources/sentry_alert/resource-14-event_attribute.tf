resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          event_attribute = {
            attribute = "message"
            match     = "co"
            value     = "bar"
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
