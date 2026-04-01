# Retrieve a VSTS integration
data "sentry_organization_integration" "vsts" {
  organization = "my-org"

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
            project        = "0389485"
            work_item_type = "Microsoft.VSTS.WorkItemTypes.Task"
          }
        }
      ]
    }
  ]
}
