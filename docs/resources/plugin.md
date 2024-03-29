---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sentry_plugin Resource - terraform-provider-sentry"
subcategory: ""
description: |-
  Sentry Plugin resource.
---

# sentry_plugin (Resource)

Sentry Plugin resource.

## Example Usage

```terraform
# Create a plugin
resource "sentry_plugin" "default" {
  organization = "my-organization"

  project = "web-app"
  plugin  = "slack"

  config = {
    webhook = "slack://webhook"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization` (String) The slug of the organization the project belongs to.
- `plugin` (String) Plugin ID.
- `project` (String) The slug of the project to create the plugin for.

### Optional

- `config` (Map of String) Plugin config.

### Read-Only

- `id` (String) The ID of this resource.
