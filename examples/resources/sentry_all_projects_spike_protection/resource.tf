# Enable spike protection for several projects in a Sentry organization.
resource "sentry_project" "web-app" {
  organization = "my-organization"

  teams = ["my-first-team"]
  name  = "web-app"
  slug  = "web-app"

  platform = "go"
}

resource "sentry_project" "mobile-app" {
  organization = "my-organization"

  teams = ["my-second-team"]
  name  = "mobile-app"
  slug  = "mobile-app"

  platform = "android"
}

resource "sentry_all_projects_spike_protection" "main" {
  organization = "my-organization"
  projects = [
    sentry_project.web-app.id,
    sentry_project.mobile-app.id,
  ]
  enabled = true
}

# Use the `sentry_all_projects` data source to get all projects in a Sentry organization and enable spike protection for all of them.
data "sentry_all_projects" "all" {
  organization = "my-organization"
}

resource "sentry_all_projects_spike_protection" "main" {
  organization = data.sentry_all_projects.all.organization
  projects     = data.sentry_all_projects.all.project_slugs
  enabled      = true
}
