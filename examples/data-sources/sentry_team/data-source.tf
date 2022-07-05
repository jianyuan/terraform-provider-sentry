# Retrieve a team
data "sentry_team" "default" {
  organization = "my-organization"

  slug = "my-team"
}
