# sentry_organization Data Source

Sentry Organization data source.

## Example Usage

```hcl
# Retrieve the organization
data "sentry_organization" "org" {
  slug = "my-organization"
}
```

## Argument Reference

The following arguments are supported:

- `slug` - (required) The unique URL slug for this organization.

## Attribute Reference

The following attributes are exported:

- `id` - The ID for this organization.
- `name` - The human readable name for this organization.
- `slug` - The unique URL slug for this organization.
