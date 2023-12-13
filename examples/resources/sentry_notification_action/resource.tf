resource "sentry_project" "default" {
  organization = "my-organization"

  teams = ["my-first-team", "my-second-team"]
  name  = "web-app"

  platform = "javascript"
}

# Create a notification action for the project
resource "sentry_notification_action" "default" {
  organization      = sentry_project.default.organization
  trigger_type      = "spike-protection"
  service_type      = "sentry_notification"
  target_identifier = "default"
  target_display    = "default"
  projects          = [sentry_project.default.id]
}
