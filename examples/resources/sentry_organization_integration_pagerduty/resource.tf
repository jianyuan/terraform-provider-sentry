# Retrieve the PagerDuty organization integration
data "sentry_organization_integration" "pagerduty" {
  organization = "my-organization"
  provider_key = "pagerduty"
  name         = "name-of-pagerduty-integration"
}

resource "sentry_organization_integration_pagerduty" "this" {
  organization    = "my-organization"
  integration_id  = data.sentry_organization_integration.pagerduty.internal_id
  service_name    = "name-of-pagerduty-service"
  integration_key = "integration-key-from-pagerduty"
}
