# Retrieve a GitHub integration
data "sentry_organization_integration" "github" {
  organization = sentry_project.test.organization

  provider_key = "github"
  name         = "GitHub"
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          github = {
            integration_id = data.sentry_organization_integration.github.id
            repo           = "terraform-provider-sentry"
            assignee       = "jianyuan"
            labels         = ["bug"]
          }
        }
      ]
    }
  ]
}
