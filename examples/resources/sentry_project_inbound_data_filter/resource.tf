resource "sentry_project" "default" {
  organization = "my-organization"

  teams = ["my-first-team", "my-second-team"]
  name  = "web-app"

  platform = "javascript"
}

# Create an inbound data filter for a project
resource "sentry_project_inbound_data_filter" "test" {
  organization = sentry_project.default.organization
  project      = sentry_project.default.id
  filter_id    = "browser-extensions"
  active       = true
}

# Create an inbound data filter with subfilters. Only applicable to the
# `legacy-browser` filter.
resource "sentry_project_inbound_data_filter" "test" {
  organization = sentry_project.default.organization
  project      = sentry_project.default.id
  filter_id    = "legacy-browser"
  subfilters   = ["ie_pre_9", "ie9"]
}
