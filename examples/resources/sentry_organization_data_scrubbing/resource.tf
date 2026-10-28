resource "sentry_organization_data_scrubbing" "default" {
  organization = "my-organization"

  data_scrubber          = true
  data_scrubber_defaults = true
  sensitive_fields       = ["email", "phone", "iban"]
  safe_fields            = ["order_id"]
  scrub_ip_addresses     = true
}
