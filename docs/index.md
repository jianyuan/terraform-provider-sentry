# Sentry Provider

Terraform provider for [Sentry](https://sentry.io).

## Example Usage

```hcl
# Configure the Sentry Provider
provider "sentry" {
  token = var.sentry_token
  base_url = var.sentry_base_url
}
```

## Argument Reference

The following arguments are supported:

- `token` - (Required) This is the Sentry authentication token. The value can be sourced from the `SENTRY_TOKEN` environment variable.
- `base_url` - (Optional) This is the target Sentry base API endpoint. The default value is `https://app.getsentry.com/api/`. The value must be provided when working with Sentry On-Premise. The value can be sourced from the `SENTRY_BASE_URL` environment variable.
