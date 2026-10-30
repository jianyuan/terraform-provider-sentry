# Map a Sentry user to their GitLab identity
data "sentry_organization_integration" "gitlab" {
  organization = "my-organization"
  provider_key = "gitlab"
  name         = "my-gitlab-group"
}

resource "sentry_organization_user_mapping" "jane_doe_gitlab" {
  organization = "my-organization"

  user_id           = 12345
  integration_id    = tonumber(data.sentry_organization_integration.gitlab.id)
  external_provider = "gitlab"
  external_id       = "67890"
  external_name     = "@jane.doe"
}
