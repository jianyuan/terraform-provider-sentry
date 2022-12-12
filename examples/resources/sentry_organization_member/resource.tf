# Create an organization member
resource "sentry_organization_member" "john_doe" {
  organization = "my-organization"

  email = "test@example.com"
  role  = "member"
  teams = ["my-team"]
}
