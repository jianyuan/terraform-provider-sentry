# Retrieve a VSTS integration
data "sentry_organization_integration" "vsts" {
  organization = sentry_project.test.organization

  provider_key = "vsts"
  name         = "Azure DevOps"
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          vsts = {
            integration_id = data.sentry_organization_integration.vsts.id
          }
        }
      ]
    }
  ]
}
