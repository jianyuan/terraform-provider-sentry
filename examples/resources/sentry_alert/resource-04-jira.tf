# Retrieve a Jira integration
data "sentry_organization_integration" "jira" {
  organization = sentry_project.test.organization

  provider_key = "jira"
  name         = "JIRA" # Name of your Jira server
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          jira = {
            integration_id = data.sentry_organization_integration.jira.id
          }
        }
      ]
    }
  ]
}
