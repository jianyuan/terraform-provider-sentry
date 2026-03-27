resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          jira_server = {
            integration_id = data.sentry_organization_integration.jira_server.id
          }
        }
      ]
    }
  ]
}
