# Retrieve a PagerDuty integration
data "sentry_organization_integration" "pagerduty" {
  organization = "my-org"
  provider_key = "pagerduty"
  name         = "PagerDuty"
}

# Configure a PagerDuty integration service and integration key
resource "sentry_integration_pagerduty" "pagerduty" {
  organization    = data.sentry_organization_integration.pagerduty.organization
  integration_id  = data.sentry_organization_integration.pagerduty.id
  service         = "issue-alert-service"
  integration_key = "issue-alert-integration-key"
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          pagerduty = {
            integration_id = sentry_integration_pagerduty.pagerduty.integration_id
            service_id     = sentry_integration_pagerduty.pagerduty.id
            service_name   = sentry_integration_pagerduty.pagerduty.service
            severity       = "default"
          }
        }
      ]
    }
  ]
}
