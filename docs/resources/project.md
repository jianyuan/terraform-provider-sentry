# sentry_project Resource

Sentry Project resource.

## Example Usage

```hcl
# Create a project
resource "sentry_project" "default" {
  organization = "my-organization"
  team         = "my-team"
  name         = "Web App"
  slug         = "web-app"
  platform     = "javascript"
  resolve_age  = 720
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the project should be created for.
- `team` - (Required) The slug of the team the project should be created for.
- `name` - (Required) The human readable name for the project.
- `slug` - (Optional) The unique URL slug for this project. If this is not provided a slug is automatically generated based on the name.
- `platform` - (Optional) The integration platform.
- `resolve_age` - (Optional) Hours in which an issue is automatically resolve if not seen after this amount of time.

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created project.
