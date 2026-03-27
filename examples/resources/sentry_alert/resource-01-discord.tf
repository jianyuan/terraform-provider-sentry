resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          discord = {
            channel_id     = "714123428994482189"
            integration_id = data.sentry_organization_integration.discord.id
            tags           = "environment, level"
          }
        }
      ]
    }
  ]
}
