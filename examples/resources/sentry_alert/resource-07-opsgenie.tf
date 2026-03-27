# Retrieve an Opsgenie integration
data "sentry_organization_integration" "opsgenie" {
  organization = sentry_project.test.organization
  provider_key = "opsgenie"
  name         = "Opsgenie"
}

# Configure an Opsgenie integration team and integration key
resource "sentry_integration_opsgenie" "opsgenie" {
  organization    = data.sentry_organization_integration.opsgenie.organization
  integration_id  = data.sentry_organization_integration.opsgenie.id
  team            = "issue-alert-team"
  integration_key = "my-integration-key"
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          opsgenie = {
            integration_id = sentry_integration_opsgenie.opsgenie.integration_id
            priority       = "P1"
            team_id        = sentry_integration_opsgenie.opsgenie.id
            team_name      = sentry_integration_opsgenie.opsgenie.team
          }
        }
      ]
    }
  ]
}
