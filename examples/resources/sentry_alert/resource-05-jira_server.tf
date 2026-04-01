# Retrieve a Jira Server integration
data "sentry_organization_integration" "jira_server" {
  organization = "my-org"
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
            project        = "349719"
            issue_type     = "1"
          }
        }
      ]
    }
  ]
}
