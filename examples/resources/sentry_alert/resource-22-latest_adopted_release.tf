resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      conditions = [
        {
          latest_adopted_release = {
            age_comparison   = "older"
            environment      = "test"
            release_age_type = "oldest"
          }
        }
      ]
      actions = [
        # ...
      ]
    }
  ]
}
