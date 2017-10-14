# terraform-provider-sentry
[![CircleCI](https://circleci.com/gh/jianyuan/terraform-provider-sentry/tree/master.svg?style=svg)](https://circleci.com/gh/jianyuan/terraform-provider-sentry/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jianyuan/terraform-provider-sentry)](https://goreportcard.com/report/github.com/jianyuan/terraform-provider-sentry)

Terraform provider for Sentry

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
}
```

##### Argument Reference

The following arguments are supported:

* `name` - (Required) The human readable name for the organization.
* `slug` - (Optional) The unique URL slug for this organization. If this is not provided a slug is automatically generated based on the name.

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
    team = "my-team"
    name = "Web App"
    slug = "web-app"
}
```

##### Argument Reference

The following arguments are supported:

* `organization` - (Required) The slug of the organization the project should be created for.
* `team` - (Required) The slug of the team the project should be created for.
* `name` - (Required) The human readable name for the project.
* `slug` - (Optional) The unique URL slug for this project. If this is not provided a slug is automatically generated based on the name.

##### Attributes Reference

The following attributes are exported:

* `id` - The ID of the created project.

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

### Import

You can import existing resources using [terraform import](https://www.terraform.io/docs/import/index.html).

Organization are directly importable using `terraform import sentry_organization.default org-slug`. Teams and project via `terraform import sentry_team.default org-slug/team-slug` and `terraform import sentry_project.default org-slug/project-slug` respectively.
