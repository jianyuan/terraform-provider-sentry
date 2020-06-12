# sentry_organization Resource

Sentry Organization resource.

## Example Usage

```hcl
# Create an organization
resource "sentry_organization" "default" {
  name        = "My Organization"
  slug        = "my-organization"
  agree_terms = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The human readable name for the organization.
- `slug` - (Optional) The unique URL slug for this organization. If this is not provided a slug is automatically generated based on the name.
- `agree_terms` - (Required) You agree to the applicable terms of service and privacy policy.

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created organization.
