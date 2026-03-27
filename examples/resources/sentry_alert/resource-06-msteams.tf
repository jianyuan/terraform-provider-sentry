# Retrieve a MS Teams integration
data "sentry_organization_integration" "msteams" {
  organization = sentry_project.test.organization

  provider_key = "msteams"
  name         = "My Team" # Name of your Microsoft Teams team
}

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
