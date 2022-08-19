# Retrieve the Github organization integration
# Organization integration should be named after the Github organization
data "sentry_organization_integration" "github" {
  organization = "my-organization"
  provider_key = "github"
  name         = "my-github-organization"
}
