# sentry_key Data Source

Sentry Key data source.

## Example Usage

```hcl
# Retrieve the Default Key
data "sentry_key" "default" {
  organization = "my-organization"
  project      = "web-app"
  name         = "Default"
}
```

## Argument Reference

The following arguments are supported:

- `organization` - (Required) The slug of the organization the key should be created for.
- `project` - (Required) The slug of the project the key should be created for.
- `name` - (Optional) The name of the key to retrieve.
- `first` - (Optional) Boolean flag indicating that we want the first key of the returned keys.

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
