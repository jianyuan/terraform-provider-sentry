# Retrieve a Jira integration
data "sentry_organization_integration" "jira" {
  organization = "my-org"
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
            project        = "349719"
            issue_type     = "1"
          }
        }
      ]
    }
  ]
}
