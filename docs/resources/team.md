# sentry_team Resource

Sentry Team resource.

## Example Usage

```hcl
# Create a team
resource "sentry_team" "default" {
    organization = "my-organization"
    name         = "My Team"
    slug         = "my-team"
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the team should be created for.
- `name` - (Required) The human readable name for the team.
- `slug` - (Optional) The unique URL slug for this team. If this is not provided a slug is automatically generated based on the name.

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created team.
