resource "sentry_project" "default" {
  organization = "my-organization"

  teams = ["my-first-team", "my-second-team"]
  name  = "web-app"

  platform = "javascript"
}

# Enable spike protection for the project
resource "sentry_project_spike_protection" "default" {
  organization = sentry_project.default.organization
  project_slug = sentry_project.default.slug
  enabled      = true
}