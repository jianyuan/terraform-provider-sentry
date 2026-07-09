resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          issue_category = {
            value = 1
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
