# GitHub
data "sentry_organization_integration" "github" {
  organization = "my-organization"
  provider_key = "github"
  name         = "my-github-organization"
}

resource "sentry_organization_repository" "github" {
  organization     = "my-organization"
  integration_type = "github"
  integration_id   = data.sentry_organization_integration.github.internal_id
  identifier       = "my-github-organization/my-github-repo"
}

# Azure DevOps
data "sentry_organization_repository" "vsts" {
  organization = "my-organization"
  provider_key = "vsts"
  name         = "my-azure-devops-organization"
}

resource "sentry_organization_repository" "vsts" {
  organization     = "my-organization"
  integration_type = "vsts"
  integration_id   = data.sentry_organization_integration.vsts.internal_id
  identifier       = "5febef5a-833d-4e14-b9c0-14cb638f91e6"
}
