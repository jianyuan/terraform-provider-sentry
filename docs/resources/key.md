# sentry_key Resource

Sentry Key resource.

## Example Usage

```hcl
# Create a key
resource "sentry_key" "default" {
  organization = "my-organization"
  project      = "web-app"
  name         = "My Key"
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the key should be created for.
- `project` - (Required) The slug of the project the key should be created for.
- `name` - (Required) The name of the key.

## Attribute Reference

The following attributes are exported:

- `id` - The ID of the created key.
- `public` - Public key portion of the client key.
- `secret` - Secret key portion of the client key.
- `project_id` - The ID of the project that the key belongs to.
- `is_active` - Flag indicating the key is active.
- `rate_limit_window` - Length of time that will be considered when checking the rate limit.
- `rate_limit_count` - Number of events that can be reported within the rate limit window.
- `dsn_secret` - DSN (Deprecated) for the key.
- `dsn_public` - DSN for the key.
- `dsn_csp` - DSN for the Content Security Policy (CSP) for the key.
