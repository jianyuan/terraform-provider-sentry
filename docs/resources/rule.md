# sentry_rule Resource

Sentry Rule resource. Note that there's no public documentation for the values of conditions, filters, and actions. You can either inspect the request payload sent when creating or editing an alert rule on Sentry or inspect [Sentry's rules registry in the source code](https://github.com/getsentry/sentry/tree/master/src/sentry/rules).

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
      id   = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
      name = "A new issue is created"
    }
  ]

  filters = [
    {
      id         = "sentry.rules.filters.assigned_to.AssignedToFilter"
      targetType = "Unassigned"
    }
  ]

  actions = [
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      name             = "Send an email to IssueOwners"
      targetIdentifier = ""
      targetType       = "IssueOwners"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the plugin should be enabled for.
- `project` - (Required) The slug of the project the plugin should be enabled for.
- `name` - (Required) Name for this alert.
- `action_match` - (Optional) Use `all` to trigger alerting when all conditions are met, and `any` when at. least a condition is met. Defaults to `any`.
- `frequency` - (Optional) Perform actions at most once every `X` minutes for this issue. Defaults to `30`.
- `environment` - (Optional) Environment for these conditions to apply to.
- `conditions` - (Required) List of conditions.
- `filters` - (Optional) List of filters.
- `actions` - (Required) List of actions.

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created rule.
- `name` - The name of the created rule.
- `frequency` - The rule's frequency.
- `environment` - The rule's environment.
