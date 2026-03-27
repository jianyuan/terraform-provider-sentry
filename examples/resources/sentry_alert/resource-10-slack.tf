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
