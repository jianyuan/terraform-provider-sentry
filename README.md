# terraform-provider-sentry
[![CircleCI](https://circleci.com/gh/jianyuan/terraform-provider-sentry/tree/master.svg?style=svg)](https://circleci.com/gh/jianyuan/terraform-provider-sentry/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jianyuan/terraform-provider-sentry)](https://goreportcard.com/report/github.com/jianyuan/terraform-provider-sentry)

Terraform provider for [Sentry](https://sentry.io).

## Installation

See the [the Provider Configuration page of the Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for instructions.

Pre-compiled binaries are available from the [Releases](https://github.com/jianyuan/terraform-provider-sentry/releases) page.

## Usage

### Provider Configuration

#### `sentry`

```
# Configure the Sentry Provider
provider "sentry" {
    token = "${var.sentry_token}"
    base_url = "${var.sentry_base_url}"
}
```

##### Argument Reference

The following arguments are supported:

* `token` - (Required) This is the Sentry authentication token. The value can be sourced from the `SENTRY_TOKEN` environment variable.
* `base_url` - (Optional) This is the target Sentry base API endpoint. The default value is `https://app.getsentry.com/api/`. The value must be provided when working with Sentry On-Premise. The value can be sourced from the `SENTRY_BASE_URL` environment variable.

### Resource Configuration

#### `sentry_organization`

##### Example Usage

```
# Create an organization
resource "sentry_organization" "default" {
    name = "My Organization"
    slug = "my-organization"
    agree_terms = true
}
```

##### Argument Reference

The following arguments are supported:

* `name` - (Required) The human readable name for the organization.
* `slug` - (Optional) The unique URL slug for this organization. If this is not provided a slug is automatically generated based on the name.
* `agree_terms` - (Required) You agree to the applicable terms of service and privacy policy.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created organization.

#### `sentry_team`

##### Example Usage

```
# Create a team
resource "sentry_team" "default" {
    organization = "my-organization"
    name = "My Team"
    slug = "my-team"
}
```

##### Argument Reference

The following arguments are supported:

* `organization` - (Required) The slug of the organization the team should be created for.
* `name` - (Required) The human readable name for the team.
* `slug` - (Optional) The unique URL slug for this team. If this is not provided a slug is automatically generated based on the name.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created team.

#### `sentry_project`

##### Example Usage

```
# Create a project
resource "sentry_project" "default" {
    organization = "my-organization"
    team     = "my-team"
    name     = "Web App"
    slug     = "web-app"
    platform = "javascript"
}
```

##### Argument Reference

The following arguments are supported:

* `organization` - (Required) The slug of the organization the project should be created for.
* `team` - (Required) The slug of the team the project should be created for.
* `name` - (Required) The human readable name for the project.
* `slug` - (Optional) The unique URL slug for this project. If this is not provided a slug is automatically generated based on the name.
* `platform` - (Optional) The integration platform.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created project.

#### `sentry_key`

##### Example Usage

```
# Create a key
resource "sentry_key" "default" {
    organization = "my-organization"
    project = "web-app"
    name = "My Key"
}
```

##### Argument Reference

The following arguments are supported:

* `organization` - (Required) The slug of the organization the key should be created for.
* `project` - (Required) The slug of the project the key should be created for.
* `name` - (Required) The name of the key.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created key.
* `public` - Public key portion of the client key.
* `secret` - Secret key portion of the client key.
* `project_id` - The ID of the project that the key belongs to.
* `is_active` - Flag indicating the key is active.
* `rate_limit_window` - Length of time that will be considered when checking the rate limit.
* `rate_limit_count` - Number of events that can be reported within the rate limit window.
* `dsn_secret` - DSN (Deprecated) for the key.
* `dsn_public` - DSN for the key.
* `dsn_csp` - DSN for the Content Security Policy (CSP) for the key.

#### `sentry_plugin`

##### Example Usage

```
# Create a plugin
resource "sentry_plugin" "default" {
    organization = "my-organization"
    project = "web-app"
    plugin = "slack"
    config = {
      webhook = "slack://webhook"
    }
}
```

##### Argument Reference

The following arguments are supported:

* `organization` - (Required) The slug of the organization the plugin should be enabled for.
* `project` - (Required) The slug of the project the plugin should be enabled for.
* `plugin` - (Required) Identifier of the plugin.
* `config` - (Optional) Configuration of the plugin.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created plugin.

### Data Source Configuration

#### `sentry_key`

##### Example Usage

```
# Retrieve the Default Key
data "sentry_key" "default" {
    organization = "my-organization"
    project = "web-app"
    name = "Default"
}
```

##### Argument Reference

The following arguments are supported:

* `organization` - (Required) The slug of the organization the key should be created for.
* `project` - (Required) The slug of the project the key should be created for.
* `name` - (Optional) The name of the key to retrieve.
* `first` - (Optional) Boolean flag indicating that we want the first key of the returned keys.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created key.
* `public` - Public key portion of the client key.
* `secret` - Secret key portion of the client key.
* `project_id` - The ID of the project that the key belongs to.
* `is_active` - Flag indicating the key is active.
* `rate_limit_window` - Length of time that will be considered when checking the rate limit.
* `rate_limit_count` - Number of events that can be reported within the rate limit window.
* `dsn_secret` - DSN (Deprecated) for the key.
* `dsn_public` - DSN for the key.
* `dsn_csp` - DSN for the Content Security Policy (CSP) for the key.

### Import

You can import existing resources using the [`terraform import`](https://www.terraform.io/docs/import/usage.html) command.

To import an organization:

```bash
$ terraform import sentry_organization.default org-slug
```

To import a team:

```bash
$ terraform import sentry_team.default org-slug/team-slug
```

To import a project:

```bash
$ terraform import sentry_project.default org-slug/project-slug
```

## Development

### Test

Test the provider by running `make test`.

Make sure to set the following environment variables:

- `SENTRY_TEST_ORGANIZATION`
- `SENTRY_TOKEN`

### Build

See the [Writing Custom Providers page of the Terraform documentation](https://www.terraform.io/docs/extend/writing-custom-providers.html#building-the-plugin) for instructions.
