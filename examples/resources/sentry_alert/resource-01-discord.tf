# Retrieve a Discord integration
data "sentry_organization_integration" "discord" {
  organization = sentry_project.test.organization

  provider_key = "discord"
  name         = "Discord Server" # Name of your Discord server
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          discord = {
            channel_id     = "123456789012345678"
            integration_id = data.sentry_organization_integration.discord.id
            tags           = "environment, level"
          }
        }
      ]
    }
  ]
}
