resource "sentry_organization_data_scrubbing" "default" {
  organization = "my-organization"

  data_scrubber          = true
  data_scrubber_defaults = true
  scrub_ip_addresses     = true
  sensitive_fields       = ["email", "phone"]
  safe_fields            = ["order_id"]

  # [Mask] [Credit card numbers] from [$message] and [Remove] [Anything] from [extra.foo]
  relay_pii_config = jsonencode({
    rules = {
      "0" = { type = "creditcard", redaction = { method = "mask" } }
      "1" = { type = "anything", redaction = { method = "remove" } }
    }
    applications = {
      "$message"  = ["0"]
      "extra.foo" = ["1"]
    }
  })
}
