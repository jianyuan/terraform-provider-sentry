# Retrieve an organization
data "sentry_organization" "org" {
  slug = "my-organization"
}
