resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          vsts = {
            integration_id = data.sentry_organization_integration.vsts.id
          }
        }
      ]
    }
  ]
}
