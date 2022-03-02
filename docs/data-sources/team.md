# sentry_team Data Source

Sentry Team data source.

## Example Usage

```hcl
# Retrieve the team
data "sentry_team" "app_team" {
  organization = "my-organization-slug"
  slug         = "some-team"
}
```

## Argument Reference

The following arguments are supported:

- `slug` - (required) The unique URL slug for this team.

## Attribute Reference

The following attributes are exported:

- `id` - The internal ID for this team.
- `name` - The human readable name for this team.
- `slug` - The unique URL slug for this team.
- `team_id` - The internal ID for this team.
