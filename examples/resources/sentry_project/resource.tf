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
}
