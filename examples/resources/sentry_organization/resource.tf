# Create an organization
resource "sentry_organization" "default" {
  name = "My Organization"
  slug = "my-organization"

  agree_terms = true
}
