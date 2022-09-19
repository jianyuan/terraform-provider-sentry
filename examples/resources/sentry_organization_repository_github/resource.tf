# Retrieve the Github organization integration
data "sentry_organization_integration" "github" {
  organization = "my-organization"
  provider_key = "github"
  name         = "my-github-organization"
}

resource "sentry_organization_repository_github" "this" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.github.internal_id
  identifier     = "my-github-organization/my-github-repo"
}
