resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          latest_release = {}
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
