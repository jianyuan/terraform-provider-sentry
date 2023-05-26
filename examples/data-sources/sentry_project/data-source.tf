# Retrieve a project
data "sentry_project" "default" {
  organization = "my-organization"

  slug = "my-project"
}
