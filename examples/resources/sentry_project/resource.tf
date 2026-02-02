# Create a project
resource "sentry_project" "default" {
  organization = "my-organization"

  teams = ["my-first-team", "my-second-team"]
  name  = "Web App"
  slug  = "web-app"

  platform    = "javascript"
  resolve_age = 720

  default_rules = false

  filters = {
    blacklisted_ips = ["127.0.0.1", "0.0.0.0/8"]
    releases        = ["1.*", "[!3].[0-9].*"]
    error_messages  = ["TypeError*", "*: integer division or modulo by zero"]
  }

  fingerprinting_rules  = <<-EOT
    # force all errors of the same type to have the same fingerprint
    error.type:DatabaseUnavailable -> system-down
    # force all memory allocation errors to be grouped together
    stack.function:malloc -> memory-allocation-error
  EOT
  grouping_enhancements = <<-EOT
    # remove all frames above a certain function from grouping
    stack.function:panic_handler ^-group
    # mark all functions following a prefix in-app
    stack.function:mylibrary_* +app
  EOT

  highlight_tags = ["release", "environment"]
}
