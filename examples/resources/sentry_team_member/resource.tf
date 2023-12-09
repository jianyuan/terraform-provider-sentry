# Add a member to a team
resource "sentry_organization_member" "default" {
  organization = "my-organization"
  email        = "test@example.com"
  role         = "member"
}

resource "sentry_team" "default" {
  organization = "my-organization"
  name         = "my-team"
  slug         = "my-team"
}

resource "sentry_team_member" "default" {
  organization = "my-organization"
  team_slug    = sentry_team.default.slug
  member_id    = sentry_organization_member.default.internal_id
}
