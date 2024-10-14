# Retrieve a list of projects available to the authenticated user
data "sentry_all_projects" "default" {
  organization = "my-organization"
}
