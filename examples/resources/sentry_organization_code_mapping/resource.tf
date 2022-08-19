# Retrieve the Github organization integration
data "sentry_organization_integration" "github" {
  organization = "my-organization"
  provider_key = "github"
  name         = "my-github-organization"
}

resource "sentry_project" "this" {
  organization = "my-organization"

  team = "my-team"
  name = "Web App"
  slug = "web-app"

  platform    = "javascript"
  resolve_age = 720
}

resource "sentry_organization_repository_github" "this" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.github.internal_id
  identifier     = "my-github-organization/my-github-repo"
}

resource "sentry_organization_code_mapping" "this" {
  organization   = "my-organization"
  integration_id = data.sentry_organization_integration.github.internal_id
  repository_id  = sentry_organization_repository_github.this.internal_id
  project_id     = sentry_project.this.internal_id

  default_branch = "main"
  stack_root     = "/"
  source_root    = "src/"
}
