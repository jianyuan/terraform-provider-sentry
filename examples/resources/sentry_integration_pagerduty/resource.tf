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

# Keep the integration key out of the Terraform state with a write-only
# attribute (Terraform 1.11 or later). Increment `integration_key_wo_version`
# to send a new key.
variable "pagerduty_integration_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "sentry_integration_pagerduty" "write_only" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.pagerduty.id

  service                    = "my-other-pagerduty-service"
  integration_key_wo         = var.pagerduty_integration_key
  integration_key_wo_version = 1
}
