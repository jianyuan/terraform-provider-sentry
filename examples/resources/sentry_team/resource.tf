# Create a team
resource "sentry_team" "default" {
  organization = "my-organization"

  name = "my-team"
  slug = "my-team"
}
