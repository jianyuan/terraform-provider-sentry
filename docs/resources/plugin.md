# sentry_plugin Resource

Sentry Plugin resource.

## Example Usage

```hcl
# Create a plugin
resource "sentry_plugin" "default" {
  organization = "my-organization"
  project      = "web-app"
  plugin       = "slack"

  config = {
    webhook = "slack://webhook"
  }
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the plugin should be enabled for.
- `project` - (Required) The slug of the project the plugin should be enabled for.
- `plugin` - (Required) Identifier of the plugin.
- `config` - (Optional) Configuration of the plugin.

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created plugin.
