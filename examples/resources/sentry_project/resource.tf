# Create a project
resource "sentry_project" "default" {
  organization = "my-organization"

  team = "my-team"
  name = "Web App"
  slug = "web-app"

  platform    = "javascript"
  resolve_age = 720
}
