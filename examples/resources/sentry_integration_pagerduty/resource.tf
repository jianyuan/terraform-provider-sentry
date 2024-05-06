# Retrieve the PagerDuty organization integration
data "sentry_organization_integration" "pagerduty" {
  organization = "my-organization"

  provider_key = "pagerduty"
  name         = "my-pagerduty-organization"
}

# Associate a PagerDuty service and integration key with a Sentry PagerDuty integration
resource "sentry_integration_pagerduty" "test" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.pagerduty.id

  service         = "my-pagerduty-service"
  integration_key = "my-pagerduty-integration-key"
}
