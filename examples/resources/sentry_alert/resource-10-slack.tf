# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = sentry_project.test.organization

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          slack = {
            integration_id = data.sentry_organization_integration.slack.id
            channel_name   = "#general"
            tags           = "environment,level"
            notes          = "Please <http://example.com|click here> for triage information."
          }
        }
      ]
    }
  ]
}
