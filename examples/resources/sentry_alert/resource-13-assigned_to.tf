resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          assigned_to = {
            target_type = "Team"
            target_id   = sentry_team.default.internal_id
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
