resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          msteams = {
            integration_id = data.sentry_organization_integration.msteams.id
            team_id        = "my-team-id"
            channel_name   = "General"
          }
        }
      ]
    }
  ]
}
