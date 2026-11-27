data "sentry_organization_integration" "slack" {
  organization = "my-organization"

  provider_key = "slack"
  name         = "Slack Workspace"
}

data "sentry_organization_integration" "github" {
  organization = "my-organization"

  provider_key = "github"
  name         = "my-github-organization"
}

resource "sentry_team" "backend" {
  organization = "my-organization"
  name         = "Backend"
  slug         = "backend"
}

# Send the team's notifications to a Slack channel
resource "sentry_external_team" "backend_slack" {
  organization   = sentry_team.backend.organization
  team           = sentry_team.backend.slug
  provider_name  = "slack"
  external_name  = "#backend-alerts"
  integration_id = data.sentry_organization_integration.slack.id
  external_id    = "C0123456789" # optional: the Slack channel ID
}

# Resolve the "@my-github-organization/backend" CODEOWNERS handle to the Sentry team
resource "sentry_external_team" "backend_github" {
  organization   = sentry_team.backend.organization
  team           = sentry_team.backend.slug
  provider_name  = "github"
  external_name  = "@my-github-organization/backend"
  integration_id = data.sentry_organization_integration.github.id
}
