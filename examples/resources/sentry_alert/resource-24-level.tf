resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          level = {
            level = 50
            match = "eq"
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
