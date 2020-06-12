# sentry_rule Resource

Sentry Rule resource.

## Example Usage

```hcl
# Create a plugin
resource "sentry_rule" "default" {
  organization = "my-organization"
  project      = "web-app"
  action_match = "any"
  frequency    = 30
  environment  = "production"

  conditions = [
    {
      id       = "sentry.rules.conditions.event_frequency.EventFrequencyCondition"
      value    = 500
      interval = "1h"
    }
  ]

  actions = [
    {
      id        = "sentry.integrations.slack.notify_action.SlackNotifyServiceAction"
      channel   = "#alerts"
      workspace = "12345"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the plugin should be enabled for.
- `project` - (Required) The slug of the project the plugin should be enabled for.
- `action_match` - (Optional) Use `all` to trigger alerting when all conditions are met, and `any` when at least a condition is met. Defaults to `any`.
- `frequency` - (Optional) Perform actions at most once every `X` minutes for this issue. Defaults to `30`.
- `environment` - (Optional) Environment name
- `actions` - (Required) List of actions
- `conditions` - (Required) List of conditions

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created rule.
- `name` - The name of the created rule.
- `actions` - The rule's actions.
- `conditions` - The rule's conditions.
- `frequency` - The rule's frequency.
- `environment` - The rule's environment.
