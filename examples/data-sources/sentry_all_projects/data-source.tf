# Retrieve a list of projects available to the authenticated user
data "sentry_projects" "default" {
  organization = "my-organization"
}
