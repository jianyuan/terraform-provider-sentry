# Create an organization member
resource "sentry_organization_member" "john_doe" {
  email = "test@example.com"
  role  = "member"
  teams = ["my-team"]
}
