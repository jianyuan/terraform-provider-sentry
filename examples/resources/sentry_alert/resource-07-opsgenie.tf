resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          opsgenie = {
            integration_id = sentry_integration_opsgenie.opsgenie.integration_id
            priority       = "P1"
            team_id        = sentry_integration_opsgenie.opsgenie.id
            team_name      = sentry_integration_opsgenie.opsgenie.team
          }
        }
      ]
    }
  ]
}
