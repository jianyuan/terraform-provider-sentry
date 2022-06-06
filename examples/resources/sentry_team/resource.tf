# Create a team
resource "sentry_team" "default" {
  organization = "my-organization"

  name = "My Team"
  slug = "my-team"
}
