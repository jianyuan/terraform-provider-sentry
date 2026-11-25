# Retrieve the Opsgenie organization integration
data "sentry_organization_integration" "opsgenie" {
  organization = "my-organization"

  provider_key = "opsgenie"
  name         = "my-pagerduty-organization"
}

# Associate a Opsgenie service and integration key with a Sentry Opsgenie integration
resource "sentry_integration_opsgenie" "test" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.opsgenie.id

  team            = "my-opsgenie-team"
  integration_key = "c6100908-5c5d-4905-8436-2448fad41bee"
}

# Keep the integration key out of the Terraform state with a write-only
# attribute (Terraform 1.11 or later). Increment `integration_key_wo_version`
# to send a new key.
variable "opsgenie_integration_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "sentry_integration_opsgenie" "write_only" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.opsgenie.id

  team                       = "my-other-opsgenie-team"
  integration_key_wo         = var.opsgenie_integration_key
  integration_key_wo_version = 1
}
