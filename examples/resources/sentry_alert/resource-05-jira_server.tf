# Retrieve a Jira Server integration
data "sentry_organization_integration" "jira_server" {
  organization = sentry_project.test.organization

  provider_key = "jira_server"
  name         = "JIRA" # Name of your Jira server
}

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
