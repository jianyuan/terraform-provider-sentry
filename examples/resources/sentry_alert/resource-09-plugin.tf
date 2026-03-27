resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          plugin = {}
        }
      ]
    }
  ]
}
